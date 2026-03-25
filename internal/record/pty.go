package record

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/creack/pty"
)

// StartPTYShell starts a login shell on a PTY. Return the PTY master and the cmd.
func StartPTYShell() (*os.File, *exec.Cmd, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	shellName := strings.ToLower(filepath.Base(shell))

	// Start a clean interactive shell for recording.
	// User startup files can inject shell-specific completion code that breaks
	// in PTY capture contexts and pollutes command output.
	// TODO: Add a --use-shell-rc flag to opt into normal shell startup files.
	var cmd *exec.Cmd
	switch shellName {
	case "bash":
		cmd = exec.Command(shell, "--noprofile", "--norc", "-i")
	case "zsh":
		cmd = exec.Command(shell, "-f", "-i")
	case "fish":
		cmd = exec.Command(shell, "--no-config", "-i")
	default:
		// Fallback for POSIX shells that may not support -i as a standalone flag.
		cmd = exec.Command(shell)
	}
	cmd.Env = append(os.Environ(), "TERM=xterm-256color", "PS1=$ ")

	master, err := pty.Start(cmd)
	if err != nil {
		return nil, nil, err
	}
	return master, cmd, nil
}
