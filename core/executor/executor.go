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
func GetExecutor(scriptPath string) (Executor, string, error) {
	ext := strings.ToLower(filepath.Ext(scriptPath))

	switch ext {
	case ".sh":
		return &BashExecutor{}, "Bash", nil
	case ".ps1":
		return &PowerShellExecutor{}, "PowerShell", nil
	case ".py":
		return &PythonExecutor{}, "Python", nil
	default:
		return nil, "", fmt.Errorf("unsupported script type: %s", ext)
	}
}

// RunScript dynamically selects the correct executor and runs the script
func RunScript(commandTemplate string, variables map[string]config.Variable) (string, error) {
	scriptPath := variables["script"].Value
	executor, executorName, err := GetExecutor(scriptPath)
	if err != nil {
		return "", err
	}

	// Log the selected executor *before* running the script
	fmt.Printf("🔹 Selected Executor: %s for script: %s\n", executorName, scriptPath)

	// Run the script
	return executor.RunScript(commandTemplate, variables)
}

// ReplacePlaceholders replaces placeholders in the command template
func ReplacePlaceholders(template string, variables map[string]config.Variable) string {
	for key, variable := range variables {
		value := quoteArguments(variable.Value) // Just quote all values

		// Replace placeholders
		template = strings.ReplaceAll(template, fmt.Sprintf("{%s}", key), value)
	}
	return template
}

// quoteArguments ensures all arguments containing spaces are properly quoted
func quoteArguments(value string) string {
	// Ensure arguments with spaces are correctly quoted
	if strings.Contains(value, " ") && !strings.HasPrefix(value, "\"") {
		return fmt.Sprintf("\"%s\"", value)
	}
	return value
}
