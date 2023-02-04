package xwslgit

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"go.uber.org/zap"

	"github.com/pkg/errors"
)

// DebugConfig is configurations for debug outputs
type DebugConfig struct {
	Enabled bool
	Logfile string
	Envs    []string
}

// Config is configurations for XWSLGit
// Considered loading from YAML file.
type Config struct {
	Debug DebugConfig
}

// Runner runs operations for xwslgit
type Runner struct {
	config *Config
	logger *zap.Logger
}

func NewRunner(config *Config) (*Runner, error) {
	zapConfig := zap.NewDevelopmentConfig()
	if config.Debug.Enabled && config.Debug.Logfile != "" {
		zapConfig.OutputPaths = []string{
			config.Debug.Logfile,
		}
	} else {
		zapConfig.OutputPaths = nil
	}
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, errors.Wrapf(err, "could not initialize logger")
	}
	return &Runner{
		config: config,
		logger: logger,
	}, nil
}

func (r *Runner) Run(args ...string) int {
	pid := os.Getpid()
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("xwslgit: could not get current working directory: %+v", err)
		return 127
	}

	var envs []string
	for _, envname := range r.config.Debug.Envs {
		envs = append(envs, envname+"="+os.Getenv(envname))
	}

	r.logger.Debug(
		"started",
		zap.Int("pid", pid),
		zap.Strings("args", args),
		zap.Strings("envs", envs),
		zap.String("cwd", cwd),
	)

	command, err := r.PrepareWindowsGit(os.Args[0], os.Args[1:]...)
	if err != nil {
		log.Printf("xwslgit: could not detect git on Windows: %+v", err)
		return 127
	}
	r.logger.Debug(
		"launch",
		zap.Int("pid", pid),
		zap.Strings("args", command),
	)
	cmd := exec.Command(command[0], command[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		log.Printf("xwslgit: command %v failed with: %+v", command, err)
		return 127
	}
	return 0
}

func findAnotherExecutable(currentExecutable, name string) (string, error) {
	currentFileInfo, err := os.Stat(currentExecutable)
	if err != nil {
		return "", errors.Wrapf(err, "could not stat %v", currentExecutable)
	}

	// `exec.LookPath()` may return myself, so search manually.
	_, currentExecutableName := filepath.Split(currentExecutable)
	var exts []string
	pathExt := os.Getenv("PATHEXT")
	if pathExt != "" {
		for _, e := range filepath.SplitList(pathExt) {
			if e == "" {
				continue
			}
			if e[0] != '.' {
				e = "." + e
			}
			exts = append(exts, e)
		}
	}
	if len(exts) == 0 {
		exts = []string{".exe"}
	}

	paths := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(paths) {
		// test whether not the place of xwslgit
		// (to skip possible .bat file in the same directory)
		possiblePath := filepath.Join(dir, currentExecutableName)
		possibleFileInfo, err := os.Stat(possiblePath)
		if err == nil && os.SameFile(currentFileInfo, possibleFileInfo) {
			continue
		}
		for _, e := range exts {
			path := filepath.Join(dir, "git"+e)
			fileInfo, err := os.Stat(path)
			if err != nil {
				continue
			}
			if os.SameFile(currentFileInfo, fileInfo) {
				continue
			}
			if fileInfo.IsDir() {
				continue
			}
			return path, nil
		}
	}
	return "", exec.ErrNotFound
}

func (r *Runner) PrepareWindowsGit(currentExecutable string, args ...string) ([]string, error) {
	gitPath, err := findAnotherExecutable(currentExecutable, "git")
	if err != nil {
		return nil, errors.Wrapf(err, "git on Windows was not found")
	}
	return append([]string{gitPath}, args...), nil
}
