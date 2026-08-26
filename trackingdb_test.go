package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenTrackingDBCreatesFileAndTable(t *testing.T) {
	dir := t.TempDir()
	db, err := OpenTrackingDB(dir)
	if err != nil {
		t.Fatalf("OpenTrackingDB: %v", err)
	}
	defer db.Close()

	// File should exist on disk
	if _, err := os.Stat(filepath.Join(dir, ".research-assistant.db")); err != nil {
		t.Fatal("expected .research-assistant.db to exist")
	}

	// Table should exist with correct columns
	rows, err := db.Query(`PRAGMA table_info(processed_files)`)
	if err != nil {
		t.Fatalf("PRAGMA table_info: %v", err)
	}
	defer rows.Close()

	cols := map[string]string{} // name -> type
	for rows.Next() {
		var cid int
		var name, typ string
		var notnull int
		var dfltValue *string
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			t.Fatalf("scan: %v", err)
		}
		cols[name] = typ
	}

	expect := map[string]string{
		"raw_path":       "TEXT",
		"formatted_path": "TEXT",
		"processed_at":   "TIMESTAMP",
	}
	for col, wantType := range expect {
		if got, ok := cols[col]; !ok {
			t.Errorf("missing column %q", col)
		} else if got != wantType {
			t.Errorf("column %q: want type %q, got %q", col, wantType, got)
		}
	}
}

func TestOpenTrackingDBIsIdempotent(t *testing.T) {
	dir := t.TempDir()

	// First open — creates DB and inserts a row
	db1, err := OpenTrackingDB(dir)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	_, err = db1.Exec(`INSERT INTO processed_files (raw_path, formatted_path) VALUES ('raw/a.txt', 'formatted/a.md')`)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	db1.Close()

	// Second open — should not error or reset table
	db2, err := OpenTrackingDB(dir)
	if err != nil {
		t.Fatalf("second open: %v", err)
	}
	defer db2.Close()

	var count int
	if err := db2.QueryRow(`SELECT COUNT(*) FROM processed_files`).Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 row after idempotent reopen, got %d", count)
	}
}

func TestInitCreatesTrackingDB(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	target := filepath.Join(dir, "my-research")

	var buf bytes.Buffer
	if err := runInit([]string{target}, &buf, fakeKeyReader("sk-test"), nil, false); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	if _, err := os.Stat(filepath.Join(target, ".research-assistant.db")); err != nil {
		t.Fatal("expected .research-assistant.db to be created by init")
	}
}

func TestInitExistingDirCreatesTrackingDB(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "raw"), 0755)
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	var buf bytes.Buffer
	if err := runInit([]string{dir}, &buf, fakeKeyReader(""), nil, false); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, ".research-assistant.db")); err != nil {
		t.Fatal("expected .research-assistant.db to be created for existing research dir")
	}
}

func TestDefaultResearchIgnoreIncludesTrackingDB(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	target := filepath.Join(dir, "my-research")

	var buf bytes.Buffer
	if err := runInit([]string{target}, &buf, fakeKeyReader("sk-test"), nil, false); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(target, ".researchignore"))
	if err != nil {
		t.Fatalf("read .researchignore: %v", err)
	}
	if !strings.Contains(string(data), ".research-assistant.db") {
		t.Error(".researchignore should contain .research-assistant.db")
	}
}
