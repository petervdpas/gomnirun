package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Define a flag for specifying a config file
	configFile := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	// Ensure that the CLI/GUI mode is correctly determined
	var mode string
	if flag.NArg() > 0 {
		mode = flag.Arg(0)
	} else {
		mode = "ui" // Default to UI mode
	}

	runExecutable(mode, *configFile)
}

// runExecutable compiles and runs a Go file dynamically with a config file
func runExecutable(mode string, configFile string) {
	var file string
	if mode == "cli" {
		file = "cmd/cli/main.go"
	} else {
		file = "cmd/fyne-ui/main.go"
	}

	// Ensure the correct order of arguments
	cmd := exec.Command("go", "run", file, mode, "-config", configFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start %s with config %s: %v", file, configFile, err)
	}
}
