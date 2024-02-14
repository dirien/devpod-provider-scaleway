package cmd

import (
	"context"

	"github.com/dirien/devpod-provider-scaleway/pkg/scaleway"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"
)

// StopCmd holds the cmd flags
type StopCmd struct{}

// NewStopCmd defines a command
func NewStopCmd() *cobra.Command {
	cmd := &StopCmd{}
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop an instance",
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

	return stopCmd
}

// Run runs the command logic
func (cmd *StopCmd) Run(
	ctx context.Context,
	providerScaleway *scaleway.ScalewayProvider,
	logs log.Logger,
) error {
	return scaleway.Stop(providerScaleway)
}
