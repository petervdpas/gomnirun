package executor

import (
	"fmt"
	"gomnirun/core/config"
	"os/exec"
	"runtime"
)

// PowerShellExecutor executes PowerShell scripts
type PowerShellExecutor struct{}

func (p *PowerShellExecutor) RunScript(commandTemplate string, variables map[string]config.Variable) (string, error) {
	scriptArgs := ReplacePlaceholders(commandTemplate, variables)

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Use Windows PowerShell
		cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", scriptArgs)
	} else {
		// Use PowerShell Core (`pwsh`) on Linux/macOS
		cmd = exec.Command("pwsh", "-ExecutionPolicy", "Bypass", "-File", scriptArgs)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("PowerShell error: %v\n%s", err, output)
	}
	return string(output), nil
}
