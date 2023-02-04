//go:build !windows

package xwslgit

// This file is actually never be used.
// Just for developing convenience.

import (
	"os"
	"os/exec"
)

func execCommand(args ...string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
