package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCreatesDirectoryStructure(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "my-research")

	var buf bytes.Buffer
	if err := runInit([]string{target}, &buf); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Verify directories exist
	for _, sub := range []string{"raw", "formatted"} {
		info, err := os.Stat(filepath.Join(target, sub))
		if err != nil {
			t.Errorf("expected %s/ to exist: %v", sub, err)
		} else if !info.IsDir() {
			t.Errorf("expected %s to be a directory", sub)
		}
	}

	// Verify doc-template.yaml exists and has content
	data, err := os.ReadFile(filepath.Join(target, "doc-template.yaml"))
	if err != nil {
		t.Fatalf("expected doc-template.yaml to exist: %v", err)
	}
	if len(data) == 0 {
		t.Error("doc-template.yaml is empty")
	}

	// Verify success message
	output := buf.String()
	if !strings.Contains(output, "Created research directory") {
		t.Errorf("expected success message, got: %s", output)
	}
	if !strings.Contains(output, "--mode backoffice") {
		t.Errorf("expected next steps in output, got: %s", output)
	}
}

func TestInitDefaultsToCurrentDir(t *testing.T) {
	dir := t.TempDir()
	// Change to empty temp dir
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	var buf bytes.Buffer
	if err := runInit(nil, &buf); err != nil {
		t.Fatalf("runInit with no args failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "raw")); err != nil {
		t.Error("expected raw/ in current directory")
	}
	if _, err := os.Stat(filepath.Join(dir, "formatted")); err != nil {
		t.Error("expected formatted/ in current directory")
	}
	if _, err := os.Stat(filepath.Join(dir, "doc-template.yaml")); err != nil {
		t.Error("expected doc-template.yaml in current directory")
	}
}

func TestInitAbortsIfDirectoryHasContent(t *testing.T) {
	dir := t.TempDir()
	// Create a file in the target so it's non-empty
	os.WriteFile(filepath.Join(dir, "existing.txt"), []byte("data"), 0644)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf)
	if err == nil {
		t.Fatal("expected error for non-empty directory")
	}
	if !strings.Contains(err.Error(), "already has content") {
		t.Errorf("expected 'already has content' error, got: %v", err)
	}
}
