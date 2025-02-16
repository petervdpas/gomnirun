package executor

import (
	"fmt"
	"gomnirun/core/config"
	"os/exec"
)

// Replace placeholders in the command template
func ReplacePlaceholders(template string, variables map[string]config.Variable) []string {
	args := []string{}

	for key, variable := range variables {
		value := variable.Value

		// Wrap value in quotes only when needed for Bash
		if variable.Type == "string" || variable.Type == "file" {
			args = append(args, fmt.Sprintf("-%s=%s", key, value))
		} else {
			args = append(args, fmt.Sprintf("-%s=%s", key, value))
		}
	}
	return args
}

// RunScript executes the script with properly formatted arguments
func RunScript(commandTemplate string, variables map[string]config.Variable) (string, error) {
	scriptArgs := ReplacePlaceholders(commandTemplate, variables)

	// Command execution: Pass script name and arguments separately
	cmd := exec.Command("bash", append([]string{"./test_script.sh"}, scriptArgs...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running script: %v\n%s", err, output)
	}
	return string(output), nil
}
