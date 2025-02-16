package executor

import (
	"fmt"
	"gomnirun/core/config"
	"os/exec"
)

// BashExecutor executes Bash scripts
type BashExecutor struct{}

func (b *BashExecutor) RunScript(commandTemplate string, variables map[string]config.Variable) (string, error) {
	scriptArgs := ReplacePlaceholders(commandTemplate, variables)
	cmd := exec.Command("bash", "-c", scriptArgs)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("bash error: %v\n%s", err, output)
	}
	return string(output), nil
}
