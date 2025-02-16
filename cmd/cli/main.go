package main

import (
	"flag"
	"fmt"
	"gomnirun/core/config"
	"gomnirun/core/executor"
	"log"
	"strings"
)

func main() {
	// Accept a dynamic config file
	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	fmt.Printf("Starting GomniRun in CLI Mode with config: %s...\n", *configFile)

	// Load the specified config file
	conf := loadConfig(*configFile)

	// Handle user variable updates via command-line arguments
	updateConfigVariables(&conf, flag.Args())

	// Save configuration if needed
	saveConfig(*configFile, &conf)

	// Execute the command
	executeCommand(conf)
}

func loadConfig(configFile string) config.Config {
	conf, err := config.Load(configFile)
	if err != nil {
		log.Println("No existing config found. Using defaults.")
		conf = config.Config{}
	}
	return conf
}

func updateConfigVariables(conf *config.Config, args []string) {
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			key, value := parts[0], parts[1]

			if _, exists := conf.Variables[key]; exists {
				conf.Variables[key] = config.Variable{Type: conf.Variables[key].Type, Value: value}
				fmt.Printf("Updated %s -> %s\n", key, value)
			} else {
				fmt.Printf("Warning: Variable %s not found in config.\n", key)
			}
		}
	}
}

func saveConfig(configFile string, conf *config.Config) {
	if conf.Overwrite {
		config.Save(configFile, *conf)
		fmt.Println("Configuration updated and saved.")
	} else {
		fmt.Println("Overwrite is disabled. Changes will not be saved.")
	}
}

func executeCommand(conf config.Config) {
	finalCommand := executor.ReplacePlaceholders(conf.CommandTemplate, conf.Variables)
	fmt.Println("Executing command:", finalCommand)

	output, err := executor.RunScript(conf.CommandTemplate, conf.Variables)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Println("Script Output:", output)
}
