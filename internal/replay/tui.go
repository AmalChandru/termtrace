package replay

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorDim    = "\033[2m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
)

// The stdout/stderr size above which replay shows a truncated view.
const maxReplayOutputBytes = 4096

func style(s, c string) string {
	return c + s + colorReset
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "<1s"
	}
	s := int(d.Round(time.Second).Seconds())
	if s < 60 {
		return fmt.Sprintf("%ds", s)
	}
	m := s / 60
	s %= 60
	if s == 0 {
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%dm %ds", m, s)
}

// Renders recorded DurationMs for the exit line (e.g. 42ms, 1.5s).
func formatStepElapsed(ms int64) string {
	if ms <= 0 {
		return ""
	}
	return (time.Duration(ms) * time.Millisecond).String()
}

// Splits s so head is at most limit bytes, preferring a break at the last newline
// within the prefix. If len(s) <= limit, tail is empty and truncated is false.
func splitTrunc(s string, limit int) (head, tail string, truncated bool) {
	if len(s) <= limit {
		return s, "", false
	}
	cut := limit
	for i := limit - 1; i >= 0; i-- {
		if s[i] == '\n' {
			cut = i + 1
			break
		}
	}
	if cut == 0 {
		cut = limit
	}
	return s[:cut], s[cut:], true
}

func ensureTrailingNewline(w io.Writer, s string) {
	if s == "" {
		return
	}
	if s[len(s)-1] != '\n' {
		fmt.Fprintln(w)
	}
}

// replayText prints content, truncating large streams. When truncated and not auto,
// prompts: type o then Enter to expand, or Enter to continue.
func replayText(w io.Writer, content string, auto bool, reader *bufio.Reader, paint func(string)) error {
	if content == "" {
		return nil
	}
	head, tail, cut := splitTrunc(content, maxReplayOutputBytes)
	paint(head)
	ensureTrailingNewline(w, head)

	if !cut {
		return nil
	}

	omit := len(tail)
	if auto {
		fmt.Fprintln(w, style(fmt.Sprintf("[... output truncated, %d bytes omitted]", omit), colorDim))
		return nil
	}

	fmt.Fprintln(w, style(fmt.Sprintf("[... output truncated, %d bytes omitted. Press o then Enter to expand, or Enter to continue]", omit), colorDim))

	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if strings.EqualFold(strings.TrimSpace(line), "o") {
		paint(tail)
		ensureTrailingNewline(w, tail)
	}
	return nil
}
