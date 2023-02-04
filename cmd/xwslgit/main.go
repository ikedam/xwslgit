package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ikedam/xwslgit/internal/xwslgit"
	"gopkg.in/yaml.v3"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	// special command: xwslgitversion
	if len(os.Args) == 2 && os.Args[1] == "xwslgitversion" {
		fmt.Printf("xwslgitversion %v:%v\n", version, commit)
		os.Exit(0)
	}
	config := readConfig()
	runner, err := xwslgit.NewRunner(config)
	if err != nil {
		fmt.Printf("Unexpected error in xwslgit: %+v\n", err)
		os.Exit(127)
	}
	code := runner.Run(os.Args...)
	os.Exit(code)
}

// readConfig reads configuration from the file
// Proceeds with default values even if something get wrong.
func readConfig() *xwslgit.Config {
	var config xwslgit.Config
	configPath := filepath.Join(os.Getenv("APPDATA"), "xwslgit", "config.yaml")
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return &config
	}
	buf, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("ERROR: Could not read from %v: %+v", configPath, err)
		return &config
	}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		log.Printf("ERROR: Could not read from %v: %+v", configPath, err)
		return &config
	}
	if config.Debug.Logfile == "" {
		config.Debug.Logfile = filepath.Join(os.Getenv("APPDATA"), "xwslgit", "debug.log")
	}
	if config.Debug.Envs == nil {
		config.Debug.Envs = []string{
			"GIT_DIR",
			"GIT_WORK_TREE",
			"GIT_SSH",
		}
	}
	return &config
}
