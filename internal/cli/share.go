package cli

import (
	"github.com/AmalChandru/termtrace/internal/app"
	"github.com/spf13/cobra"
)

func newShareCmd() *cobra.Command {
	var output string
	var trimOutput int
	var redacts []string
	var noDefaultRedact bool

	c := &cobra.Command{
		Use:   "share <workflow.wf>",
		Short: "Create a sanitized copy of a workflow for sharing",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return app.ShareWorkflow(args[0], app.ShareOptions{
				Output:          output,
				TrimOutputBytes: trimOutput,
				RedactRules:     redacts,
				NoDefaultRedact: noDefaultRedact,
			})
		},
	}

	c.Flags().StringVarP(&output, "output", "o", "", "path to write cleaned workflow (default: <input>.shared.wf)")
	c.Flags().IntVar(&trimOutput, "trim-output", 0, "max bytes for stdout/stderr per step (0 = no trim)")
	c.Flags().StringArrayVar(&redacts, "redact", nil, "custom redaction rule in pattern=replacement form (repeatable)")
	c.Flags().BoolVar(&noDefaultRedact, "no-default-redact", false, "disable built-in redaction rules")

	return c
}
