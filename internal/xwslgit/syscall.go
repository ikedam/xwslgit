//go:build generate

package xwslgit

// To let go:generate generate a file for windows even on non-windows.

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output syscall_windows.go $GOFILE
//sys	GetConsoleWindow() (consoleWindow windows.HWND) = Kernel32.GetConsoleWindow
