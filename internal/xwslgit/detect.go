package xwslgit

import (
	"strings"

	"go.uber.org/zap"
)

// DetectDistribution detects WSL distribution to use
func (r *Runner) DetectDistribution(args ...string) string {
	if r.config.Detection.UseArguments {
		distro := r.detectDistributionFromArguments(args...)
		if distro != "" {
			return distro
		}
	}
	distro := detectDistributionFromPath(r.cwd)
	if distro == "" {
		return ""
	}
	r.logger.Debug(
		"detect distribution from current work directory",
		zap.Int("pid", r.pid),
		zap.String("cwd", r.cwd),
		zap.String("distro", distro),
	)
	return distro
}

func (r *Runner) detectDistributionFromArguments(args ...string) string {
	for _, arg := range args {
		distro := detectDistributionFromPath(arg)
		if distro != "" {
			r.logger.Debug(
				"detect distribution from argument",
				zap.Int("pid", r.pid),
				zap.String("arg", arg),
				zap.String("distro", distro),
			)
			return distro
		}
	}
	return ""
}

func detectDistributionFromPath(path string) string {
	lowerPath := strings.ToLower(path)
	for _, prefix := range wslPathPrefixes {
		if !strings.HasPrefix(lowerPath, prefix) {
			continue
		}
		rest := path[len(prefix):]
		split := strings.SplitN(rest, "\\", 2)
		return split[0]
	}
	return ""
}
