//go:build windows

package xwslgit

import (
	"os"
	"os/exec"
	"syscall"
)

func hasConsole() bool {
	hwnd := GetConsoleWindow()
	return hwnd != 0
}

// execCommand executes a command
func execCommand(args ...string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)
	if !hasConsole() {
		// sets `CREATE_NO_WINDOW` not to open console.
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
