package xwslgit

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// PrepareCommandForDistro builds command and prepare exec.Cmd
// Errors are processed inside this function (e.g. printing an error message) and just returns nil.
func (r *Runner) PrepareCommandForDistro(distro string, args ...string) *exec.Cmd {
	if distro == "" {
		command, err := r.PrepareWindowsGit(os.Args[0], os.Args[1:]...)
		if err != nil {
			log.Printf("xwslgit: could not detect git on Windows: %+v", err)
			return nil
		}
		return execCommand(command...)
	}
	command, err := r.PrepareWSLGit(distro, os.Args[1:]...)
	if err != nil {
		log.Printf("xwslgit: could not prepare git on WSL %v: %+v", distro, err)
		return nil
	}
	return execCommand(command...)
}

func (r *Runner) PrepareWindowsGit(currentExecutable string, args ...string) ([]string, error) {
	gitPath, err := findAnotherExecutable(currentExecutable, "git")
	if err != nil {
		return nil, errors.Wrapf(err, "git on Windows was not found")
	}
	return append([]string{gitPath}, args...), nil
}

func (r *Runner) PrepareWSLGit(distro string, args ...string) ([]string, error) {
	replacedArgs := make([]string, len(args))
	for i, arg := range args {
		replacedArgs[i] = r.ConvertPathToWSL(distro, arg)
	}
	return append([]string{
		"wsl",
		"-d",
		distro,
		"--shell-type",
		"none",
		"--",
		"git",
	}, replacedArgs...), nil
}

func (r *Runner) ConvertPathToWSL(distro string, path string) string {
	lowerPath := strings.ToLower(path)
	for _, prefix := range wslPathPrefixes {
		if !strings.HasPrefix(lowerPath, prefix) {
			continue
		}
		rest := path[len(prefix):]
		split := strings.SplitN(rest, "\\", 2)
		if distro != split[0] {
			return "/mnt/wsl/" + split[0] + "/" + strings.ReplaceAll(split[1], "\\", "/")
		}
		return "/" + strings.ReplaceAll(split[1], "\\", "/")
	}
	// not considered a path. return as-is.
	return path
}
