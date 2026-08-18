---
id: init-respects-ignore
parent: research-ignore
created: 2026-08-18T18:00:00Z
priority: 2
status: done
---

# Init Respects .researchignore

When `runInit()` adopts existing files into `raw/`, it loads `.researchignore` from the target directory and skips matching files.

## Success Criteria

- If a `.researchignore` file exists in the target directory, it is loaded before the adoption loop
- Files matching any pattern are **not** moved into `raw/` — they stay in place
- The `.researchignore` file itself is never moved into `raw/`
- The moved-file count reflects only files that were actually moved
- If `.researchignore` does not exist, all files are adopted as before (backwards compatible)
