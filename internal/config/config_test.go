package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNonExistentFile(t *testing.T) {
	origDefaultPath := defaultPathFunc
	defaultPathFunc = func() string {
		return filepath.Join(t.TempDir(), "nonexistent.yml")
	}
	defer func() { defaultPathFunc = origDefaultPath }()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Projects == nil {
		t.Error("Expected Projects to be initialized, got nil")
	}
}

func TestLoadInvalidYaml(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.yml")
	os.WriteFile(tmpFile, []byte("invalid: yaml: content:"), 0644)

	origDefaultPath := defaultPathFunc
	defaultPathFunc = func() string { return tmpFile }
	defer func() { defaultPathFunc = origDefaultPath }()

	cfg, err := Load()
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
	// Config may be partially populated or empty on error - that's acceptable
	_ = cfg
}

func TestLoadValidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yml")
	validYaml := `projects:
  testproject:
    path: /tmp/test
    profiles:
      dev:
        tools:
          - name: vscode
            path: code
            args: ["."]
`
	os.WriteFile(tmpFile, []byte(validYaml), 0644)

	origDefaultPath := defaultPathFunc
	defaultPathFunc = func() string { return tmpFile }
	defer func() { defaultPathFunc = origDefaultPath }()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.Projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(cfg.Projects))
	}
	proj, ok := cfg.Projects["testproject"]
	if !ok {
		t.Error("Expected project 'testproject' to exist")
	}
	if proj.Path != "/tmp/test" {
		t.Errorf("Expected path '/tmp/test', got %q", proj.Path)
	}
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yml")

	origDefaultPath := defaultPathFunc
	defaultPathFunc = func() string { return tmpFile }
	defer func() { defaultPathFunc = origDefaultPath }()

	cfg := &Config{
		Projects: map[string]Project{
			"myproject": {
				Path: "/path/to/project",
				Profiles: map[string]Profile{
					"dev": {
						Tools: []Tool{
							{Name: "vscode", Path: "code"},
						},
					},
				},
			},
		},
	}

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}

	cfgloaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfgloaded.Projects) != 1 {
		t.Errorf("Expected 1 project after reload, got %d", len(cfgloaded.Projects))
	}
}
