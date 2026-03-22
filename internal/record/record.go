package record

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
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
	firstLine := true

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
			}
		}
	}()

	stdin := os.Stdin
	lineBuf := make([]byte, 0, 256)
	readOne := make([]byte, 1)

	for {
		n, readErr := stdin.Read(readOne)
		if n == 1 {
			b := readOne[0]
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
					steps = append(steps, workflow.Step{
						Command:   pendingCmd,
						Stdout:    out,
						Stderr:    "",
						ExitCode:  0,
						Timestamp: time.Now().UTC(),
					})
				}
				pendingCmd = line

				if _, werr := io.WriteString(ptyFile, line+"\n"); werr != nil {
					_ = ptyFile.Close()
					cancel()
					wg.Wait()
					return fmt.Errorf("record: write pty: %w", werr)
				}
				continue
			}
			lineBuf = append(lineBuf, b)
			continue
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return fmt.Errorf("record: stdin: %w", readErr)
		}
	}

	// Flush last line as a step.
	// E.g. user typed a command then Ctrl+D without newline.
	if pendingCmd != "" {
		mu.Lock()
		out := pendingOut.String()
		mu.Unlock()
		steps = append(steps, workflow.Step{
			Command:   pendingCmd,
			Stdout:    out,
			Stderr:    "",
			ExitCode:  0,
			Timestamp: time.Now().UTC(),
		})
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
