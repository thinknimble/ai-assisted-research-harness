package main

import (
	"database/sql"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const trackingDBFile = ".research-assistant.db"

// OpenTrackingDB opens (or creates) the tracking database at the given project
// root and ensures the schema exists. It is safe to call multiple times — the
// CREATE TABLE uses IF NOT EXISTS, so existing data is never reset.
func OpenTrackingDB(root string) (*sql.DB, error) {
	dbPath := filepath.Join(root, trackingDBFile)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS processed_files (
		raw_path      TEXT PRIMARY KEY,
		formatted_path TEXT NOT NULL,
		processed_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
