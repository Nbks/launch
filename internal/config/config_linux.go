package config

import (
	"os"
	"path/filepath"
)

func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "launch", "config.yml")
}
