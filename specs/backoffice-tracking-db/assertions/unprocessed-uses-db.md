---
id: unprocessed-uses-db
parent: backoffice-tracking-db
created: 2026-08-26T00:00:00Z
priority: 1
status: done
depends-on: write-stub-records-mapping
---

# unprocessedRawFiles queries the tracking DB instead of matching basenames

The filesystem basename comparison is replaced entirely by a DB lookup.

## Success Criteria

- `unprocessedRawFiles()` lists all files in `raw/`, then excludes any whose path appears in the `processed_files` table
- The old basename-matching logic is removed — no more `formattedSet` map
- A raw file with a DB entry is never returned as unprocessed, regardless of what the formatted stub is named
- A raw file without a DB entry is returned as unprocessed, even if a formatted stub with a similar name exists in `formatted/`
- A test: raw file `raw/meetings/foo.md` with DB entry mapping to `formatted/meetings-foo.md` does not appear in the unprocessed list
