package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	// Define a flag for specifying a config file
	configFile := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	// Determine the mode (CLI or GUI)
	var mode string
	if flag.NArg() > 0 {
		mode = flag.Arg(0)
	} else {
		mode = "ui" // Default to GUI mode
	}

	runExecutable(mode, *configFile)
}

// runExecutable executes the correct binary depending on mode and OS
func runExecutable(mode string, configFile string) {
	var binary string
	if mode == "cli" {
		binary = "./builds/gomnirun-cli"
	} else {
		binary = "./builds/gomnirun-ui"
	}

	// On Windows, add .exe extension
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	// Ensure the binary exists
	if _, err := os.Stat(binary); os.IsNotExist(err) {
		log.Fatalf("Error: Binary %s not found. Please build the project first.", binary)
	}

	cmd := exec.Command(binary, "-config", configFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start %s with config %s: %v", binary, configFile, err)
	}
}
