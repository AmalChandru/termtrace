package replay

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	startIdx := opts.StartStep - 1
	replayed := total - startIdx

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	reader := bufio.NewReader(os.Stdin)
	i := startIdx
stepLoop:
	for i < total {
		select {
		case <-ctx.Done():
			fmt.Fprintln(os.Stderr, "\nreplay interrupted")
			return nil
		default:
		}

		step := wf.Steps[i]
		stepHasErr := step.ExitCode != 0 || step.Stderr != ""
		cmdColor := colorCyan
		if stepHasErr {
			cmdColor = colorYellow
		}

		stepLabel := fmt.Sprintf("[%d/%d]", i+1, total)
		fmt.Printf("%s %s %s %s\n",
			style(stepLabel, colorDim),
			style(">", colorDim),
			style("$", cmdColor),
			style(step.Command, cmdColor),
		)

		if err := replayText(os.Stdout, step.Stdout, opts.Auto, reader, func(s string) { fmt.Print(s) }); err != nil {
			return fmt.Errorf("replay: stdin: %w", err)
		}
		if step.Stderr != "" {
			fmt.Fprintln(os.Stderr, style("[stderr]", colorDim))
			if err := replayText(os.Stderr, step.Stderr, opts.Auto, reader, func(s string) {
				fmt.Fprint(os.Stderr, style(s, colorRed))
			}); err != nil {
				return fmt.Errorf("replay: stdin: %w", err)
			}
		}

		exitText := fmt.Sprintf("[+] exit code: %d", step.ExitCode)
		if d := formatStepElapsed(step.DurationMs); d != "" {
			exitText = fmt.Sprintf("%s (%s)", exitText, d)
		}
		if step.ExitCode == 0 {
			fmt.Println(style(exitText, colorGreen))
		} else {
			fmt.Println(style(exitText, colorRed))
		}

		if opts.Auto {
			i++
			fmt.Println()
			continue
		}

		if i < total-1 {
		navPrompt:
			for {
				nextCmd := wf.Steps[i+1].Command
				fmt.Println(style(fmt.Sprintf("(next: %s)", nextCmd), colorDim))
				fmt.Print(style("Enter = next | b = previous | q = quit: ", colorDim))
				line, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("replay: read input: %w", err)
				}
				t := strings.TrimSpace(strings.ToLower(line))
				switch {
				case t == "q" || t == "quit":
					fmt.Fprintln(os.Stderr, "\nreplay quit")
					return nil
				case t == "b" || t == "back":
					if i > startIdx {
						i--
						fmt.Println()
						continue stepLoop
					}
					fmt.Println(style("Already at first step in this replay.", colorDim))
					continue navPrompt
				case t == "":
					i++
					fmt.Println()
					break navPrompt
				default:
					fmt.Println(style("Use Enter (next), b (previous), or q (quit).", colorDim))
					continue navPrompt
				}
			}
		} else {
			i++
			fmt.Println()
		}
	}

	failures := 0
	for j := startIdx; j < total; j++ {
		if wf.Steps[j].ExitCode != 0 {
			failures++
		}
	}

	// Prefer sum of per-step duration_ms (matches exit-line timings). Timestamp span excludes
	// the first step's window because Step.Timestamp is set when the step is finalized.
	// TODO: Get a better way for this.
	var recordedSpan time.Duration
	for j := startIdx; j < total; j++ {
		recordedSpan += time.Duration(wf.Steps[j].DurationMs) * time.Millisecond
	}
	if recordedSpan == 0 && replayed > 1 {
		recordedSpan = wf.Steps[total-1].Timestamp.Sub(wf.Steps[startIdx].Timestamp)
	}

	fmt.Println()
	fmt.Println(style(fmt.Sprintf("[+] Completed %d steps", replayed), colorGreen))
	fmt.Println(style(fmt.Sprintf("[*] Total time: %s", formatDuration(recordedSpan)), colorDim))
	if failures == 0 {
		fmt.Println(style("[F] Failures: 0", colorGreen))
	} else {
		fmt.Println(style(fmt.Sprintf("[F] Failures: %d", failures), colorRed))
	}

	return nil
}
