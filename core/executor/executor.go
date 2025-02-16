package executor

import (
	"fmt"
	"gomnirun/core/config"
	"path/filepath"
	"strings"
)

// Executor interface for different script execution types
type Executor interface {
	RunScript(commandTemplate string, variables map[string]config.Variable) (string, error)
}

// GetExecutor dynamically returns the correct executor based on file extension
func GetExecutor(scriptPath string) (Executor, error) {
	ext := strings.ToLower(filepath.Ext(scriptPath))
	switch ext {
	case ".sh":
		return &BashExecutor{}, nil
	case ".ps1":
		return &PowerShellExecutor{}, nil
	case ".py":
		return &PythonExecutor{}, nil
	default:
		return nil, fmt.Errorf("unsupported script type: %s", ext)
	}
}

// RunScript dynamically selects the correct executor and runs the script
func RunScript(commandTemplate string, variables map[string]config.Variable) (string, error) {
	scriptPath := variables["script"].Value
	executor, err := GetExecutor(scriptPath)
	if err != nil {
		return "", err
	}
	return executor.RunScript(commandTemplate, variables)
}

// ReplacePlaceholders replaces placeholders in the command template
func ReplacePlaceholders(template string, variables map[string]config.Variable) string {
	for key, variable := range variables {
		value := variable.Value

		// Wrap value in quotes if it contains spaces
		if strings.Contains(value, " ") {
			value = fmt.Sprintf("\"%s\"", value)
		}

		// Replace placeholders
		template = strings.ReplaceAll(template, fmt.Sprintf("{%s}", key), value)
	}
	return template
}
