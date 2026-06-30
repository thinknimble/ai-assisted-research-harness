package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunRepos_NoConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	var buf bytes.Buffer
	if err := runRepos(&buf); err != nil {
		t.Fatalf("runRepos: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "research-assistant init") {
		t.Errorf("expected suggestion to run init, got: %s", out)
	}
}

func TestRunRepos_ListsReposWithDefault(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg := &GlobalConfig{
		Repos: map[string]string{
			"alpha": "/path/to/alpha",
			"beta":  "/path/to/beta",
		},
		Default: "alpha",
	}
	if err := SaveGlobalConfig(cfg); err != nil {
		t.Fatalf("SaveGlobalConfig: %v", err)
	}

	var buf bytes.Buffer
	if err := runRepos(&buf); err != nil {
		t.Fatalf("runRepos: %v", err)
	}

	out := buf.String()

	// Default repo should be marked with *
	if !strings.Contains(out, "* alpha") {
		t.Errorf("expected default marker for alpha, got:\n%s", out)
	}
	// Non-default should not have *
	if !strings.Contains(out, "  beta") {
		t.Errorf("expected beta without marker, got:\n%s", out)
	}
	// Both paths should appear
	if !strings.Contains(out, "/path/to/alpha") {
		t.Errorf("expected alpha path in output")
	}
	if !strings.Contains(out, "/path/to/beta") {
		t.Errorf("expected beta path in output")
	}
}
