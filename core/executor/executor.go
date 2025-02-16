package executor

import (
	"fmt"
	"gomnirun/core/config"
	"os/exec"
	"strings"
)

// Replace placeholders in the command template
func ReplacePlaceholders(commandTemplate string, variables map[string]config.Variable) string {
	for key, variable := range variables {
		placeholder := fmt.Sprintf("{%s}", key) // e.g., {script}, {var1}
		commandTemplate = strings.ReplaceAll(commandTemplate, placeholder, variable.Value)
	}
	return commandTemplate
}

// RunScript replaces variables in the command and executes it
func RunScript(commandTemplate string, variables map[string]config.Variable) (string, error) {
	finalCommand := ReplacePlaceholders(commandTemplate, variables)

	// Split command into parts for exec.Command
	cmdParts := strings.Fields(finalCommand)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running script: %v", err)
	}
	return string(output), nil
}
