package workflow_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/AmalChandru/termtrace/internal/workflow"
)

func TestRoundTrip(t *testing.T) {
	ts := time.Date(2025, 3, 21, 12, 0, 0, 0, time.UTC)
	w := &workflow.Workflow{
		Version: workflow.FormatVersion,
		Steps: []workflow.Step{
			{
				Command:   "echo hi",
				Stdout:    "hi\n",
				Stderr:    "",
				ExitCode:  0,
				Timestamp: ts,
			},
		},
	}

	data, err := workflow.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}

	got, err := workflow.Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != w.Version {
		t.Fatalf("version: got %q want %q", got.Version, w.Version)
	}
	if len(got.Steps) != 1 {
		t.Fatalf("steps: got %d want 1", len(got.Steps))
	}
	if !got.Steps[0].Timestamp.Equal(ts) {
		t.Fatalf("timestamp: got %v want %v", got.Steps[0].Timestamp, ts)
	}
}

func TestUnmarshal_invalidJSON(t *testing.T) {
	_, err := workflow.Unmarshal([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUnmarshal_unsupportedVersion(t *testing.T) {
	raw := `{
		"version": "99",
		"steps": [
			{
				"command": "true",
				"stdout": "",
				"stderr": "",
				"exit_code": 0,
				"timestamp": "2025-03-21T12:00:00Z"
			}
		]
	}`
	_, err := workflow.Unmarshal([]byte(raw))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWriteToFile_ReadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.wf")

	w := &workflow.Workflow{
		Version: workflow.FormatVersion,
		Steps: []workflow.Step{
			{
				Command:   "true",
				Stdout:    "",
				Stderr:    "",
				ExitCode:  0,
				Timestamp: time.Now().UTC().Truncate(time.Second),
			},
		},
	}
	if err := workflow.WriteToFile(w, path); err != nil {
		t.Fatal(err)
	}

	got, err := workflow.ReadFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Steps[0].Command != w.Steps[0].Command {
		t.Fatalf("command: got %q", got.Steps[0].Command)
	}
}
