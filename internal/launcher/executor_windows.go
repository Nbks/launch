package launcher

import (
	"fmt"
	"launch/internal/config"
	"launch/internal/workspace"
	"os"
	"os/exec"
)

func expandVars(s string, vars map[string]string) string {
	result := s
	result = replaceVar(result, "%USERPROFILE%", os.Getenv("USERPROFILE"))
	result = replaceVar(result, "%PROJECT_PATH%", vars["PROJECT_PATH"])
	return result
}

func launchByType(t config.Tool, cmd *exec.Cmd) error {
	switch t.Type {
	case "terminal":
		// windows terminal abre una nueva ventana con el comando adentro
		termCmd := exec.Command("cmd.exe", append([]string{"/c", "start", "cmd.exe", "/k", cmd.Path}, cmd.Args[1:]...)...)
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
