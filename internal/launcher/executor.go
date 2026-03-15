package launcher

import (
	"os"
	"os/exec"
	"sync"

	"launch/internal/config"
	"launch/internal/logger"
)

func LaunchProfile(profile config.Profile, projectPath string) error {
	vars := map[string]string{"PROJECT_PATH": projectPath}

	logger.Log.Debug("launching profile",
		"projectPath", projectPath,
		"toolsCount", len(profile.Tools),
	)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, tool := range profile.Tools {
		wg.Add(1)

		go func(t config.Tool) {
			defer wg.Done()

			logger.Log.Debug("starting tool",
				"name", t.Name,
				"type", t.Type,
			)

			if err := launchTool(t, projectPath, vars, profile.Env); err != nil {
				logger.Log.Error("tool execution failed",
					"name", t.Name,
					"type", t.Type,
					"error", err,
				)

				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(tool)
	}

	wg.Wait()

	if len(errs) > 0 {
		logger.Log.Warn("profile finished with errors",
			"errorsCount", len(errs),
		)
		return joinErrors(errs)
	}

	logger.Log.Debug("profile finished successfully")
	return nil
}

func mapToEnv(m map[string]string) []string {
	env := make([]string, 0, len(m))
	for k, v := range m {
		env = append(env, k+"="+v)
	}
	return env
}

func launchTool(t config.Tool, projectPath string, vars map[string]string, env map[string]string) error {
	path := expandVars(t.Path, vars)

	args := make([]string, len(t.Args))
	for i, arg := range t.Args {
		args[i] = expandVars(arg, vars)
	}

	logger.Log.Debug("launching tool",
		"name", t.Name,
		"path", path,
		"args", args,
	)

	cmd := exec.Command(path, args...)
	cmd.Dir = projectPath
	cmd.Env = append(os.Environ(), mapToEnv(env)...)

	return launchByType(t, cmd)
}
