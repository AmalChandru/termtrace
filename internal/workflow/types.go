package workflow

import (
	"errors"
	"fmt"
	"time"
)

// FormatVersion is the only workflow format version this build understands.
const FormatVersion = "1"

// Workflow is the on-disk root object for a .wf file.
type Workflow struct {
	Version string `json:"version"`
	Steps   []Step `json:"steps"`
}

// Step is one captured command and its result.
type Step struct {
	Command    string    `json:"command"`
	Stdout     string    `json:"stdout"`
	Stderr     string    `json:"stderr"`
	ExitCode   int       `json:"exit_code"`
	Timestamp  time.Time `json:"timestamp"`
	// DurationMs is wall time from sending this command to the PTY until the next line
	// is submitted or the session ends (includes output, reading, and typing—not shell CPU time).
	DurationMs int64 `json:"duration_ms,omitempty"`
}

func (w *Workflow) Validate() error {
	if w == nil {
		return errors.New("workflow: nil workflow")
	}
	if w.Version == "" {
		return errors.New("workflow: version is required")
	}
	if w.Version != FormatVersion {
		return fmt.Errorf("workflow: unsupported version %q (want %q)", w.Version, FormatVersion)
	}
	if len(w.Steps) == 0 {
		return errors.New("workflow: steps must be non-empty")
	}
	for i, s := range w.Steps {
		if s.Command == "" {
			return fmt.Errorf("workflow: step %d: command is required", i)
		}
		if s.Timestamp.IsZero() {
			return fmt.Errorf("workflow: step %d: timestamp is required", i)
		}
	}
	return nil
}
