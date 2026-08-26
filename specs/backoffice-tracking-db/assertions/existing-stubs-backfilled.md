---
id: existing-stubs-backfilled
parent: backoffice-tracking-db
created: 2026-08-26T00:00:00Z
priority: 2
status: not_started
depends-on: tracking-db-exists
---

# Existing formatted stubs are backfilled into the tracking DB on first run

Research repos that already have formatted stubs but no DB entries get migrated automatically.

## Success Criteria

- When backoffice mode starts and the `processed_files` table is empty but `formatted/` contains stubs, a backfill runs before processing
- Backfill reads each formatted stub's YAML frontmatter `path` field to determine the raw file it was created from
- For each stub with a valid `path` field pointing to an existing raw file, a row is inserted into `processed_files`
- Stubs without a `path` field or whose `path` points to a non-existent raw file are skipped (logged to stderr as a warning)
- Backfill runs only once — subsequent runs see a non-empty table and skip it
- A test: a repo with 3 formatted stubs and an empty DB has 3 rows in `processed_files` after backfill
