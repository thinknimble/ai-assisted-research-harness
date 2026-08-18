package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// setupTestProject creates a temp project tree and sets projectRoot.
// Returns a cleanup function.
func setupTestProject(t *testing.T, files []string) string {
	t.Helper()
	dir := t.TempDir()
	for _, f := range files {
		p := filepath.Join(dir, f)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	projectRoot = dir
	return dir
}

func TestListDir_NoIgnoreFile(t *testing.T) {
	setupTestProject(t, []string{
		"raw/a.md",
		"raw/b.md",
		"raw/sub/c.md",
	})

	got, err := listDir("raw")
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	want := []string{"raw/a.md", "raw/b.md", "raw/sub/c.md"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d]=%q, want %q", i, got[i], want[i])
		}
	}
}

func TestListDir_GlobPattern(t *testing.T) {
	dir := setupTestProject(t, []string{
		"raw/notes.md",
		"raw/data.tmp",
		"raw/other.tmp",
	})
	os.WriteFile(filepath.Join(dir, ".researchignore"), []byte("*.tmp\n"), 0644)

	got, err := listDir("raw")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "raw/notes.md" {
		t.Fatalf("expected only raw/notes.md, got %v", got)
	}
}

func TestListDir_DirectoryPattern(t *testing.T) {
	dir := setupTestProject(t, []string{
		"raw/keep.md",
		"raw/drafts/wip.md",
		"raw/drafts/nested/deep.md",
	})
	os.WriteFile(filepath.Join(dir, ".researchignore"), []byte("raw/drafts/\n"), 0644)

	got, err := listDir("raw")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "raw/keep.md" {
		t.Fatalf("expected only raw/keep.md, got %v", got)
	}
}

func TestListDir_NegationPattern(t *testing.T) {
	dir := setupTestProject(t, []string{
		"raw/a.tmp",
		"raw/b.tmp",
		"raw/keep-this.tmp",
	})
	os.WriteFile(filepath.Join(dir, ".researchignore"), []byte("*.tmp\n!keep-this.tmp\n"), 0644)

	got, err := listDir("raw")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "raw/keep-this.tmp" {
		t.Fatalf("expected only raw/keep-this.tmp, got %v", got)
	}
}

func TestListDir_RecursiveGlob(t *testing.T) {
	dir := setupTestProject(t, []string{
		"raw/a.log",
		"raw/sub/b.log",
		"raw/sub/deep/c.log",
		"raw/keep.md",
	})
	os.WriteFile(filepath.Join(dir, ".researchignore"), []byte("**/*.log\n"), 0644)

	got, err := listDir("raw")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "raw/keep.md" {
		t.Fatalf("expected only raw/keep.md, got %v", got)
	}
}

func TestListDir_EmptyIgnoreFile(t *testing.T) {
	dir := setupTestProject(t, []string{
		"raw/a.md",
		"raw/b.md",
	})
	os.WriteFile(filepath.Join(dir, ".researchignore"), []byte(""), 0644)

	got, err := listDir("raw")
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	if len(got) != 2 {
		t.Fatalf("expected 2 files, got %v", got)
	}
}

func TestListDir_CommentsIgnored(t *testing.T) {
	dir := setupTestProject(t, []string{
		"raw/a.md",
		"raw/b.tmp",
	})
	os.WriteFile(filepath.Join(dir, ".researchignore"), []byte("# this is a comment\n*.tmp\n"), 0644)

	got, err := listDir("raw")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "raw/a.md" {
		t.Fatalf("expected only raw/a.md, got %v", got)
	}
}

func TestListDir_FormattedDirAlsoFiltered(t *testing.T) {
	dir := setupTestProject(t, []string{
		"formatted/doc.md",
		"formatted/skip.tmp",
	})
	os.WriteFile(filepath.Join(dir, ".researchignore"), []byte("*.tmp\n"), 0644)

	got, err := listDir("formatted")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "formatted/doc.md" {
		t.Fatalf("expected only formatted/doc.md, got %v", got)
	}
}
