package scaleway

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dirien/devpod-provider-scaleway/pkg/options"
	"github.com/loft-sh/devpod/pkg/client"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/pkg/errors"
	"github.com/scaleway/scaleway-sdk-go/api/account/v2"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type ScalewayProvider struct {
	Config           *options.Options
	InstanceAPI      *instance.API
	AccountAPI       *account.API
	Log              log.Logger
	WorkingDirectory string
}

func NewProvider(logs log.Logger, init bool) (*ScalewayProvider, error) {
	scwAccessKey := os.Getenv("SCW_ACCESS_KEY")
	if scwAccessKey == "" {
		return nil, errors.Errorf("SCW_ACCESS_KEY is not set")
	}

	scwSecretKey := os.Getenv("SCW_SECRET_KEY")
	if scwSecretKey == "" {
		return nil, errors.Errorf("SCW_SECRET_KEY is not set")
	}

	config, err := options.FromEnv(init)

	if err != nil {
		return nil, err
	}

	zone, err := scw.ParseZone(config.Zone)
	if err != nil {
		return nil, err
	}

	client, err := scw.NewClient(
		scw.WithAuth(scwAccessKey, scwSecretKey),
		scw.WithDefaultOrganizationID(config.OrganizationID),
		scw.WithDefaultZone(zone),
	)
	if err != nil {
		return nil, err
	}
	provider := &ScalewayProvider{
		Config:      config,
		Log:         logs,
		InstanceAPI: instance.NewAPI(client),
		AccountAPI:  account.NewAPI(client),
	}
	return provider, nil
}

func GetDevpodInstance(scalewayProvider *ScalewayProvider) (*instance.GetServerResponse, error) {
	servers, err := scalewayProvider.InstanceAPI.ListServers(&instance.ListServersRequest{
		Tags: []string{scalewayProvider.Config.MachineID},
	})
	if err != nil {
		return nil, err
	}

	return scalewayProvider.InstanceAPI.GetServer(&instance.GetServerRequest{
		ServerID: servers.Servers[0].ID,
	})
}

func Create(scalewayProvider *ScalewayProvider) error {
	publicKeyBase, err := ssh.GetPublicKeyBase(scalewayProvider.Config.MachineFolder)
	if err != nil {
		return err
	}
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase)
	if err != nil {
		return err
	}

	sizeGB, _ := strconv.Atoi(scalewayProvider.Config.DiskSizeGB)

	server, err := scalewayProvider.InstanceAPI.CreateServer(&instance.CreateServerRequest{
		CommercialType: scalewayProvider.Config.CommercialType,
		Image:          scalewayProvider.Config.Image,
		Tags:           []string{scalewayProvider.Config.MachineID},
		Volumes: map[string]*instance.VolumeServerTemplate{
			"0": {
				Size: scw.Size(sizeGB) * scw.GB,
			},
		},
		DynamicIPRequired: scw.BoolPtr(true),
	})
	if err != nil {
		return err
	}

	err = scalewayProvider.InstanceAPI.SetServerUserData(&instance.SetServerUserDataRequest{
		ServerID: server.Server.ID,
		Key:      "cloud-init",
		Content: strings.NewReader(fmt.Sprintf(`#cloud-config
users:
- name: devpod
  shell: /bin/bash
  groups: [ sudo ]
  ssh_authorized_keys:
  - %s
`, string(publicKey))),
	})
	if err != nil {
		return err
	}
	duration := 2 * time.Second
	err = scalewayProvider.InstanceAPI.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
		ServerID:      server.Server.ID,
		Action:        instance.ServerActionPoweron,
		RetryInterval: &duration,
	})
	if err != nil {
		return err
	}

	return nil
}

func Delete(scalewayProvider *ScalewayProvider) error {
	devPodInstance, err := GetDevpodInstance(scalewayProvider)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	duration := 2 * time.Second
	err = scalewayProvider.InstanceAPI.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
		ServerID:      devPodInstance.Server.ID,
		Action:        instance.ServerActionPoweroff,
		RetryInterval: &duration,
	})
	if err != nil {
		return err
	}
	err = scalewayProvider.InstanceAPI.DeleteServer(&instance.DeleteServerRequest{
		ServerID: devPodInstance.Server.ID,
	})
	for _, volume := range devPodInstance.Server.Volumes {
		err = scalewayProvider.InstanceAPI.DeleteVolume(&instance.DeleteVolumeRequest{
			VolumeID: volume.ID,
		})
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func Start(scalewayProvider *ScalewayProvider) error {
	devPodInstance, err := GetDevpodInstance(scalewayProvider)
	if err != nil {
		return err
	}

	duration := 2 * time.Second
	err = scalewayProvider.InstanceAPI.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
		ServerID:      devPodInstance.Server.ID,
		Action:        instance.ServerActionPoweron,
		RetryInterval: &duration,
	})
	if err != nil {
		return err
	}

	return nil
}

func Status(scalewayProvider *ScalewayProvider) (client.Status, error) {
	devPodInstance, err := GetDevpodInstance(scalewayProvider)
	if err != nil {
		return client.StatusNotFound, nil
	}

	switch {
	case devPodInstance.Server.State == instance.ServerStateRunning:
		return client.StatusRunning, nil
	case devPodInstance.Server.State == instance.ServerStateStopped:
		return client.StatusStopped, nil
	default:
		return client.StatusBusy, nil
	}
}

func Stop(scalewayProvider *ScalewayProvider) error {
	devPodInstance, err := GetDevpodInstance(scalewayProvider)
	if err != nil {
		return err
	}

	duration := 2 * time.Second
	err = scalewayProvider.InstanceAPI.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
		ServerID:      devPodInstance.Server.ID,
		Action:        instance.ServerActionPoweroff,
		RetryInterval: &duration,
	})
	if err != nil {
		return err
	}
	return nil
}

func Init(scalewayProvider *ScalewayProvider) error {
	_, err := scalewayProvider.InstanceAPI.ListServers(&instance.ListServersRequest{})
	if err != nil {
		return err
	}
	return nil
}
