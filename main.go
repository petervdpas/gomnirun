package main

import (
	"flag"
	"fmt"
	"log"

	"gomnirun/cmd/cli"
	fyne_ui "gomnirun/cmd/fyne-ui"
	"gomnirun/core/config"
)

func main() {
	// Define command-line flags
	configFile := flag.String("config", "config.json", "Path to configuration file")
	mode := flag.String("mode", "ui", "Mode to run: 'cli' or 'ui'")
	help := flag.Bool("help", false, "Show usage information")

	// Parse flags first
	flag.Parse()

	// Display help message if requested
	if *help {
		fmt.Println(`
Usage: gomnirun [options]

Options:
  -config <file>   Specify configuration file (default: config.json)
  -mode <cli|ui>   Choose whether to run in CLI or GUI mode (default: ui)
  -help            Show this help message and exit

Examples:
  ./gomnirun -config=config-linux.json -mode=cli
  ./gomnirun -config=config-windows.json -mode=ui
		`)
		return
	}

	// Validate mode
	if *mode != "cli" && *mode != "ui" {
		log.Fatalf("❌ Invalid mode: '%s'. Use 'cli' or 'ui'.", *mode)
	}

	// Load the configuration file
	conf, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Run in the selected mode
	if *mode == "cli" {
		fmt.Printf("🔹 Running CLI mode with config file: %s\n", *configFile)
		cli.RunCLI(&conf, *configFile)
	} else {
		fmt.Printf("🖥️  Running GUI mode with config file: %s\n", *configFile)
		fyne_ui.RunGUI(&conf, *configFile)
	}
}
