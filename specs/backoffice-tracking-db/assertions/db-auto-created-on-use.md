---
id: db-auto-created-on-use
parent: backoffice-tracking-db
created: 2026-08-26T00:00:00Z
priority: 1
status: done
depends-on: tracking-db-exists
---

# Tracking DB is auto-created when backoffice runs, not only via init

The backoffice command creates `.research-assistant.db` if it does not exist, so users upgrading from a pre-DB version never hit an error or need to re-run `init`.

## Success Criteria

- Running `research-assistant backoffice` on a repo with no `.research-assistant.db` creates the DB with the `processed_files` schema before processing begins
- The auto-created DB is identical in schema to one created by `init`
- No user action or flag is required — the upgrade path is invisible
- If the DB already exists, it is opened as-is (no reset, no error)
