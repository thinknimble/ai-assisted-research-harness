package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

// BackfillFromStubs populates the processed_files table from existing formatted
// stubs when the table is empty. Each stub's YAML frontmatter is parsed for a
// "path" field pointing to the raw file. Stubs without a valid path are skipped
// with a warning on stderr.
func BackfillFromStubs(db *sql.DB, root string, stderr io.Writer) (int, error) {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM processed_files").Scan(&count); err != nil {
		return 0, fmt.Errorf("count processed_files: %w", err)
	}
	if count > 0 {
		return 0, nil // table already has data, skip backfill
	}

	formattedDir := filepath.Join(root, "formatted")
	entries, err := os.ReadDir(formattedDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("read formatted/: %w", err)
	}

	var backfilled int
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		stubPath := filepath.Join(formattedDir, entry.Name())
		rawPath, err := extractFrontmatterPath(stubPath)
		if err != nil || rawPath == "" {
			fmt.Fprintf(stderr, "backfill: skipping %s: no path field in frontmatter\n", entry.Name())
			continue
		}

		absRaw := rawPath
		if !filepath.IsAbs(rawPath) {
			absRaw = filepath.Join(root, rawPath)
		}
		if _, err := os.Stat(absRaw); err != nil {
			fmt.Fprintf(stderr, "backfill: skipping %s: raw file %s does not exist\n", entry.Name(), rawPath)
			continue
		}

		formattedRel := filepath.Join("formatted", entry.Name())
		_, err = db.Exec(
			`INSERT OR IGNORE INTO processed_files (raw_path, formatted_path) VALUES (?, ?)`,
			rawPath, formattedRel,
		)
		if err != nil {
			return backfilled, fmt.Errorf("insert %s: %w", rawPath, err)
		}
		backfilled++
	}
	return backfilled, nil
}

// extractFrontmatterPath reads a file's YAML frontmatter (between --- delimiters)
// and returns the value of the "path" field, or empty string if not found.
func extractFrontmatterPath(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return "", nil // no frontmatter
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break // end of frontmatter
		}
		if strings.HasPrefix(line, "path:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "path:"))
			// Strip surrounding quotes if present
			val = strings.Trim(val, `"'`)
			return val, nil
		}
	}
	return "", scanner.Err()
}
