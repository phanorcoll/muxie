// Package config provides configuration structures and utilities for the application.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the root configuration structure for muxie.
type Config struct {
	Sessions []Session `yaml:"sessions"`
}

// Session defines a session with a name, working directory, and associated windows.
type Session struct {
	Name      string   `yaml:"name"`
	Directory string   `yaml:"directory"`
	Windows   []Window `yaml:"windows"`
}

// Window represents a window within a session, containing multiple panes split in a layout.
type Window struct {
	Name      string `yaml:"name"`
	Directory string `yaml:"directory"`
	Panes     []Pane `yaml:"panes"`
	Layout    string `yaml:"layout"`
}

// Pane defines a single pane within a window, with its associated command.
type Pane struct {
	Command   string `yaml:"command"`
	Directory string `yaml:"directory"`
}

// createExampleConfigFile creates an example configuration file if one does not already exist.
func createExampleConfigFile(configDir string) error {
	exampleConfigFile := filepath.Join(configDir, "config_example.yaml")
	_, err := os.Stat(exampleConfigFile)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("could not stat example config file: %w", err)
	}

	exampleContent := `sessions:
  - name: "My Awesome Project"
    directory: "~/projects/my-awesome-project"
    windows:
      - name: "Code"
        layout: "vertical"
        panes:
          - command: "nvim"
          - command: "git status"
      - name: "Server"
        panes:
          - command: "npm run dev"
  - name: "Another Project"
    directory: "~/projects/another-project"
    windows:
      - name: "Editor"
        panes:
          - command: "vim"
`
	if err := os.WriteFile(exampleConfigFile, []byte(exampleContent), 0644); err != nil {
		return fmt.Errorf("could not write example config file: %w", err)
	}

	return nil
}

// Load reads the muxie configuration file from the user's home directory and returns a Config struct.
// If the configuration file does not exist, it returns a default Config.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get user home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "muxie")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.yml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := createExampleConfigFile(configDir); err != nil {
			return nil, err
		}
		return &Config{}, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal config yaml: %w", err)
	}

	return &config, nil
}
