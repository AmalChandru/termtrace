package share

import (
	"strings"
	"testing"
	"time"

	"github.com/AmalChandru/termtrace/internal/workflow"
)

func TestParseCustomRules_valid(t *testing.T) {
	t.Parallel()

	rules, err := parseCustomRules([]string{`(?i)token=\S+=token=<REDACTED>`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("got %d rules, want 1", len(rules))
	}
}

func TestParseCustomRules_invalidFormat(t *testing.T) {
	t.Parallel()

	_, err := parseCustomRules([]string{"token"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "pattern=replacement") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseCustomRules_invalidRegex(t *testing.T) {
	t.Parallel()

	_, err := parseCustomRules([]string{`(token=<REDACTED>`})
	if err == nil {
		t.Fatal("expected regex error, got nil")
	}
}

func TestDefaultRules_redactAuthBearer(t *testing.T) {
	t.Parallel()

	rules, err := defaultRules()
	if err != nil {
		t.Fatalf("defaultRules error: %v", err)
	}

	got := applyRules("Authorization: Bearer abc123", rules)
	want := "Authorization: Bearer <REDACTED>"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestDefaultRules_redactKeyValueSecrets(t *testing.T) {
	t.Parallel()

	rules, err := defaultRules()
	if err != nil {
		t.Fatalf("defaultRules error: %v", err)
	}

	got := applyRules("password=hunter2 token=abc api_key=xyz", rules)
	want := "password=<REDACTED> token=<REDACTED> api_key=<REDACTED>"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestApplyRedactions_allStepFields(t *testing.T) {
	t.Parallel()

	rules, err := defaultRules()
	if err != nil {
		t.Fatalf("defaultRules error: %v", err)
	}

	wf := &workflow.Workflow{
		Version: workflow.FormatVersion,
		Steps: []workflow.Step{
			{
				Command:   "curl -H 'Authorization: Bearer abc123' https://x",
				Stdout:    "token=secret123",
				Stderr:    "password=letmein",
				ExitCode:  0,
				Timestamp: time.Now().UTC(),
			},
		},
	}

	applyRedactions(wf, rules)

	step := wf.Steps[0]
	if strings.Contains(step.Command, "abc123") {
		t.Fatalf("command was not redacted: %q", step.Command)
	}
	if strings.Contains(step.Stdout, "secret123") {
		t.Fatalf("stdout was not redacted: %q", step.Stdout)
	}
	if strings.Contains(step.Stderr, "letmein") {
		t.Fatalf("stderr was not redacted: %q", step.Stderr)
	}
}
