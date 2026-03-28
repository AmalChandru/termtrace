package cli

import (
	"github.com/spf13/cobra"
)

const version = "0.1.0-dev"

func Execute() error {
	return NewRootCmd().Execute()
}

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "termtrace",
		Short:   "Record and replay terminal workflows",
		Version: version,
	}

	root.AddCommand(
		newRecordCmd(),
		newStopCmd(),
		newReplayCmd(),
	)
	return root
}
