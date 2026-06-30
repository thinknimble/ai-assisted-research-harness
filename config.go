package main

import (
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
