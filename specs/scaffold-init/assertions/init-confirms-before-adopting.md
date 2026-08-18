---
id: init-confirms-before-adopting
parent: scaffold-init
created: 2026-08-18T17:00:00Z
priority: 1
status: done
depends-on: init-adopts-existing-files
---

# Init Confirms Before Adopting Files

When `research-assistant init` detects files to move into `raw/`, it lists them and asks for confirmation before moving anything.

## Success Criteria

- Init prints the list of files it will move before moving them
- Prompt asks `Proceed? [y/N]` — defaults to "no" (safe by default)
- Only `y` or `Y` proceeds; any other input (including empty/Enter) aborts without moving files
- A tip line mentions `.researchignore`: `Tip: add a .researchignore file to exclude files from adoption.`
- `--no-input` flag on the `init` subcommand bypasses the confirmation and moves immediately (for scripting/CI)
- Existing research directories (already have `raw/` + `formatted/`) skip the confirmation entirely — nothing is moved
- Fresh empty directories skip the confirmation — nothing to move
- Aborting prints a message like `Aborted. No files were moved.` and exits with code 0 (not an error)

**Tests:** init_test.go (TestInitConfirmation*, TestInitNoInputSkipsConfirmation)
