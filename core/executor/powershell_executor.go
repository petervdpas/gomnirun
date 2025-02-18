package executor

import (
	"fmt"
	"gomnirun/core/config"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// PowerShellExecutor executes PowerShell scripts
type PowerShellExecutor struct{}

// RunScript executes a PowerShell script with correctly formatted arguments
func (p *PowerShellExecutor) RunScript(commandTemplate string, variables map[string]config.Variable) (string, error) {
	// Replace placeholders in the command template
	scriptCommand := ReplacePlaceholders(commandTemplate, variables)

	// Split script and arguments
	args := strings.Fields(scriptCommand) // ["pwsh", "./Convert.ps1", "-arg1", "value1"]
	if len(args) < 2 {
		return "", fmt.Errorf("invalid PowerShell command format, expected at least a script path")
	}

	// Get script path (must be absolute)
	scriptPath, err := filepath.Abs(args[1])
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for script: %v", err)
	}

	// Extract arguments (excluding script path)
	scriptArgs := args[2:]

	// Choose PowerShell version based on OS
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	} else {
		cmd = exec.Command("pwsh", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	}

	// Append arguments separately (keys without quotes, values with quotes)
	for i := 0; i < len(scriptArgs); i += 2 {
		key := scriptArgs[i]
		if i+1 < len(scriptArgs) {
			value := scriptArgs[i+1]
			// Wrap only values with spaces in quotes
			if strings.Contains(value, " ") {
				value = fmt.Sprintf(`"%s"`, value)
			}
			cmd.Args = append(cmd.Args, key, value)
		} else {
			// If there's a key without a value, just add it
			cmd.Args = append(cmd.Args, key)
		}
	}

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("PowerShell error: %v\n%s", err, output)
	}

	return string(output), nil
}
