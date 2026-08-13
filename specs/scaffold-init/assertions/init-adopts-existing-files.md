---
id: init-adopts-existing-files
parent: scaffold-init
created: 2026-07-03T00:00:00Z
priority: 1
status: done
depends-on: init-creates-directory-structure
---

# Init Adopts Existing Files Into Raw

When `research-assistant init` targets a directory that has files but no `raw/`
or `formatted/` subdirectories, it creates the research structure and moves all
existing files into `raw/` rather than aborting.

This lets users point `init` at a folder of downloaded documents (emails, PDFs,
etc.) and have them automatically organized into the research directory layout.

## Success Criteria

- If the target directory has files but no `raw/` or `formatted/`, init creates
  both directories and moves all existing entries into `raw/`
- Subdirectories in the target are also moved into `raw/`
- A message reports how many files were moved: "Moved N existing files into raw/"
- After adoption the directory is registered in `~/.research-assistant/config.yaml`
- `.env` and `doc-template.yaml` are created as normal
