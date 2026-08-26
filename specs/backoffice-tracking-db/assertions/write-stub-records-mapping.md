---
id: write-stub-records-mapping
parent: backoffice-tracking-db
created: 2026-08-26T00:00:00Z
priority: 1
status: not_started
depends-on: tracking-db-exists
---

# write_formatted_stub records the raw-to-formatted mapping in the tracking DB

When backoffice mode writes a formatted stub, the mapping from raw file to formatted stub is persisted in the database.

## Success Criteria

- `handleBackofficeTool("write_formatted_stub")` inserts a row into `processed_files` with the raw file path and the formatted stub path after a successful write
- The raw file path comes from the current processing context (the file path passed in the user message to the tool loop)
- If a row for the same `raw_path` already exists, it is replaced (UPSERT) — this handles re-processing gracefully
- A test: after `write_formatted_stub` succeeds, querying `processed_files` for the raw path returns the correct formatted path
