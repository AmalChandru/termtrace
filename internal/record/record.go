package record

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/AmalChandru/termtrace/internal/workflow"
)

// RunRecord runs an interactive shell on a PTY, records line-based steps, and writes path.
// Each line you type (terminated by Enter) is one Step.Command; Step.Stdout is PTY output
// until the next line or end of session.

// TODO: ExitCode is always 0 for now; shell does not expose per-command exit status.
// Shell integration parsing can help (?)
func RunRecord(outputPath string) error {
	if outputPath == "" {
		outputPath = "session.wf"
	}

	ptyFile, cmd, err := StartPTYShell()
	if err != nil {
		return fmt.Errorf("record: start shell: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		mu         sync.Mutex
		pendingOut bytes.Buffer // PTY bytes for current command's stdout window
	)

	// PTY → real stdout + capture buffer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 32*1024)
		for {
			n, readErr := ptyFile.Read(buf)
			if n > 0 {
				_, _ = os.Stdout.Write(buf[:n])
				mu.Lock()
				_, _ = pendingOut.Write(buf[:n])
				mu.Unlock()
			}
			if readErr != nil {
				if readErr != io.EOF {
					// PTY closed or error; stop stdin side
				}
				cancel()
				return
			}
		}
	}()

	var steps []workflow.Step
	var pendingCmd string
	var lastCmdStart time.Time
	firstLine := true
	appendStep := func(command, out string) {
		if strings.TrimSpace(command) == "" {
			return
		}
		var durMs int64
		if !lastCmdStart.IsZero() {
			durMs = time.Since(lastCmdStart).Milliseconds()
		}
		steps = append(steps, workflow.Step{
			Command:    command,
			Stdout:     out,
			Stderr:     "",
			ExitCode:   0,
			Timestamp:  time.Now().UTC(),
			DurationMs: durMs,
		})
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case sig := <-sigCh:
				if cmd.Process != nil {
					_ = cmd.Process.Signal(sig)
				}
				// Ensure the record loop can exit even if stdin read is blocked.
				_ = ptyFile.Close()
				cancel()
			}
		}
	}()

	stdin := os.Stdin
	lineBuf := make([]byte, 0, 256)
	readOne := make([]byte, 1)
	type stdinEvent struct {
		n   int
		b   byte
		err error
	}
	stdinCh := make(chan stdinEvent, 16)
	go func() {
		for {
			n, readErr := stdin.Read(readOne)
			ev := stdinEvent{n: n, err: readErr}
			if n == 1 {
				ev.b = readOne[0]
			}
			select {
			case <-ctx.Done():
				return
			case stdinCh <- ev:
			}
			if readErr != nil {
				return
			}
		}
	}()

recordLoop:
	for {
		select {
		case <-ctx.Done():
			break recordLoop
		case ev := <-stdinCh:
			if ev.n != 1 {
				if ev.err == io.EOF {
					break recordLoop
				}
				if ev.err != nil {
					return fmt.Errorf("record: stdin: %w", ev.err)
				}
				continue
			}
			b := ev.b
			if b == '\n' || b == '\r' {
				line := string(lineBuf)
				lineBuf = lineBuf[:0]
				if b == '\r' {
					// consume LF in CRLF
					_, _ = stdin.Read(readOne)
				}

				if firstLine {
					mu.Lock()
					pendingOut.Reset() // drop shell banner before first command
					mu.Unlock()
					firstLine = false
				} else {
					mu.Lock()
					out := pendingOut.String()
					pendingOut.Reset()
					mu.Unlock()
					appendStep(pendingCmd, cleanStepOutput(out, pendingCmd))
				}
				pendingCmd = line

				if _, werr := io.WriteString(ptyFile, line+"\n"); werr != nil {
					_ = ptyFile.Close()
					cancel()
					wg.Wait()
					return fmt.Errorf("record: write pty: %w", werr)
				}
				lastCmdStart = time.Now()
				continue
			}
			lineBuf = append(lineBuf, b)
		}
	}

	// Flush last line as a step.
	// E.g. user typed a command then Ctrl+D without newline.
	if pendingCmd != "" {
		mu.Lock()
		out := pendingOut.String()
		mu.Unlock()
		appendStep(pendingCmd, cleanStepOutput(out, pendingCmd))
	}

	_ = ptyFile.Close()
	cancel()
	wg.Wait()

	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	_ = cmd.Wait()

	if len(steps) == 0 {
		return fmt.Errorf("record: no commands recorded")
	}

	wf := &workflow.Workflow{
		Version: workflow.FormatVersion,
		Steps:   steps,
	}
	if err := workflow.WriteToFile(wf, outputPath); err != nil {
		return err
	}
	return nil
}

func Start(outputPath string) error {
	return RunRecord(outputPath)
}

func Stop() error {
	return fmt.Errorf("record: stop is not implemented yet")
}
