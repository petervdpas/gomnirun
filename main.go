package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Define flags
	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Determine mode: CLI or GUI
	mode := "ui" // Default to GUI
	if len(flag.Args()) > 0 {
		arg := strings.ToLower(flag.Arg(0))
		if arg == "cli" || arg == "ui" {
			mode = arg
		}
	}

	// Build everything first if needed
	checkAndBuild("cli", "./builds/gomnirun-cli", "./cmd/cli/")
	checkAndBuild("ui", "./builds/gomnirun-ui", "./cmd/fyne-ui/")

	// Define which binary to execute
	var binaryPath string
	if mode == "cli" {
		binaryPath = "./builds/gomnirun-cli"
	} else {
		binaryPath = "./builds/gomnirun-ui"
	}

	// Construct command with config file flag
	cmdArgs := []string{}
	if *configFile != "config.json" {
		cmdArgs = append(cmdArgs, "-config", *configFile)
	}

	// Run the correct binary
	cmd := exec.Command(binaryPath, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start %s: %v", binaryPath, err)
	}
}

// checkAndBuild ensures binaries exist and builds them if missing
func checkAndBuild(mode string, binaryPath string, sourcePath string) {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		fmt.Printf("⚠️  Binary %s not found, building it now...\n", binaryPath)
		buildCmd := exec.Command("go", "build", "-o", binaryPath, sourcePath)
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		err := buildCmd.Run()
		if err != nil {
			log.Fatalf("❌ Failed to build %s: %v", mode, err)
		}
		fmt.Printf("✅ Successfully built %s\n", mode)
	}
}
