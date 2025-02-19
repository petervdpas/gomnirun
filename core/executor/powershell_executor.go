package executor

import (
	"fmt"
	"gomnirun/core/config"
	"os/exec"
	"runtime"
	"strings"
)

// PowerShellExecutor executes PowerShell scripts
type PowerShellExecutor struct{}

// RunScript properly formats and executes a PowerShell script
func (p *PowerShellExecutor) RunScript(commandTemplate string, variables map[string]config.Variable) (string, error) {
	scriptPath := variables["script"].Value
	args := ReplacePlaceholders(commandTemplate, variables)

	// Ensure script path is properly quoted if it contains spaces
	if strings.Contains(scriptPath, " ") && !strings.HasPrefix(scriptPath, "\"") {
		scriptPath = fmt.Sprintf("\"%s\"", scriptPath)
	}

	// Ensure all arguments are properly quoted to handle spaces
	if strings.Contains(args, " ") && !strings.HasPrefix(args, "\"") {
		args = fmt.Sprintf("\"%s\"", args)
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// ✅ Windows: Use `-Command "& 'script.ps1' args"` to handle spaces correctly
		fullCommand := fmt.Sprintf("& %s %s", scriptPath, args)
		cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", fullCommand)
	} else {
		// ✅ Linux/macOS: Use `-Command "& 'script.ps1' args"` (same fix as Windows)
		fullCommand := fmt.Sprintf("& %s %s", scriptPath, args)
		cmd = exec.Command("pwsh", "-ExecutionPolicy", "Bypass", "-Command", fullCommand)
	}

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("PowerShell error: %v\n%s", err, output)
	}
	return string(output), nil
}
