package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fakeKeyReader(key string) func() (string, error) {
	return func() (string, error) { return key, nil }
}

func TestInitCreatesDirectoryStructure(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "my-research")

	var buf bytes.Buffer
	if err := runInit([]string{target}, &buf, fakeKeyReader("sk-test")); err != nil {
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
}

func TestInitDefaultsToCurrentDir(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	var buf bytes.Buffer
	if err := runInit(nil, &buf, fakeKeyReader("")); err != nil {
		t.Fatalf("runInit with no args failed: %v", err)
	}

	for _, name := range []string{"raw", "formatted", "doc-template.yaml", ".env"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("expected %s in current directory", name)
		}
	}
}

func TestInitAbortsIfDirectoryHasNonResearchContent(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "existing.txt"), []byte("data"), 0644)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader(""))
	if err == nil {
		t.Fatal("expected error for non-research directory with content")
	}
	if !strings.Contains(err.Error(), "not a research directory") {
		t.Errorf("expected 'not a research directory' error, got: %v", err)
	}
}

func TestInitAllowsExistingResearchDirectory(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "raw"), 0755)
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"))
	if err != nil {
		t.Fatalf("expected init to succeed for existing research directory, got: %v", err)
	}
}

func TestInitExistingResearchDirRegistersAndPrintsMessage(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	dir := filepath.Join(tmp, "my-papers")
	os.MkdirAll(filepath.Join(dir, "raw"), 0755)
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader(""))
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	output := buf.String()

	// Should print the registration message
	if !strings.Contains(output, "Found existing research directory") {
		t.Errorf("expected 'Found existing research directory' message, got: %s", output)
	}
	if !strings.Contains(output, `registered as "my-papers"`) {
		t.Errorf("expected registration name in message, got: %s", output)
	}

	// Should NOT print "Created research directory"
	if strings.Contains(output, "Created research directory") {
		t.Errorf("should not say 'Created' for existing directory, got: %s", output)
	}

	// Should be registered in global config
	cfg, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("LoadGlobalConfig: %v", err)
	}
	absDir, _ := filepath.Abs(dir)
	if cfg.Repos["my-papers"] != absDir {
		t.Errorf("repo not registered: repos=%v", cfg.Repos)
	}
}

func TestInitExistingDirFillsMissingTemplate(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	os.MkdirAll(filepath.Join(dir, "raw"), 0755)
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)
	// No doc-template.yaml exists

	var buf bytes.Buffer
	if err := runInit([]string{dir}, &buf, fakeKeyReader("")); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// doc-template.yaml should be created
	data, err := os.ReadFile(filepath.Join(dir, "doc-template.yaml"))
	if err != nil {
		t.Fatalf("expected doc-template.yaml to be created: %v", err)
	}
	if len(data) == 0 {
		t.Error("doc-template.yaml should have content")
	}
}

func TestInitExistingDirNeverOverwritesFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	os.MkdirAll(filepath.Join(dir, "raw"), 0755)
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	// Pre-existing doc-template.yaml with custom content
	customTemplate := "my-custom-template-content"
	os.WriteFile(filepath.Join(dir, "doc-template.yaml"), []byte(customTemplate), 0644)
	// Pre-existing .env with custom content
	customEnv := "MY_CUSTOM_VAR=hello"
	os.WriteFile(filepath.Join(dir, ".env"), []byte(customEnv), 0644)

	var buf bytes.Buffer
	if err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test")); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// doc-template.yaml should NOT be overwritten
	data, _ := os.ReadFile(filepath.Join(dir, "doc-template.yaml"))
	if string(data) != customTemplate {
		t.Errorf("doc-template.yaml was overwritten, got: %s", string(data))
	}

	// .env should NOT be overwritten
	data, _ = os.ReadFile(filepath.Join(dir, ".env"))
	if string(data) != customEnv {
		t.Errorf(".env was overwritten, got: %s", string(data))
	}
}

func TestInitNonResearchDirWithContentWarnsAndExits(t *testing.T) {
	dir := t.TempDir()
	// Has content but no raw/ or formatted/
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("stuff"), 0644)
	os.MkdirAll(filepath.Join(dir, "other-dir"), 0755)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader(""))
	if err == nil {
		t.Fatal("expected error for directory with content but no raw/formatted")
	}
	if !strings.Contains(err.Error(), "not a research directory") {
		t.Errorf("expected 'not a research directory' error, got: %v", err)
	}
}

func TestInitPromptsForAPIKey(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "project")

	var buf bytes.Buffer
	if err := runInit([]string{target}, &buf, fakeKeyReader("sk-ant-my-secret-key")); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Verify prompt was printed
	output := buf.String()
	if !strings.Contains(output, "Enter your Anthropic API key:") {
		t.Errorf("expected API key prompt, got: %s", output)
	}

	// Verify .env contains the key
	data, err := os.ReadFile(filepath.Join(target, ".env"))
	if err != nil {
		t.Fatalf("expected .env to exist: %v", err)
	}
	if !strings.Contains(string(data), "ANTHROPIC_API_KEY=sk-ant-my-secret-key") {
		t.Errorf("expected API key in .env, got: %s", string(data))
	}

	// Verify output says "Next: run:" (not asking to set key)
	if !strings.Contains(output, "Next: run:") {
		t.Errorf("expected 'Next: run:' in output, got: %s", output)
	}
}

func TestInitRegistersRepoInGlobalConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	target := filepath.Join(tmp, "ai-papers")

	var buf bytes.Buffer
	if err := runInit([]string{target}, &buf, fakeKeyReader("sk-test")); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	cfg, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("LoadGlobalConfig: %v", err)
	}

	// Should be registered with name derived from directory
	if cfg.Repos["ai-papers"] != target {
		t.Errorf("repo not registered: repos=%v", cfg.Repos)
	}

	// First repo should be the default
	if cfg.Default != "ai-papers" {
		t.Errorf("expected default %q, got %q", "ai-papers", cfg.Default)
	}

	// Output should mention registration
	if !strings.Contains(buf.String(), `registered as "ai-papers"`) {
		t.Errorf("expected registration message in output, got: %s", buf.String())
	}
}

func TestInitSkippedKeyWritesPlaceholder(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "project")

	var buf bytes.Buffer
	if err := runInit([]string{target}, &buf, fakeKeyReader("")); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(target, ".env"))
	if err != nil {
		t.Fatalf("expected .env to exist: %v", err)
	}
	content := string(data)

	// Should have placeholder comment, not a real key
	if !strings.Contains(content, "console.anthropic.com") {
		t.Errorf("expected placeholder with console URL, got: %s", content)
	}
	if strings.Contains(content, "ANTHROPIC_API_KEY=sk-") {
		t.Errorf("should not have a real key, got: %s", content)
	}

	// Output should tell user to add key
	output := buf.String()
	if !strings.Contains(output, "add your ANTHROPIC_API_KEY") {
		t.Errorf("expected instruction to add key, got: %s", output)
	}
}
