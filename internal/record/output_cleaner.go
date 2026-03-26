package record

import (
	"regexp"
	"strings"
)

var (
	ansiCSI = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)
	ansiOSC = regexp.MustCompile(`\x1b\][^\x07]*(\x07|\x1b\\)`)
)

func cleanStepOutput(raw, command string) string {
	s := normalizeNewlines(raw)
	s = stripBackspaces(s)
	s = ansiOSC.ReplaceAllString(s, "")
	s = ansiCSI.ReplaceAllString(s, "")
	s = stripControlChars(s)

	s = stripEchoedCommand(s, command)
	s = stripTrailingPromptLine(s)

	return strings.TrimRight(s, "\n")
}

func normalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func stripBackspaces(s string) string {
	r := make([]rune, 0, len(s))
	for _, ch := range s {
		if ch == '\b' {
			if len(r) > 0 {
				r = r[:len(r)-1]
			}
			continue
		}
		r = append(r, ch)
	}
	return string(r)
}

func stripControlChars(s string) string {
	var b strings.Builder
	for _, ch := range s {
		if ch == '\n' || ch == '\t' || ch >= 0x20 {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func stripEchoedCommand(out, cmd string) string {
	out = strings.TrimLeft(out, "\n")
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return out
	}
	if strings.HasPrefix(out, cmd+"\n") {
		return strings.TrimPrefix(out, cmd+"\n")
	}
	return out
}

func stripTrailingPromptLine(out string) string {
	lines := strings.Split(out, "\n")
	if len(lines) == 0 {
		return out
	}
	last := strings.TrimSpace(lines[len(lines)-1])
	switch last {
	case "$", "#", "%", "❯":
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}
