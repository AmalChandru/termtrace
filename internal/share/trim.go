package share

import "github.com/AmalChandru/termtrace/internal/workflow"

func trimWorkflowOutput(wf *workflow.Workflow, max int) {
	if wf == nil || max <= 0 {
		return
	}
	for i := range wf.Steps {
		wf.Steps[i].Stdout = trimBytes(wf.Steps[i].Stdout, max)
		wf.Steps[i].Stderr = trimBytes(wf.Steps[i].Stderr, max)
	}
}

func trimBytes(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}
