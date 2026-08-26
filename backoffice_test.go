package main

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBackofficeAutoCreatesDBMatchingInit(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Create a DB via init
	initDir := filepath.Join(t.TempDir(), "init-project")
	var buf bytes.Buffer
	if err := runInit([]string{initDir}, &buf, fakeKeyReader("sk-test"), nil, false); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	// Simulate backoffice auto-creation on a dir with no DB
	backofficeDir := t.TempDir()
	os.MkdirAll(filepath.Join(backofficeDir, "raw"), 0755)

	dbPath := filepath.Join(backofficeDir, ".research-assistant.db")
	if _, err := os.Stat(dbPath); err == nil {
		t.Fatal("DB should not exist before auto-creation")
	}

	autoDB, err := OpenTrackingDB(backofficeDir)
	if err != nil {
		t.Fatalf("OpenTrackingDB: %v", err)
	}
	defer autoDB.Close()

	if _, err := os.Stat(dbPath); err != nil {
		t.Fatal("expected DB to be auto-created")
	}

	// Verify schema matches what init produces
	initDB, err := OpenTrackingDB(initDir)
	if err != nil {
		t.Fatalf("open init DB: %v", err)
	}
	defer initDB.Close()

	initSchema := tableSchema(t, initDB, "init")
	autoSchema := tableSchema(t, autoDB, "backoffice")

	if len(initSchema) != len(autoSchema) {
		t.Fatalf("column count mismatch: init=%d, backoffice=%d", len(initSchema), len(autoSchema))
	}
	for col, wantType := range initSchema {
		if got, ok := autoSchema[col]; !ok {
			t.Errorf("backoffice DB missing column %q", col)
		} else if got != wantType {
			t.Errorf("column %q type mismatch: init=%q, backoffice=%q", col, wantType, got)
		}
	}
}

func TestWriteFormattedStubRecordsMapping(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	oldRoot := projectRoot
	oldRawFile := currentRawFile
	t.Cleanup(func() {
		projectRoot = oldRoot
		currentRawFile = oldRawFile
	})
	projectRoot = dir
	currentRawFile = "raw/survey-webhooks/typeform-webhooks.md"

	input := `{"filename":"typeform-webhooks.md","content":"---\ntitle: test\n---"}`
	result, err := handleBackofficeTool("write_formatted_stub", []byte(input))
	if err != nil {
		t.Fatalf("handleBackofficeTool: %v", err)
	}
	if !strings.Contains(result, "formatted/typeform-webhooks.md") {
		t.Fatalf("unexpected result: %s", result)
	}

	db, err := OpenTrackingDB(dir)
	if err != nil {
		t.Fatalf("OpenTrackingDB: %v", err)
	}
	defer db.Close()

	var formattedPath string
	err = db.QueryRow("SELECT formatted_path FROM processed_files WHERE raw_path = ?", currentRawFile).Scan(&formattedPath)
	if err != nil {
		t.Fatalf("query processed_files: %v", err)
	}
	if formattedPath != "formatted/typeform-webhooks.md" {
		t.Errorf("got formatted_path=%q, want %q", formattedPath, "formatted/typeform-webhooks.md")
	}
}

func TestWriteFormattedStubUpsertsMapping(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "formatted"), 0755)

	oldRoot := projectRoot
	oldRawFile := currentRawFile
	t.Cleanup(func() {
		projectRoot = oldRoot
		currentRawFile = oldRawFile
	})
	projectRoot = dir
	currentRawFile = "raw/doc.md"

	// First write
	input := `{"filename":"doc.md","content":"v1"}`
	if _, err := handleBackofficeTool("write_formatted_stub", []byte(input)); err != nil {
		t.Fatalf("first write: %v", err)
	}

	// Second write with same raw path — should UPSERT
	input2 := `{"filename":"doc-v2.md","content":"v2"}`
	if _, err := handleBackofficeTool("write_formatted_stub", []byte(input2)); err != nil {
		t.Fatalf("second write: %v", err)
	}

	db, err := OpenTrackingDB(dir)
	if err != nil {
		t.Fatalf("OpenTrackingDB: %v", err)
	}
	defer db.Close()

	var formattedPath string
	err = db.QueryRow("SELECT formatted_path FROM processed_files WHERE raw_path = ?", "raw/doc.md").Scan(&formattedPath)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if formattedPath != "formatted/doc-v2.md" {
		t.Errorf("got %q, want %q — UPSERT should have replaced the old path", formattedPath, "formatted/doc-v2.md")
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM processed_files WHERE raw_path = ?", "raw/doc.md").Scan(&count)
	if count != 1 {
		t.Errorf("expected 1 row for raw_path, got %d", count)
	}
}

func tableSchema(t *testing.T, db *sql.DB, label string) map[string]string {
	t.Helper()
	rows, err := db.Query("PRAGMA table_info(processed_files)")
	if err != nil {
		t.Fatalf("%s PRAGMA: %v", label, err)
	}
	defer rows.Close()
	cols := map[string]string{}
	for rows.Next() {
		var cid, notnull, pk int
		var name, typ string
		var dflt *string
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dflt, &pk); err != nil {
			t.Fatalf("%s scan: %v", label, err)
		}
		cols[name] = typ
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("%s rows.Err: %v", label, err)
	}
	if len(cols) == 0 {
		t.Fatalf("%s: no columns in processed_files", label)
	}
	return cols
}
