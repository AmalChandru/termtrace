package share

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AmalChandru/termtrace/internal/workflow"
)

type Options struct {
	Output          string
	TrimOutputBytes int
	RedactRules     []string
	NoDefaultRedact bool
}

func Run(inputPath string, opts Options) error {
	if strings.TrimSpace(inputPath) == "" {
		return fmt.Errorf("share: missing workflow file path")
	}

	wf, err := workflow.ReadFromFile(inputPath)
	if err != nil {
		return fmt.Errorf("share: read workflow: %w", err)
	}

	//1. Build rules (defaults + custom).
	rules, err := compileRules(opts.RedactRules, opts.NoDefaultRedact)
	if err != nil {
		return err
	}

	//2. Apply redaction to cmd/stdout/stderr.
	applyRedactions(wf, rules)

	//3. Output trimming (if asked).
	if opts.TrimOutputBytes > 0 {
		trimWorkflowOutput(wf, opts.TrimOutputBytes)
	}

	//4. Validate transformed workflow.
	if err := wf.Validate(); err != nil {
		return fmt.Errorf("share: invalid workflow after transform: %w", err)
	}

	//5. Resolve output path and write.
	out := opts.Output
	if strings.TrimSpace(out) == "" {
		out = defaultOutputPath(inputPath)
	}
	if err := workflow.WriteToFile(wf, out); err != nil {
		return fmt.Errorf("share: write workflow: %w", err)
	}

	return nil
}

func defaultOutputPath(inputPath string) string {
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(inputPath, ext)
	if ext == "" {
		return inputPath + ".shared.wf"
	}
	return base + ".shared" + ext
}
