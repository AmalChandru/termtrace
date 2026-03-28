// Package app wires domain packages for CLI and tests without Cobra.
package app

import (
	"github.com/AmalChandru/termtrace/internal/record"
	"github.com/AmalChandru/termtrace/internal/replay"
)

type ReplayOptions struct {
	Auto      bool
	StartStep int
}

func RecordSession(outputPath string) error {
	return record.Start(outputPath)
}

func StopRecording() error {
	return record.Stop()
}

func ReplayWorkflow(path string, opts ReplayOptions) error {
	return replay.Run(path, replay.Options{
		Auto:      opts.Auto,
		StartStep: opts.StartStep,
	})
}
