//go:build !windows

package xwslgit

// This file is actually never be used.
// Just for developing convenience.

import (
	"os"
	"os/exec"
)

func execCommand(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
