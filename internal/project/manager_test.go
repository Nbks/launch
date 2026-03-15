package project

import (
	"os"
	"path/filepath"
	"testing"

	"launch/internal/config"
)

func setupTestConfig(t *testing.T) string {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yml")

	initialConfig := `projects: {}
`
	os.WriteFile(tmpFile, []byte(initialConfig), 0644)

	config.SetPath(tmpFile)

	return tmpFile
}

func TestAddNewProject(t *testing.T) {
	setupTestConfig(t)

	err := Add("newproject", "/path/to/newproject")
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	proj, exists := cfg.Projects["newproject"]
	if !exists {
		t.Error("Expected project 'newproject' to exist")
	}
	if proj.Path != "/path/to/newproject" {
		t.Errorf("Expected path '/path/to/newproject', got %q", proj.Path)
	}
	if _, exists := proj.Profiles["dev"]; !exists {
		t.Error("Expected default 'dev' profile to exist")
	}
}

func TestAddDuplicateProject(t *testing.T) {
	setupTestConfig(t)

	err := Add("existing", "/path/one")
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = Add("existing", "/path/two")
	if err == nil {
		t.Error("Expected error for duplicate project, got nil")
	}
}

func TestListEmpty(t *testing.T) {
	setupTestConfig(t)

	projects, err := List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(projects) != 0 {
		t.Errorf("Expected 0 projects, got %d", len(projects))
	}
}

func TestListWithProjects(t *testing.T) {
	tmpFile := setupTestConfig(t)

	existingConfig := `projects:
  project1:
    path: /path/one
    profiles:
      dev:
        tools:
          - name: vscode
            path: code
  project2:
    path: /path/two
    profiles:
      dev:
        tools:
          - name: vscode
            path: code
`
	os.WriteFile(tmpFile, []byte(existingConfig), 0644)

	projects, err := List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}

	names := make(map[string]bool)
	for _, p := range projects {
		names[p.Name] = true
	}
	if !names["project1"] || !names["project2"] {
		t.Error("Expected both project1 and project2 in results")
	}
}

func TestOpenProjectNotFound(t *testing.T) {
	setupTestConfig(t)

	err := Open("nonexistent", "dev")
	if err == nil {
		t.Error("Expected error for nonexistent project, got nil")
	}
}

func TestOpenProfileNotFound(t *testing.T) {
	tmpFile := setupTestConfig(t)

	existingConfig := `projects:
  myproject:
    path: /path/to/project
    profiles:
      dev:
        tools:
          - name: vscode
            path: code
`
	os.WriteFile(tmpFile, []byte(existingConfig), 0644)

	err := Open("myproject", "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent profile, got nil")
	}
}
