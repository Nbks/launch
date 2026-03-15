package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var defaultPathFunc = DefaultPath

func SetPath(path string) {
	defaultPathFunc = func() string { return path }
}

type Config struct {
	Projects map[string]Project `yaml:"projects"`
}

type Project struct {
	Path     string             `yaml:"path"`
	Profiles map[string]Profile `yaml:"profiles"`
}

type Profile struct {
	Tools []Tool            `yaml:"tools"`
	Env   map[string]string `yaml:"env,omitempty"`
}

type Tool struct {
	Name    string   `yaml:"name"`
	Path    string   `yaml:"path"`
	Args    []string `yaml:"args,omitempty"`
	Type    string   `yaml:"type,omitempty"`    // "gui", "terminal", "background", "browser"
	Layout  string   `yaml:"layout,omitempty"`  // fullscreen, top-half, bottom-half, right-half, left-half
	Display string   `yaml:"display,omitempty"` // display 1 display 2 etc..
}

func Load() (*Config, error) {
	data, err := os.ReadFile(defaultPathFunc())
	if err != nil {
		return &Config{Projects: make(map[string]Project)}, nil
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(defaultPathFunc()), 0755)
	return os.WriteFile(defaultPathFunc(), data, 0644)
}
