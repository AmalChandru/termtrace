package replay

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestSplitTrunc(t *testing.T) {
	t.Parallel()
	short := "hello\nworld"
	head, tail, cut := splitTrunc(short, 100)
	if cut || tail != "" || head != short {
		t.Fatalf("short string: head=%q tail=%q cut=%v", head, tail, cut)
	}

	long := strings.Repeat("a", 100) + "\n" + strings.Repeat("b", 100)
	head, tail, cut = splitTrunc(long, 50)
	if !cut {
		t.Fatal("expected truncation")
	}
	if len(head) > 50 {
		t.Fatalf("head len %d > 50", len(head))
	}
	if !strings.HasPrefix(long, head) {
		t.Fatal("head should be prefix of input")
	}
	if !strings.HasPrefix(long[len(head):], tail) {
		t.Fatal("tail should continue from head")
	}
}

func TestFormatDuration(t *testing.T) {
	t.Parallel()
	if got := formatDuration(500 * time.Millisecond); got != "<1s" {
		t.Fatalf("500ms: got %q", got)
	}
	if got := formatDuration(3 * time.Second); got != "3s" {
		t.Fatalf("3s: got %q", got)
	}
	if got := formatDuration(2 * time.Minute); got != "2m" {
		t.Fatalf("2m: got %q", got)
	}
	if got := formatDuration(2*time.Minute + 5*time.Second); got != "2m 5s" {
		t.Fatalf("2m5s: got %q", got)
	}
}

func TestFormatStepElapsed(t *testing.T) {
	t.Parallel()
	if formatStepElapsed(0) != "" || formatStepElapsed(-1) != "" {
		t.Fatal("non-positive should be empty")
	}
	if got := formatStepElapsed(42); got == "" {
		t.Fatal("expected non-empty for 42ms")
	}
}

func TestEnsureTrailingNewline(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	ensureTrailingNewline(&buf, "")
	if buf.Len() != 0 {
		t.Fatalf("empty: got %q", buf.String())
	}
	buf.Reset()
	ensureTrailingNewline(&buf, "hi\n")
	if buf.Len() != 0 {
		t.Fatalf("already newline: got %q", buf.String())
	}
	buf.Reset()
	ensureTrailingNewline(&buf, "hi")
	if buf.String() != "\n" {
		t.Fatalf("want single newline, got %q", buf.String())
	}
}

func TestReplayText_empty(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	r := bufio.NewReader(strings.NewReader(""))
	if err := replayText(&out, "", false, r, func(s string) { out.WriteString(s) }); err != nil {
		t.Fatal(err)
	}
}

func TestReplayText_shortNoPrompt(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	r := bufio.NewReader(strings.NewReader(""))
	content := "line1\nline2\n"
	if err := replayText(&out, content, false, r, func(s string) { out.WriteString(s) }); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "line1") {
		t.Fatalf("output: %q", out.String())
	}
}

func TestReplayText_truncAuto(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	r := bufio.NewReader(strings.NewReader(""))
	// Force truncation: content longer than maxReplayOutputBytes without early newline break
	content := strings.Repeat("x", maxReplayOutputBytes+50)
	if err := replayText(&out, content, true, r, func(s string) { out.WriteString(s) }); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "truncated") {
		t.Fatalf("expected truncated banner: %q", out.String())
	}
}

func TestReplayText_truncExpandO(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	prefix := strings.Repeat("a", maxReplayOutputBytes) + "\n"
	suffix := "TAIL_ONLY"
	content := prefix + suffix
	r := bufio.NewReader(strings.NewReader("o\n"))
	if err := replayText(&out, content, false, r, func(s string) { out.WriteString(s) }); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "TAIL_ONLY") {
		t.Fatalf("expected expanded tail: %q", out.String())
	}
}
