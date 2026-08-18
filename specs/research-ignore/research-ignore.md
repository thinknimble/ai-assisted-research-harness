---
id: research-ignore
created: 2026-08-18T17:00:00Z
priority: 1
---

# Research Ignore

A `.researchignore` file in each research project root controls which files are excluded from directory listings and from init adoption. Uses gitignore-compatible pattern syntax. Applied at the `listDir()` level in `docs.go` (so both backoffice and reception modes inherit filtering) and during `runInit()` file adoption in `init.go`.
