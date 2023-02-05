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
		command, err := r.prepareWindowsGit(r.executable, os.Args[1:]...)
		if err != nil {
			log.Printf("xwslgit: could not detect git on Windows: %+v", err)
			return nil
		}
		return execCommand(command...)
	}
	command, err := r.prepareWSLGit(distro, os.Args[1:]...)
	if err != nil {
		log.Printf("xwslgit: could not prepare git on WSL %v: %+v", distro, err)
		return nil
	}
	return execCommand(command...)
}

func (r *Runner) prepareWindowsGit(currentExecutable string, args ...string) ([]string, error) {
	if r.config.WindowsGit.Path != "" {
		return append([]string{r.config.WindowsGit.Path}, args...), nil
	}
	gitPath, err := findAnotherExecutable(currentExecutable, "git")
	if err != nil {
		return nil, errors.Wrapf(err, "git on Windows was not found")
	}
	return append([]string{gitPath}, args...), nil
}

func (r *Runner) prepareWSLGit(distro string, args ...string) ([]string, error) {
	var config DistributionConfig
	if r.config.Distributions != nil {
		// Note: zero value is set if no configuration
		config = r.config.Distributions[distro]
	}
	replacedArgs := make([]string, len(args))
	for i, arg := range args {
		replacedArgs[i] = r.ConvertPathToWSL(distro, arg)
		if config.EscapeArguments {
			replacedArgs[i] = escapeArgument(replacedArgs[i])
		}
	}
	var command []string
	if len(config.Command) > 0 {
		command = make([]string, len(config.Command), len(config.Command)+len(replacedArgs))
		copy(command, config.Command)
	} else {
		// Not optimal (requires memory allocation for `append`), but easy to read
		command = []string{
			"wsl",
			"-d",
			distro,
			// not to have special letters escaped
			"--shell-type",
			"none",
			"--",
			"git",
		}
	}
	return append(command, replacedArgs...), nil
}

// ConvertPathToWSL convert a value to a WSL path if the value is a Windows path pointing WSL filesystem
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

var shellSpecialLetters = ` #$&*()\|[]{};'"<>?!` + "`"

func escapeArgument(arg string) string {
	if len(arg) == 0 {
		return `''`
	}
	if !strings.ContainsAny(arg, shellSpecialLetters) {
		return arg
	}
	var sb strings.Builder
	for _, letter := range arg {
		if strings.ContainsRune(shellSpecialLetters, letter) {
			sb.WriteRune('\\')
		}
		sb.WriteRune(letter)
	}
	return sb.String()
}
