//go:build windows

package xwslgit

import (
	"os"
	"os/exec"
	"syscall"
)

// execCommand executes a command
func execCommand(args ...string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)
	// sets `CREATE_NO_WINDOW` not to open console.
	// build with -ldflags="-H=windowsgui"
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
