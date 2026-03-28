package workflow_test

import (
	"testing"
	"time"

	"github.com/AmalChandru/termtrace/internal/workflow"
)

func TestValidate_nil(t *testing.T) {
	t.Parallel()
	var w *workflow.Workflow
	if err := w.Validate(); err == nil {
		t.Fatal("expected error for nil")
	}
}

func TestValidate_emptyVersion(t *testing.T) {
	t.Parallel()
	w := &workflow.Workflow{Version: "", Steps: []workflow.Step{{Command: "x", Timestamp: time.Now().UTC()}}}
	if err := w.Validate(); err == nil {
		t.Fatal("expected error for empty version")
	}
}

func TestValidate_badVersion(t *testing.T) {
	t.Parallel()
	w := &workflow.Workflow{
		Version: "0",
		Steps:   []workflow.Step{{Command: "x", Timestamp: time.Now().UTC()}},
	}
	if err := w.Validate(); err == nil {
		t.Fatal("expected error for wrong version")
	}
}

func TestValidate_noSteps(t *testing.T) {
	t.Parallel()
	w := &workflow.Workflow{Version: workflow.FormatVersion, Steps: nil}
	if err := w.Validate(); err == nil {
		t.Fatal("expected error for no steps")
	}
}

func TestValidate_emptyCommand(t *testing.T) {
	t.Parallel()
	w := &workflow.Workflow{
		Version: workflow.FormatVersion,
		Steps: []workflow.Step{
			{Command: "", Timestamp: time.Now().UTC()},
		},
	}
	if err := w.Validate(); err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestValidate_zeroTimestamp(t *testing.T) {
	t.Parallel()
	w := &workflow.Workflow{
		Version: workflow.FormatVersion,
		Steps: []workflow.Step{
			{Command: "true", Timestamp: time.Time{}},
		},
	}
	if err := w.Validate(); err == nil {
		t.Fatal("expected error for zero timestamp")
	}
}

func TestValidate_ok(t *testing.T) {
	t.Parallel()
	w := &workflow.Workflow{
		Version: workflow.FormatVersion,
		Steps: []workflow.Step{
			{Command: "true", Timestamp: time.Now().UTC()},
		},
	}
	if err := w.Validate(); err != nil {
		t.Fatal(err)
	}
}
