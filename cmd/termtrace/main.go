package main

import (
	"fmt"
	"os"

	"github.com/AmalChandru/termtrace/internal/record"
	"github.com/AmalChandru/termtrace/internal/replay"
	"github.com/spf13/cobra"
)

const version = "0.1.0-dev"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
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

func newRecordCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "record",
		Short: "Start recording a terminal session",
		RunE: func(cmd *cobra.Command, args []string) error {
			return record.Start()
		},
	}
}

func newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the current recording",
		RunE: func(cmd *cobra.Command, args []string) error {
			return record.Stop()
		},
	}
}

func newReplayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "replay [workflow.wf]",
		Short: "Replay a recorded workflow file",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := ""
			if len(args) > 0 {
				path = args[0]
			}
			return replay.Run(path)
		},
	}
}
