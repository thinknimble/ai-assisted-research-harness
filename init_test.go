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

func fakeLineReader(line string) func() (string, error) {
	return func() (string, error) { return line, nil }
}

func TestInitCreatesDirectoryStructure(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	target := filepath.Join(dir, "my-research")

	var buf bytes.Buffer
	if err := runInit([]string{target}, &buf, fakeKeyReader("sk-test"), nil, false); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Verify directories exist
	for _, sub := range []string{"raw", "formatted", "output"} {
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

	// Verify success message includes output/
	output := buf.String()
	if !strings.Contains(output, "Created research directory") {
		t.Errorf("expected success message, got: %s", output)
	}
	if !strings.Contains(output, "output/") {
		t.Errorf("expected output/ in directory listing, got: %s", output)
	}
}

func TestInitDefaultsToCurrentDir(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	var buf bytes.Buffer
	if err := runInit(nil, &buf, fakeKeyReader(""), nil, false); err != nil {
		t.Fatalf("runInit with no args failed: %v", err)
	}

	for _, name := range []string{"raw", "formatted", "output", "doc-template.yaml", ".env"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("expected %s in current directory", name)
		}
	}
}

func TestInitAdoptsExistingFilesIntoRaw(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "existing.txt"), []byte("data"), 0644)
	os.WriteFile(filepath.Join(dir, "notes.pdf"), []byte("pdf"), 0644)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), nil, true)
	if err != nil {
		t.Fatalf("expected init to adopt files, got error: %v", err)
	}

	// Files should be moved into raw/
	if _, err := os.Stat(filepath.Join(dir, "raw", "existing.txt")); err != nil {
		t.Error("expected existing.txt in raw/")
	}
	if _, err := os.Stat(filepath.Join(dir, "raw", "notes.pdf")); err != nil {
		t.Error("expected notes.pdf in raw/")
	}
	// Original locations should be gone
	if _, err := os.Stat(filepath.Join(dir, "existing.txt")); err == nil {
		t.Error("existing.txt should have been moved from root")
	}
	// formatted/ and output/ should exist
	if _, err := os.Stat(filepath.Join(dir, "formatted")); err != nil {
		t.Error("expected formatted/ to be created")
	}
	if _, err := os.Stat(filepath.Join(dir, "output")); err != nil {
		t.Error("expected output/ to be created")
	}

	output := buf.String()
	if !strings.Contains(output, "Moved 2 existing files into raw/") {
		t.Errorf("expected adoption message, got: %s", output)
	}
}

func TestInitAllowsExistingResearchDirectory(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "raw"), 0755)
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), nil, false)
	if err != nil {
		t.Fatalf("expected init to succeed for existing research directory, got: %v", err)
	}
}

func TestInitExistingResearchDirGetsOutputDir(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "raw"), 0755)
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)
	// No output/ exists yet

	var buf bytes.Buffer
	if err := runInit([]string{dir}, &buf, fakeKeyReader(""), nil, false); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	info, err := os.Stat(filepath.Join(dir, "output"))
	if err != nil {
		t.Fatalf("expected output/ to be created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected output to be a directory")
	}
}

func TestInitExistingResearchDirRegistersAndPrintsMessage(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	dir := filepath.Join(tmp, "my-papers")
	os.MkdirAll(filepath.Join(dir, "raw"), 0755)
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader(""), nil, false)
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
	if err := runInit([]string{dir}, &buf, fakeKeyReader(""), nil, false); err != nil {
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
	if err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), nil, false); err != nil {
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

func TestInitAdoptsFilesAndSubdirs(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("stuff"), 0644)
	os.MkdirAll(filepath.Join(dir, "other-dir"), 0755)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), nil, true)
	if err != nil {
		t.Fatalf("expected init to adopt content, got error: %v", err)
	}

	// Both files and subdirs should move into raw/
	if _, err := os.Stat(filepath.Join(dir, "raw", "notes.txt")); err != nil {
		t.Error("expected notes.txt in raw/")
	}
	if _, err := os.Stat(filepath.Join(dir, "raw", "other-dir")); err != nil {
		t.Error("expected other-dir in raw/")
	}
}

func TestInitPromptsForAPIKey(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	target := filepath.Join(dir, "project")

	var buf bytes.Buffer
	if err := runInit([]string{target}, &buf, fakeKeyReader("sk-ant-my-secret-key"), nil, false); err != nil {
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
	if err := runInit([]string{target}, &buf, fakeKeyReader("sk-test"), nil, false); err != nil {
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

func TestInitRespectsResearchIgnore(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "keep.txt"), []byte("keep"), 0644)
	os.WriteFile(filepath.Join(dir, "skip.log"), []byte("skip"), 0644)
	os.WriteFile(filepath.Join(dir, "also-skip.tmp"), []byte("skip"), 0644)
	os.WriteFile(filepath.Join(dir, ".researchignore"), []byte("*.log\n*.tmp\n"), 0644)

	var buf bytes.Buffer
	if err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), nil, true); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// keep.txt should be moved into raw/
	if _, err := os.Stat(filepath.Join(dir, "raw", "keep.txt")); err != nil {
		t.Error("expected keep.txt in raw/")
	}

	// Ignored files should NOT be moved
	if _, err := os.Stat(filepath.Join(dir, "skip.log")); err != nil {
		t.Error("skip.log should remain in root (ignored)")
	}
	if _, err := os.Stat(filepath.Join(dir, "also-skip.tmp")); err != nil {
		t.Error("also-skip.tmp should remain in root (ignored)")
	}

	// .researchignore itself should stay in root
	if _, err := os.Stat(filepath.Join(dir, ".researchignore")); err != nil {
		t.Error(".researchignore should remain in root")
	}

	// Only 1 file should have been moved
	if !strings.Contains(buf.String(), "Moved 1 existing files into raw/") {
		t.Errorf("expected 1 file moved, got: %s", buf.String())
	}
}

