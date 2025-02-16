package executor

import (
	"fmt"
	"gomnirun/core/config"
	"os/exec"
)

// PythonExecutor executes Python scripts
type PythonExecutor struct{}

func (p *PythonExecutor) RunScript(commandTemplate string, variables map[string]config.Variable) (string, error) {
	scriptArgs := ReplacePlaceholders(commandTemplate, variables)
	cmd := exec.Command("python3", "-c", scriptArgs)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("python error: %v\n%s", err, output)
	}
	return string(output), nil
}
