package main

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
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
