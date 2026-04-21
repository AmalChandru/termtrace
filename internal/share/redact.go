package share

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AmalChandru/termtrace/internal/workflow"
)

type redactRule struct {
	re          *regexp.Regexp
	replacement string
}

func compileRules(custom []string, noDefault bool) ([]redactRule, error) {
	var rules []redactRule

	if !noDefault {
		def, err := defaultRules()
		if err != nil {
			return nil, err
		}
		rules = append(rules, def...)
	}

	userRules, err := parseCustomRules(custom)
	if err != nil {
		return nil, err
	}
	rules = append(rules, userRules...)

	return rules, nil
}

func defaultRules() ([]redactRule, error) {
	// Redacts "bare minimum" sensitive values by:
	// Replaces OAuth-style "Authorization: Bearer <token>" with "Authorization: Bearer <REDACTED>"
	// Replaces "<secret-key>=<value>" pairs (password/passwd/pwd/token/api[_-]?key/secret) with "<key>=<REDACTED>"
	// matching is case-insensitive; values are treated as non-whitespace sequences.
	pairs := [][2]string{
		{`(?i)(authorization:\s*bearer\s+)\S+`, `${1}<REDACTED>`},
		{`(?i)\b(password|passwd|pwd|token|api[_-]?key|secret)\s*=\s*\S+`, `$1=<REDACTED>`},
	}

	out := make([]redactRule, 0, len(pairs))
	for _, p := range pairs {
		re, err := regexp.Compile(p[0])
		if err != nil {
			return nil, fmt.Errorf("share: invalid default redact pattern %q: %w", p[0], err)
		}
		out = append(out, redactRule{re: re, replacement: p[1]})
	}

	return out, nil
}

func parseCustomRules(vals []string) ([]redactRule, error) {
	out := make([]redactRule, 0, len(vals))
	for _, v := range vals {
		i := strings.Index(v, "=")
		if i <= 0 {
			return nil, fmt.Errorf("share: invalid --redact %q (want pattern=replacement)", v)
		}
		pattern := v[:i]
		repl := v[i+1:]

		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("share: invalid --redact pattern %q: %w", pattern, err)
		}
		out = append(out, redactRule{re: re, replacement: repl})
	}
	return out, nil
}

func applyRedactions(wf *workflow.Workflow, rules []redactRule) {
	if wf == nil || len(rules) == 0 {
		return
	}
	for i := range wf.Steps {
		wf.Steps[i].Command = applyRules(wf.Steps[i].Command, rules)
		wf.Steps[i].Stdout = applyRules(wf.Steps[i].Stdout, rules)
		wf.Steps[i].Stderr = applyRules(wf.Steps[i].Stderr, rules)
	}
}

func applyRules(s string, rules []redactRule) string {
	out := s
	for _, r := range rules {
		out = r.re.ReplaceAllString(out, r.replacement)
	}
	return out
}
