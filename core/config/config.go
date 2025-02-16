package config

import (
	"encoding/json"
	"os"
)

// Variable defines a value and its type
type Variable struct {
	Type  string `json:"type"`  // "file", "directory", "string", "bool", etc.
	Value string `json:"value"` // Actual value
}

// Config stores the command template, variables, and overwrite flag
type Config struct {
	CommandTemplate string              `json:"command"`
	Variables       map[string]Variable `json:"variables"`
	Overwrite       bool                `json:"overwrite"` // New field to control overwriting
}

// Load configuration from file
func Load(filename string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	return config, err
}

// Save configuration to file
func Save(filename string, config Config) error {
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, file, 0644)
}
