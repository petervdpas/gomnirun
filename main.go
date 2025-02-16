package main

import (
	"flag"
	"fmt"
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
	}

	if mode == "cli" {
		fmt.Printf("🔹 Launching GomniRun CLI Mode with config: %s...\n", *configFile)
		runExecutable("cmd/cli/main.go", *configFile)
	} else {
		fmt.Printf("🔹 Launching GomniRun GUI Mode with config: %s...\n", *configFile)
		runExecutable("cmd/fyne-ui/main.go", *configFile)
	}
}

// runExecutable compiles and runs a Go file dynamically with a config file
func runExecutable(file string, configFile string) {
	cmd := exec.Command("go", "run", file, "-config", configFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start %s with config %s: %v", file, configFile, err)
	}
}
