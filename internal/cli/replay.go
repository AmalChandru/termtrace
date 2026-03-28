package cli

import (
	"github.com/AmalChandru/termtrace/internal/app"
	"github.com/spf13/cobra"
)

func newReplayCmd() *cobra.Command {
	var auto bool
	var startStep int
	c := &cobra.Command{
		Use:   "replay <workflow.wf>",
		Short: "Replay a recorded workflow file",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return app.ReplayWorkflow(args[0], app.ReplayOptions{
				Auto:      auto,
				StartStep: startStep,
			})
		},
	}
	c.Flags().BoolVarP(&auto, "auto", "y", false, "replay all steps without waiting for Enter")
	c.Flags().IntVar(&startStep, "step", 1, "start replay from step number (1-based)")
	return c
}
