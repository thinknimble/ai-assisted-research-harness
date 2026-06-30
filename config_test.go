package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGlobalConfig_CreatesOnFirstUse(t *testing.T) {
	// Use a temp dir as HOME so we don't touch the real config
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("LoadGlobalConfig: %v", err)
	}

	// Config should have empty repos map and no default
	if len(cfg.Repos) != 0 {
		t.Errorf("expected empty repos, got %v", cfg.Repos)
	}
	if cfg.Default != "" {
		t.Errorf("expected empty default, got %q", cfg.Default)
	}

	// File and directory should exist on disk
	expectedPath := filepath.Join(tmp, ".research-assistant", "config.yaml")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Fatalf("config file not created at %s: %v", expectedPath, err)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	original := &GlobalConfig{
		Repos: map[string]string{
			"project-a": "/home/user/research/project-a",
			"project-b": "/home/user/research/project-b",
		},
		Default: "project-a",
	}

	if err := SaveGlobalConfig(original); err != nil {
		t.Fatalf("SaveGlobalConfig: %v", err)
	}

	loaded, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("LoadGlobalConfig: %v", err)
	}

	if loaded.Default != original.Default {
		t.Errorf("default: got %q, want %q", loaded.Default, original.Default)
	}
	for name, path := range original.Repos {
		if loaded.Repos[name] != path {
			t.Errorf("repos[%q]: got %q, want %q", name, loaded.Repos[name], path)
		}
	}
	if len(loaded.Repos) != len(original.Repos) {
		t.Errorf("repos length: got %d, want %d", len(loaded.Repos), len(original.Repos))
	}
}

func TestConfigDir_IsUnderHome(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := configDir()
	expected := filepath.Join(tmp, ".research-assistant")
	if dir != expected {
		t.Errorf("configDir: got %q, want %q", dir, expected)
	}
}
