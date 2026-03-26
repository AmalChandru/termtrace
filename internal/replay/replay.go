package replay

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AmalChandru/termtrace/internal/workflow"
)

type Options struct {
	Auto      bool
	StartStep int
}

const (
	colorReset = "\033[0m"
	colorDim   = "\033[2m"
	colorCyan  = "\033[36m"
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
)

func style(s, c string) string {
	return c + s + colorReset
}

func Run(path string, opts Options) error {

	if path == "" {
		return fmt.Errorf("replay: missing workflow file path")
	}

	wf, err := workflow.ReadFromFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("replay: file not found: %s", path)
		}
		return fmt.Errorf("replay: %w", err)
	}

	if opts.StartStep < 1 {
		opts.StartStep = 1
	}
	total := len(wf.Steps)
	if opts.StartStep > total {
		return fmt.Errorf("replay: --step out of range: %d (workflow has %d steps)", opts.StartStep, total)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	reader := bufio.NewReader(os.Stdin)
	for i := opts.StartStep - 1; i < total; i++ {
		select {
		case <-ctx.Done():
			fmt.Fprintln(os.Stderr, "\nreplay interrupted")
			return nil
		default:
		}

		step := wf.Steps[i]

		stepLabel := fmt.Sprintf("[%d/%d]", i+1, total)
		fmt.Printf("%s %s %s\n",
			style(stepLabel, colorDim),
			style("$", colorCyan),
			style(step.Command, colorCyan),
		)

		if step.Stdout != "" {
			fmt.Print(step.Stdout)
			if step.Stdout[len(step.Stdout)-1] != '\n' {
				fmt.Println()
			}
		}
		if step.Stderr != "" {
			fmt.Fprint(os.Stderr, step.Stderr)
			if step.Stderr[len(step.Stderr)-1] != '\n' {
				fmt.Fprintln(os.Stderr)
			}
		}

		exitText := fmt.Sprintf("exit code: %d", step.ExitCode)
		if step.ExitCode == 0 {
			fmt.Println(style(exitText, colorGreen))
		} else {
			fmt.Println(style(exitText, colorRed))
		}

		if !opts.Auto && i < total-1 {
			fmt.Print(style("Press Enter for next step...", colorDim))
			if _, err := reader.ReadString('\n'); err != nil {
				return fmt.Errorf("replay: read input: %w", err)
			}
		}
		fmt.Println()
	}
	return nil
}
