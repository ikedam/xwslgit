package xwslgit

import (
	"log"
	"os"
	"os/exec"

	"go.uber.org/zap"

	"github.com/pkg/errors"
)

var wslPathPrefixes = []string{
	`\\wsl$\`,
	`\\wsl.localhost\`,
}

// Config is configurations for XWSLGit
// Considered loading from YAML file.
type Config struct {
	// Debug configures debug outputs
	Debug DebugConfig
	// Detection configures how to detect distributions
	Detection DetectionConfig
}

// DebugConfig is configurations for debug outputs
type DebugConfig struct {
	// Enabled enables debug logging.
	Enabled bool
	// Logfile is the path to output logs.
	Logfile string
	// Envs contains environment variables to reacord.
	Envs []string
}

// DetectinoConfig is configurations about how to detect distribution
type DetectionConfig struct {
	UseArguments bool
}

// Runner runs operations for xwslgit
type Runner struct {
	config     *Config
	logger     *zap.Logger
	pid        int
	cwd        string
	executable string
}

// NewRunner creates a new runner
// This also records running environments.
func NewRunner(config *Config) (*Runner, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrapf(err, "could not get current working directory")
	}
	executable, err := os.Executable()
	if err != nil {
		return nil, errors.Wrapf(err, "could not get my executable path")
	}
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
		config:     config,
		logger:     logger,
		pid:        os.Getpid(),
		cwd:        cwd,
		executable: executable,
	}, nil
}

// Run executes xwslgit operations
func (r *Runner) Run(args ...string) int {
	var envs []string
	for _, envname := range r.config.Debug.Envs {
		envs = append(envs, envname+"="+os.Getenv(envname))
	}

	r.logger.Debug(
		"started",
		zap.Int("pid", r.pid),
		zap.Strings("args", args),
		zap.Strings("envs", envs),
		zap.String("cwd", r.cwd),
	)

	distro := r.DetectDistribution(args...)
	cmd := r.PrepareCommandForDistro(distro, args...)
	if cmd == nil {
		// The error message is output inside `PrepareCommandForDistro()`
		return 127
	}

	r.logger.Debug(
		"launch",
		zap.Int("pid", r.pid),
		zap.Strings("args", cmd.Args),
	)
	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		log.Printf("xwslgit: command %v failed with: %+v", cmd.Args, err)
		return 127
	}
	return 0
}
