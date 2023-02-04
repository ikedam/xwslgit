package xwslgit

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

// findAnotherExecutable searches paths for specified executable, but not currentExecutable
// implemented just for Windows
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
