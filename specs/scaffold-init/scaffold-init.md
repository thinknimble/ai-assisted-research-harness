---
id: scaffold-init
created: 2026-06-30T00:00:00Z
priority: 1
---

# Scaffold Init Command

Non-technical users need a single command to set up their research directory. Today
they'd have to manually create folders, copy the template, and hand-edit a `.env` file.
That's too much friction for sharing this tool.

After this spec is complete, a user downloads the binary, runs `research-assistant init`,
and is ready to go — no manual file creation, no config editing, one binary with
everything embedded.

## Context

- `doc-template.yaml` defines the frontmatter schema and must be available at runtime
- `.env` holds `ANTHROPIC_API_KEY` — the only required config
- `raw/` is where users place their research documents
- `formatted/` is managed by the backoffice mode
- The tool is a single Go binary — external file dependencies should be eliminated
