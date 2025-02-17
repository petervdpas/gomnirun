package main

import (
	"flag"
	"log"
	"strings"

	"gomnirun/core/config"
	"gomnirun/cmd/cli"
	"gomnirun/cmd/fyne-ui"
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

	// Load configuration
	conf, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Run CLI or GUI directly
	if mode == "cli" {
		cli.RunCLI(&conf)
	} else {
		fyne_ui.RunGUI(&conf)
	}
}
