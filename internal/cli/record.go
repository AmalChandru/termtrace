package cli

import (
	"github.com/AmalChandru/termtrace/internal/app"
	"github.com/spf13/cobra"
)

func newRecordCmd() *cobra.Command {
	var output string
	c := &cobra.Command{
		Use:   "record",
		Short: "Start recording a terminal session",
		RunE: func(_ *cobra.Command, _ []string) error {
			return app.RecordSession(output)
		},
	}
	c.Flags().StringVarP(&output, "output", "o", "session.wf", "path to write the workflow file")
	return c
}
