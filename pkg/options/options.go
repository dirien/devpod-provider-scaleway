package options

import (
	"fmt"
	"os"
)

var (
	SCW_DEFAULT_PROJECT_ID      = "SCW_DEFAULT_PROJECT_ID"
	SCW_DEFAULT_ORGANIZATION_ID = "SCW_DEFAULT_ORGANIZATION_ID"
	SCW_DEFAULT_REGION          = "SCW_DEFAULT_REGION"
	SCW_DEFAULT_ZONE            = "SCW_DEFAULT_ZONE"
	SCW_COMMERCIAL_TYPE         = "SCW_COMMERCIAL_TYPE"
	SCW_IMAGE                   = "SCW_IMAGE"
	SCW_DISK_SIZE               = "SCW_DISK_SIZE"
	MACHINE_ID                  = "MACHINE_ID"
	MACHINE_FOLDER              = "MACHINE_FOLDER"
)

type Options struct {
	Image          string
	CommercialType string
	DiskSizeGB     string
	Zone           string
	OrganizationID string
	ProjectID      string
	ServerID       string

	MachineID     string
	MachineFolder string
}

func ConfigFromEnv() (Options, error) {
	return Options{
		Image:          os.Getenv(SCW_IMAGE),
		CommercialType: os.Getenv(SCW_COMMERCIAL_TYPE),
		Zone:           os.Getenv(SCW_DEFAULT_ZONE),
	}, nil
}

func FromEnv(init bool) (*Options, error) {
	retOptions := &Options{}

	var err error

	retOptions.Image, err = fromEnvOrError(SCW_IMAGE)
	if err != nil {
		return nil, err
	}
	retOptions.DiskSizeGB, err = fromEnvOrError(SCW_DISK_SIZE)
	if err != nil {
		return nil, err
	}
	retOptions.CommercialType, err = fromEnvOrError(SCW_COMMERCIAL_TYPE)
	if err != nil {
		return nil, err
	}

	retOptions.Zone, err = fromEnvOrError(SCW_DEFAULT_ZONE)
	if err != nil {
		return nil, err
	}
	retOptions.OrganizationID, err = fromEnvOrError(SCW_DEFAULT_ORGANIZATION_ID)
	if err != nil {
		return nil, err
	}
	retOptions.ProjectID, err = fromEnvOrError(SCW_DEFAULT_PROJECT_ID)
	if err != nil {
		return nil, err
	}

	// Return eraly if we're just doing init
	if init {
		return retOptions, nil
	}

	retOptions.MachineID, err = fromEnvOrError(MACHINE_ID)
	if err != nil {
		return nil, err
	}
	// prefix with devpod-
	retOptions.MachineID = "devpod-" + retOptions.MachineID

	retOptions.MachineFolder, err = fromEnvOrError(MACHINE_FOLDER)
	if err != nil {
		return nil, err
	}
	return retOptions, nil
}

func fromEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf(
			"couldn't find option %s in environment, please make sure %s is defined",
			name,
			name,
		)
	}

	return val, nil
}
