package record

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// StartPTYShell starts a login shell on a PTY. Return the PTY master and the cmd.
func StartPTYShell() (*os.File, *exec.Cmd, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	// TODO: login shell; change flags per shell later
	cmd := exec.Command(shell, "-l")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	master, err := pty.Start(cmd)
	if err != nil {
		return nil, nil, err
	}
	return master, cmd, nil
}
