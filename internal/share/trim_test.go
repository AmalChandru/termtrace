package share

import (
	"testing"
	"time"

	"github.com/AmalChandru/termtrace/internal/workflow"
)

func TestTrimBytes_underLimit(t *testing.T) {
	t.Parallel()

	got := trimBytes("abc", 5)
	if got != "abc" {
		t.Fatalf("got %q want %q", got, "abc")
	}
}

func TestTrimBytes_exactLimit(t *testing.T) {
	t.Parallel()

	got := trimBytes("abcde", 5)
	if got != "abcde" {
		t.Fatalf("got %q want %q", got, "abcde")
	}
}

func TestTrimBytes_overLimit(t *testing.T) {
	t.Parallel()

	got := trimBytes("abcdef", 3)
	if got != "abc" {
		t.Fatalf("got %q want %q", got, "abc")
	}
}

func TestTrimWorkflowOutput_noopOnNil(t *testing.T) {
	t.Parallel()

	trimWorkflowOutput(nil, 10) // shouldn't panic
}

func TestTrimWorkflowOutput_noopWhenMaxNonPositive(t *testing.T) {
	t.Parallel()

	wf := &workflow.Workflow{
		Version: workflow.FormatVersion,
		Steps: []workflow.Step{
			{
				Command:   "echo hi",
				Stdout:    "hello",
				Stderr:    "err",
				ExitCode:  0,
				Timestamp: time.Now().UTC(),
			},
		},
	}

	trimWorkflowOutput(wf, 0)

	if wf.Steps[0].Stdout != "hello" || wf.Steps[0].Stderr != "err" {
		t.Fatalf("expected no trim when max <= 0, got stdout=%q stderr=%q", wf.Steps[0].Stdout, wf.Steps[0].Stderr)
	}
}
