package executor

import (
	"fmt"
	"gomnirun/core/config"
	"os/exec"
	"runtime"
)

// PowerShellExecutor executes PowerShell scripts
type PowerShellExecutor struct{}

// RunScript properly formats and executes a PowerShell script
func (p *PowerShellExecutor) RunScript(commandTemplate string, variables map[string]config.Variable) (string, error) {
	scriptPath := variables["script"].Value
	args := ReplacePlaceholders(commandTemplate, variables)

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// ✅ FIX: Use -Command instead of -File for Windows
		cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", "& '"+scriptPath+"' "+args)
	} else {
		// 🔹 Keep Linux/macOS working as before
		cmd = exec.Command("pwsh", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	}

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("PowerShell error: %v\n%s", err, output)
	}
	return string(output), nil
}