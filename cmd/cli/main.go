package main

import (
	"fmt"
	"gomnirun/core/config"
	"gomnirun/core/executor"
	"log"
	"os"
	"strings"
)

const configFile = "config.json"

func main() {
	fmt.Println("Starting GomniRun in CLI Mode...")

	// Load configuration
	conf, err := config.Load(configFile)
	if err != nil {
		log.Println("No existing config found. Using defaults.")
		conf = config.Config{}
	}

	// Check if user provided variable updates
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg, "=", 2)
				key := parts[0]
				value := parts[1]

				// Update the variable if it exists in config
				if _, exists := conf.Variables[key]; exists {
					conf.Variables[key] = config.Variable{Type: conf.Variables[key].Type, Value: value}
					fmt.Printf("Updated %s -> %s\n", key, value)
				} else {
					fmt.Printf("Warning: Variable %s not found in config.\n", key)
				}
			}
		}

		// Save the updated config ONLY if overwrite is true
		if conf.Overwrite {
			config.Save(configFile, conf)
			fmt.Println("Configuration updated and saved.")
		} else {
			fmt.Println("Overwrite is disabled. Changes will not be saved.")
		}
	}

	// Generate final command with variables replaced
	finalCommand := executor.ReplacePlaceholders(conf.CommandTemplate, conf.Variables)
	fmt.Println("Executing command:", finalCommand)

	// Run the script
	output, err := executor.RunScript(conf.CommandTemplate, conf.Variables)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Println("Script Output:", output)
}
