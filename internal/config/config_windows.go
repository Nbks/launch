package config

import (
	"os"
	"path/filepath"
)

func DefaultPath() string {
	appData := os.Getenv("APPDATA")
	return filepath.Join(appData, "launch", "config.yml")
}
