package project

import (
	"fmt"
	"launch/internal/config"
	"launch/internal/launcher"
)

type ProjectInfo struct {
	Name string
	config.Project
}

func Add(name, path string) error {
	cfg, err := config.Load()

	if err != nil {
		return err
	}

	if _, exists := cfg.Projects[name]; exists {
		return fmt.Errorf("project '%s' already exists", name)
	}

	// basic config
	newProject := config.Project{
		Path: path,
		Profiles: map[string]config.Profile{
			"dev": {
				Tools: []config.Tool{
					{
						Name: "vscode",
						Path: "code",
						Args: []string{"."},
					},
				},
			},
		},
	}

	cfg.Projects[name] = newProject

	return cfg.Save()

}

func Remove(name string) error {
	return fmt.Errorf("dont done yet")
}

func List() ([]ProjectInfo, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	var projects []ProjectInfo

	for name, proj := range cfg.Projects {
		projects = append(projects, ProjectInfo{
			Name:    name,
			Project: proj,
		})
	}

	return projects, nil
}

func Open(name, profileName string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	project, exists := cfg.Projects[name]
	if !exists {
		return fmt.Errorf("project '%s' not found", name)
	}

	profile, exists := project.Profiles[profileName]
	if !exists {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	return launcher.LaunchProfile(profile, project.Path)
}
