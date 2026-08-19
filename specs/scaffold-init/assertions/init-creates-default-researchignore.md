---
id: init-creates-default-researchignore
parent: scaffold-init
created: 2026-08-19T00:00:00Z
priority: 1
status: done
---

# Init Creates a Default .researchignore

`runInit()` writes a `.researchignore` file to the project root containing common OS and editor junk patterns.

## Success Criteria

- `runInit()` creates `.researchignore` in the project root with these default patterns:
  ```
  .DS_Store
  Thumbs.db
  *.swp
  *~
  .env
  ```
- If `.researchignore` already exists, it is not overwritten
- The tip line "Tip: add a .researchignore file to exclude files from adoption." is replaced with "Tip: edit .researchignore to customize which files are excluded from adoption."
