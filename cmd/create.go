package cmd

import (
	"context"

	"github.com/dirien/devpod-provider-scaleway/pkg/scaleway"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"
)

// CreateCmd holds the cmd flags
type CreateCmd struct{}

// NewCreateCmd defines a command
func NewCreateCmd() *cobra.Command {
	cmd := &CreateCmd{}
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			scalewayProvider, err := scaleway.NewProvider(log.Default, false)
			if err != nil {
				return err
			}

			return cmd.Run(
				context.Background(),
				scalewayProvider,
				log.Default,
			)
		},
	}

	return createCmd
}

// Run runs the command logic
func (cmd *CreateCmd) Run(
	ctx context.Context,
	providerScaleway *scaleway.ScalewayProvider,
	logs log.Logger,
) error {
	return scaleway.Create(providerScaleway)
}
