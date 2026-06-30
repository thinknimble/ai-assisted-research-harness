package main

import (
	"os"
	"path/filepath"
	"strings"
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

func TestRegisterRepo_DerivesNameFromPath(t *testing.T) {
	cfg := &GlobalConfig{Repos: map[string]string{}}

	name := RegisterRepo(cfg, "/home/user/work/ai-papers")

	if name != "ai-papers" {
		t.Errorf("got name %q, want %q", name, "ai-papers")
	}
	if cfg.Repos["ai-papers"] != "/home/user/work/ai-papers" {
		t.Errorf("repo path not registered correctly")
	}
}

func TestRegisterRepo_FirstRepoBecomeDefault(t *testing.T) {
	cfg := &GlobalConfig{Repos: map[string]string{}}

	name := RegisterRepo(cfg, "/home/user/work/ai-papers")

	if cfg.Default != name {
		t.Errorf("first repo should be default, got default=%q", cfg.Default)
	}
}

func TestRegisterRepo_SecondRepoDoesNotOverrideDefault(t *testing.T) {
	cfg := &GlobalConfig{
		Repos:   map[string]string{"first": "/some/path"},
		Default: "first",
	}

	RegisterRepo(cfg, "/home/user/work/second-project")

	if cfg.Default != "first" {
		t.Errorf("default should remain %q, got %q", "first", cfg.Default)
	}
}

func TestRegisterRepo_DuplicateNameGetsSuffix(t *testing.T) {
	cfg := &GlobalConfig{
		Repos: map[string]string{
			"ai-papers": "/old/path/ai-papers",
		},
		Default: "ai-papers",
	}

	name := RegisterRepo(cfg, "/new/path/ai-papers")

	if name != "ai-papers-2" {
		t.Errorf("got name %q, want %q", name, "ai-papers-2")
	}
	// Original should be untouched
	if cfg.Repos["ai-papers"] != "/old/path/ai-papers" {
		t.Errorf("original repo was overwritten")
	}
	if cfg.Repos["ai-papers-2"] != "/new/path/ai-papers" {
		t.Errorf("new repo not registered at suffixed name")
	}
}

func TestRegisterRepo_MultipleDuplicatesIncrementSuffix(t *testing.T) {
	cfg := &GlobalConfig{
		Repos: map[string]string{
			"project":   "/path/1/project",
			"project-2": "/path/2/project",
		},
		Default: "project",
	}

	name := RegisterRepo(cfg, "/path/3/project")

	if name != "project-3" {
		t.Errorf("got name %q, want %q", name, "project-3")
	}
}

func TestResolveProjectRoot_RepoFlagSelectsRegisteredRepo(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg := &GlobalConfig{
		Repos: map[string]string{
			"ai-papers": "/home/user/research/ai-papers",
			"ml-notes":  "/home/user/research/ml-notes",
		},
		Default: "ml-notes",
	}
	if err := SaveGlobalConfig(cfg); err != nil {
		t.Fatalf("SaveGlobalConfig: %v", err)
	}

	root, err := ResolveProjectRoot("ai-papers")
	if err != nil {
		t.Fatalf("ResolveProjectRoot: %v", err)
	}
	if root != "/home/user/research/ai-papers" {
		t.Errorf("got %q, want %q", root, "/home/user/research/ai-papers")
	}
}

func TestResolveProjectRoot_NoFlagUsesDefault(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg := &GlobalConfig{
		Repos: map[string]string{
			"ai-papers": "/home/user/research/ai-papers",
			"ml-notes":  "/home/user/research/ml-notes",
		},
		Default: "ml-notes",
	}
	if err := SaveGlobalConfig(cfg); err != nil {
		t.Fatalf("SaveGlobalConfig: %v", err)
	}

	root, err := ResolveProjectRoot("")
	if err != nil {
		t.Fatalf("ResolveProjectRoot: %v", err)
	}
	if root != "/home/user/research/ml-notes" {
		t.Errorf("got %q, want %q", root, "/home/user/research/ml-notes")
	}
}

func TestResolveProjectRoot_NoConfigFallsToCwd(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	// No config file exists — should fall back to cwd

	cwd, _ := os.Getwd()
	root, err := ResolveProjectRoot("")
	if err != nil {
		t.Fatalf("ResolveProjectRoot: %v", err)
	}
	if root != cwd {
		t.Errorf("got %q, want cwd %q", root, cwd)
	}
}

func TestResolveProjectRoot_UnknownRepoErrors(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg := &GlobalConfig{
		Repos: map[string]string{
			"ai-papers": "/home/user/research/ai-papers",
		},
		Default: "ai-papers",
	}
	if err := SaveGlobalConfig(cfg); err != nil {
		t.Fatalf("SaveGlobalConfig: %v", err)
	}

	_, err := ResolveProjectRoot("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown repo, got nil")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention the unknown name: %v", err)
	}
	if !strings.Contains(err.Error(), "ai-papers") {
		t.Errorf("error should list available repos: %v", err)
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
