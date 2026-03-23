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
	var output string
	c := &cobra.Command{
		Use:   "record",
		Short: "Start recording a terminal session",
		RunE: func(_ *cobra.Command, _ []string) error {
			return record.Start(output)
		},
	}
	c.Flags().StringVarP(&output, "output", "o", "session.wf", "path to write the workflow file")
	return c
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
	var auto bool
	var startStep int
	c := &cobra.Command{
		Use:   "replay <workflow.wf>",
		Short: "Replay a recorded workflow file",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return replay.Run(args[0], replay.Options{
				Auto:      auto,
				StartStep: startStep,
			})
		},
	}
	c.Flags().BoolVarP(&auto, "auto", "y", false, "replay all steps without waiting for Enter")
	c.Flags().IntVar(&startStep, "step", 1, "start replay from step number (1-based)")
	return c
}
