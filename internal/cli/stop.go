package cli

import (
	"github.com/AmalChandru/termtrace/internal/app"
	"github.com/spf13/cobra"
)

func newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the current recording",
		RunE: func(_ *cobra.Command, _ []string) error {
			return app.StopRecording()
		},
	}
}
