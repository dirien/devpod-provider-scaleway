package cmd

import (
	"context"

	"github.com/dirien/devpod-provider-scaleway/pkg/scaleway"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"
)

// InitCmd holds the cmd flags
type InitCmd struct{}

// NewInitCmd defines a init
func NewInitCmd() *cobra.Command {
	cmd := &InitCmd{}
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init account",
		RunE: func(_ *cobra.Command, args []string) error {
			scalewayProvider, err := scaleway.NewProvider(log.Default, true)
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

	return initCmd
}

// Run runs the init logic
func (cmd *InitCmd) Run(
	ctx context.Context,
	scalewayProvider *scaleway.ScalewayProvider,
	logs log.Logger,
) error {
	return scaleway.Init(scalewayProvider)
}