func TestInitSkippedKeyWritesPlaceholder(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	target := filepath.Join(dir, "project")

	var buf bytes.Buffer
	if err := runInit([]string{target}, &buf, fakeKeyReader(""), nil, false); err != nil {
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

// --- .researchignore creation tests ---

func TestInitCreatesDefaultResearchIgnore(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	target := filepath.Join(dir, "my-research")

	var buf bytes.Buffer
	if err := runInit([]string{target}, &buf, fakeKeyReader("sk-test"), nil, false); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(target, ".researchignore"))
	if err != nil {
		t.Fatalf("expected .researchignore to exist: %v", err)
	}

	content := string(data)
	for _, pattern := range []string{".DS_Store", "Thumbs.db", "*.swp", "*~", ".env"} {
		if !strings.Contains(content, pattern) {
			t.Errorf("expected .researchignore to contain %q, got: %s", pattern, content)
		}
	}
}

func TestInitDoesNotOverwriteExistingResearchIgnore(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	custom := "my-custom-pattern\n"
	os.WriteFile(filepath.Join(dir, ".researchignore"), []byte(custom), 0644)
	// Need a file to adopt so we go through the adoption path
	os.WriteFile(filepath.Join(dir, "file.txt"), []byte("data"), 0644)

	var buf bytes.Buffer
	if err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), nil, true); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, ".researchignore"))
	if string(data) != custom {
		t.Errorf(".researchignore was overwritten, got: %s", string(data))
	}
}

func TestInitShowsEditTipWhenAdoptingFiles(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "paper.pdf"), []byte("pdf"), 0644)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), nil, true)
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "Tip: add a .researchignore") {
		t.Error("should not show old 'add' tip")
	}
	if !strings.Contains(output, "Tip: edit .researchignore to customize which files are excluded from adoption.") {
		t.Errorf("expected updated tip, got: %s", output)
	}
}

// --- Confirmation prompt tests ---

func TestInitConfirmationListsFilesAndShowsTip(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "paper.pdf"), []byte("pdf"), 0644)
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("notes"), 0644)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), fakeLineReader("y"), false)
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Files to move into raw/:") {
		t.Errorf("expected file listing header, got: %s", output)
	}
	if !strings.Contains(output, "paper.pdf") {
		t.Errorf("expected paper.pdf in listing, got: %s", output)
	}
	if !strings.Contains(output, "notes.txt") {
		t.Errorf("expected notes.txt in listing, got: %s", output)
	}
	if !strings.Contains(output, "Tip: edit .researchignore to customize which files are excluded from adoption.") {
		t.Errorf("expected updated .researchignore tip, got: %s", output)
	}
	if !strings.Contains(output, "Proceed? [y/N]") {
		t.Errorf("expected confirmation prompt, got: %s", output)
	}
}

func TestInitConfirmationAbortsOnNo(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "paper.pdf"), []byte("pdf"), 0644)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), fakeLineReader("n"), false)
	if err != nil {
		t.Fatalf("expected no error on abort, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Aborted. No files were moved.") {
		t.Errorf("expected abort message, got: %s", output)
	}

	// File should NOT have been moved
	if _, err := os.Stat(filepath.Join(dir, "paper.pdf")); err != nil {
		t.Error("paper.pdf should still be in root after abort")
	}
	// raw/ should NOT exist
	if _, err := os.Stat(filepath.Join(dir, "raw")); err == nil {
		t.Error("raw/ should not be created after abort")
	}
}

func TestInitConfirmationAbortsOnEmptyInput(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "paper.pdf"), []byte("pdf"), 0644)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), fakeLineReader(""), false)
	if err != nil {
		t.Fatalf("expected no error on empty input, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Aborted. No files were moved.") {
		t.Errorf("expected abort message on empty input, got: %s", output)
	}

	// File should still be in root
	if _, err := os.Stat(filepath.Join(dir, "paper.pdf")); err != nil {
		t.Error("paper.pdf should still be in root after empty-input abort")
	}
}

func TestInitConfirmationProceedsOnY(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "paper.pdf"), []byte("pdf"), 0644)

	var buf bytes.Buffer
	err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), fakeLineReader("Y"), false)
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// File should be moved
	if _, err := os.Stat(filepath.Join(dir, "raw", "paper.pdf")); err != nil {
		t.Error("expected paper.pdf in raw/ after confirming with Y")
	}
	if _, err := os.Stat(filepath.Join(dir, "paper.pdf")); err == nil {
		t.Error("paper.pdf should not be in root after confirming")
	}
}

func TestInitNoInputSkipsConfirmation(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "paper.pdf"), []byte("pdf"), 0644)

	var buf bytes.Buffer
	// readLine is nil — should never be called with noInput=true
	err := runInit([]string{dir}, &buf, fakeKeyReader("sk-test"), nil, true)
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	output := buf.String()
	// Should NOT contain the confirmation prompt
	if strings.Contains(output, "Proceed? [y/N]") {
		t.Errorf("should not prompt when --no-input is set, got: %s", output)
	}
	// File should be moved
	if _, err := os.Stat(filepath.Join(dir, "raw", "paper.pdf")); err != nil {
		t.Error("expected paper.pdf in raw/ with --no-input")
	}
}
