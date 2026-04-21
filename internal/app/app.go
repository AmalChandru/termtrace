// Package app wires domain packages for CLI and tests without Cobra.
package app

import (
	"github.com/AmalChandru/termtrace/internal/record"
	"github.com/AmalChandru/termtrace/internal/replay"
	"github.com/AmalChandru/termtrace/internal/share"
)

type ReplayOptions struct {
	Auto      bool
	StartStep int
}

type ShareOptions struct {
	Output          string
	TrimOutputBytes int
	RedactRules     []string
	NoDefaultRedact bool
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

func ShareWorkflow(path string, opts ShareOptions) error {
	return share.Run(path, share.Options{
		Output:          opts.Output,
		TrimOutputBytes: opts.TrimOutputBytes,
		RedactRules:     opts.RedactRules,
		NoDefaultRedact: opts.NoDefaultRedact,
	})
}
