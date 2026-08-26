package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestBackfillFromStubsPopulatesEmptyDB(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "raw", "subdir"), 0755)
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	// Create 3 raw files
	os.WriteFile(filepath.Join(dir, "raw", "alpha.md"), []byte("raw"), 0644)
	os.WriteFile(filepath.Join(dir, "raw", "beta.md"), []byte("raw"), 0644)
	os.WriteFile(filepath.Join(dir, "raw", "subdir", "gamma.md"), []byte("raw"), 0644)

	// Create 3 formatted stubs with path frontmatter
	stubs := map[string]string{
		"alpha.md": "---\ntitle: Alpha\npath: raw/alpha.md\n---\n",
		"beta.md":  "---\ntitle: Beta\npath: raw/beta.md\n---\n",
		"gamma.md": "---\ntitle: Gamma\npath: raw/subdir/gamma.md\n---\n",
	}
	for name, content := range stubs {
		os.WriteFile(filepath.Join(dir, "formatted", name), []byte(content), 0644)
	}

	db, err := OpenTrackingDB(dir)
	if err != nil {
		t.Fatalf("OpenTrackingDB: %v", err)
	}
	defer db.Close()

	var stderr bytes.Buffer
	n, err := BackfillFromStubs(db, dir, &stderr)
	if err != nil {
		t.Fatalf("BackfillFromStubs: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 backfilled, got %d", n)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM processed_files").Scan(&count)
	if count != 3 {
		t.Errorf("expected 3 rows in processed_files, got %d", count)
	}
}

func TestBackfillSkipsStubsWithoutPath(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "raw"), 0755)
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	os.WriteFile(filepath.Join(dir, "raw", "good.md"), []byte("raw"), 0644)

	// One stub with path, one without
	os.WriteFile(filepath.Join(dir, "formatted", "good.md"), []byte("---\ntitle: Good\npath: raw/good.md\n---\n"), 0644)
	os.WriteFile(filepath.Join(dir, "formatted", "no-path.md"), []byte("---\ntitle: No Path\n---\n"), 0644)

	db, err := OpenTrackingDB(dir)
	if err != nil {
		t.Fatalf("OpenTrackingDB: %v", err)
	}
	defer db.Close()

	var stderr bytes.Buffer
	n, err := BackfillFromStubs(db, dir, &stderr)
	if err != nil {
		t.Fatalf("BackfillFromStubs: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 backfilled, got %d", n)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("no-path.md")) {
		t.Error("expected warning about no-path.md on stderr")
	}
}

func TestBackfillSkipsNonExistentRawFile(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	// Stub points to a raw file that doesn't exist
	os.WriteFile(filepath.Join(dir, "formatted", "ghost.md"), []byte("---\ntitle: Ghost\npath: raw/gone.md\n---\n"), 0644)

	db, err := OpenTrackingDB(dir)
	if err != nil {
		t.Fatalf("OpenTrackingDB: %v", err)
	}
	defer db.Close()

	var stderr bytes.Buffer
	n, err := BackfillFromStubs(db, dir, &stderr)
	if err != nil {
		t.Fatalf("BackfillFromStubs: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 backfilled, got %d", n)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("does not exist")) {
		t.Error("expected warning about non-existent raw file")
	}
}

func TestBackfillSkipsWhenDBNotEmpty(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "raw"), 0755)
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	os.WriteFile(filepath.Join(dir, "raw", "a.md"), []byte("raw"), 0644)
	os.WriteFile(filepath.Join(dir, "formatted", "a.md"), []byte("---\ntitle: A\npath: raw/a.md\n---\n"), 0644)

	db, err := OpenTrackingDB(dir)
	if err != nil {
		t.Fatalf("OpenTrackingDB: %v", err)
	}
	defer db.Close()

	// Pre-populate DB
	db.Exec(`INSERT INTO processed_files (raw_path, formatted_path) VALUES ('raw/existing.md', 'formatted/existing.md')`)

	var stderr bytes.Buffer
	n, err := BackfillFromStubs(db, dir, &stderr)
	if err != nil {
		t.Fatalf("BackfillFromStubs: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 backfilled (table not empty), got %d", n)
	}

	// Should still have only the 1 original row
	var count int
	db.QueryRow("SELECT COUNT(*) FROM processed_files").Scan(&count)
	if count != 1 {
		t.Errorf("expected 1 row (pre-existing), got %d", count)
	}
}
