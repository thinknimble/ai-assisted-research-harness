package main

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

// GlobalConfig tracks registered research repos and the default.
type GlobalConfig struct {
	Repos   map[string]string `yaml:"repos"`
	Default string            `yaml:"default,omitempty"`
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".research-assistant")
}

func configPath() string {
	return filepath.Join(configDir(), "config.yaml")
}

// LoadGlobalConfig reads the global config file. If the file or directory
// does not exist, it creates them with an empty config.
func LoadGlobalConfig() (*GlobalConfig, error) {
	path := configPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// File doesn't exist — create it
		cfg := &GlobalConfig{Repos: map[string]string{}}
		if err := SaveGlobalConfig(cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	var cfg GlobalConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Repos == nil {
		cfg.Repos = map[string]string{}
	}
	return &cfg, nil
}

// RegisterRepo adds a repo to the global config with a name derived from the
// directory basename. If it's the first repo, it becomes the default. If a repo
// with the same name already exists, a numeric suffix is appended.
func RegisterRepo(cfg *GlobalConfig, absPath string) string {
	base := filepath.Base(absPath)
	name := base

	// If name already taken, append -2, -3, etc.
	if _, exists := cfg.Repos[name]; exists {
		for i := 2; ; i++ {
			candidate := fmt.Sprintf("%s-%d", base, i)
			if _, exists := cfg.Repos[candidate]; !exists {
				name = candidate
				break
			}
		}
	}

	cfg.Repos[name] = absPath

	// First registered repo becomes the default
	if len(cfg.Repos) == 1 {
		cfg.Default = name
	}

	return name
}

// SaveGlobalConfig writes the config to disk, creating the directory if needed.
func SaveGlobalConfig(cfg *GlobalConfig) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}
