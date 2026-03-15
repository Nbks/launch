package launcher

import (
	"fmt"
	"launch/internal/config"
	"launch/internal/workspace"
	"os"
	"os/exec"
)

func getShell() string {
	return "/bin/bash"
}

func init() {
	if err := checkDisplayServer(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		os.Exit(1)
	}
}

func checkDisplayServer() error {
	session := os.Getenv("XDG_SESSION_TYPE")
	if session == "wayland" {
		return fmt.Errorf("Wayland isn't sopported yet, please login with X11 for use this app")
	}
	return nil
}

func expandVars(s string, vars map[string]string) string {
	result := s
	result = replaceVar(result, "$HOME", os.Getenv("HOME"))
	result = replaceVar(result, "$PROJECT_PATH", vars["PROJECT_PATH"])
	return result
}

func launchByType(t config.Tool, cmd *exec.Cmd) error {
	switch t.Type {
	case "terminal":
		termCmd := exec.Command("xterm", append([]string{"-e", cmd.Path}, cmd.Args[1:]...)...)
		termCmd.Dir = cmd.Dir
		termCmd.Env = cmd.Env
		if err := termCmd.Start(); err != nil {
			return fmt.Errorf("failed to launch terminal '%s': %w", t.Name, err)
		}
		workspace.ApplyLayout(t)

	case "background":
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to launch background '%s': %w", t.Name, err)
		}

	default: // gui
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to launch '%s': %w", t.Name, err)
		}

		workspace.ApplyLayout(t)
	}
	return nil
}
