package record

import "testing"

func TestCleanStepOutput_normalizeNewlines(t *testing.T) {
	t.Parallel()
	got := cleanStepOutput("a\r\nb\rc", "x")
	if got != "a\nb\nc" {
		t.Fatalf("got %q", got)
	}
}

func TestCleanStepOutput_stripBackspaces(t *testing.T) {
	t.Parallel()
	got := cleanStepOutput("ab\b\bc", "x")
	if got != "c" {
		t.Fatalf("got %q want c", got)
	}
}

func TestCleanStepOutput_stripANSI(t *testing.T) {
	t.Parallel()
	raw := "\x1b[31mred\x1b[0m\n"
	got := cleanStepOutput(raw, "x")
	if got != "red" {
		t.Fatalf("got %q", got)
	}
}

func TestCleanStepOutput_stripEchoedCommand(t *testing.T) {
	t.Parallel()
	got := cleanStepOutput("echo hi\nline\n", "echo hi")
	if got != "line" {
		t.Fatalf("got %q want line", got)
	}
}

func TestCleanStepOutput_stripTrailingPrompt(t *testing.T) {
	t.Parallel()
	got := cleanStepOutput("out\n$", "cmd")
	if got != "out" {
		t.Fatalf("got %q want out", got)
	}
}

func TestStripTrailingPromptLine_specialPrompts(t *testing.T) {
	t.Parallel()
	for _, p := range []string{"$", "#", "%", "❯"} {
		got := stripTrailingPromptLine("line\n" + p)
		if got != "line" {
			t.Fatalf("prompt %q: got %q", p, got)
		}
	}
}
