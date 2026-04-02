package record

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/creack/pty"
)

const (
	// A machine-readble prompt marker.
	exitCodeMarkerPrefix = "__TT_RC__:"
	defaultPS1           = "$ "
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
	env := append(os.Environ(), "TERM=xterm-256color")

	switch shellName {
	case "bash":
		cmd = exec.Command(shell, "--noprofile", "--norc", "-i")
		// Somewhat a hacky way of getting exit codes,
		// make shell to print a machine-readable marker before each prompt
		// containing the previous command status.
		env = append(env,
			"PROMPT_COMMAND=printf '"+exitCodeMarkerPrefix+"%d\\n' $?",
			"PS1="+defaultPS1,
		)
	case "zsh":
		cmd = exec.Command(shell, "-f", "-i")
		// Define precmd hook to print previous command exit code for zsh.
		env = append(env,
			"ZDOTDIR=/nonexistent",
			"PS1="+defaultPS1,
			`ZSH_EXIT_CODE_HOOK=precmd(){ printf '`+exitCodeMarkerPrefix+`%d\n' $?; }`,
		)
	default:
		// Fallback for POSIX shells that may not support -i as a standalone flag.
		cmd = exec.Command(shell)
		env = append(env, "PS1="+defaultPS1)
	}

	cmd.Env = env

	master, err := pty.Start(cmd)
	if err != nil {
		return nil, nil, err
	}
	// Best-effort: disable echo on the PTY slave side to avoid duplicated
	// command input in captured/live output
	// https://github.com/creack/pty/issues/204
	_ = disablePTYEcho(master)
	return master, cmd, nil
}
