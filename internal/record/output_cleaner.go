package record

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	ansiCSI = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)
	ansiOSC = regexp.MustCompile(`\x1b\][^\x07]*(\x07|\x1b\\)`)
)

func cleanStepOutput(raw, command string) string {
	_, out := extractExitCodeAndCleanOutput(raw, command)
	return out
}

func extractExitCodeAndCleanOutput(raw, command string) (int, string) {
	s := normalizeNewlines(raw)
	s = stripBackspaces(s)
	s = ansiOSC.ReplaceAllString(s, "")
	s = ansiCSI.ReplaceAllString(s, "")
	s = stripControlChars(s)

	exitCode := 0
	s, exitCode = stripAndParseExitCodeMarkers(s)

	s = stripEchoedCommand(s, command)
	s = stripTrailingPromptLine(s)
	s = strings.TrimRight(s, "\n")
	return exitCode, s
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

// Full line marker format -> "__TT_RC__:<int>"
func stripAndParseExitCodeMarkers(s string) (string, int) {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	exitCode := 0

	for _, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, exitCodeMarkerPrefix) {
			raw := strings.TrimPrefix(t, exitCodeMarkerPrefix)
			if n, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil {
				exitCode = n // last valid marker wins
				continue
			}
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n"), exitCode
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
