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
		fmt.Printf("[%d/%d] $ %s\n", i+1, total, step.Command)
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
		fmt.Printf("exit code: %d\n", step.ExitCode)
		if !opts.Auto && i < total-1 {
			fmt.Print("Press Enter for next step...")
			if _, err := reader.ReadString('\n'); err != nil {
				return fmt.Errorf("replay: read input: %w", err)
			}
		}
		fmt.Println()
	}
	return nil
}
