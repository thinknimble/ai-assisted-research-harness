---
id: tracking-db-exists
parent: backoffice-tracking-db
created: 2026-08-26T00:00:00Z
priority: 1
status: not_started
---

# Each research repo has a .research-assistant.db SQLite file

The tracking database lives at the root of each research repo alongside `raw/`, `formatted/`, and `output/`.

## Success Criteria

- `research-assistant init` creates `.research-assistant.db` at the project root if it does not exist
- The DB contains a `processed_files` table with columns: `raw_path` (TEXT PRIMARY KEY), `formatted_path` (TEXT NOT NULL), `processed_at` (TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)
- Existing research directories get the DB created on next init if it does not already exist
- `.research-assistant.db` is added to the default `.researchignore` so it does not appear in `list_raw_files` or `list_formatted_files`
- Opening the DB is idempotent — running init twice does not error or reset the table
