---
id: init-registers-repo
parent: multi-repo-config
created: 2026-06-30T00:00:00Z
priority: 1
status: done
depends-on: global-config-file
---

# Init Auto-Registers New Repos in Config

When `research-assistant init` scaffolds a directory, it automatically registers
that repo in the global config so the user never has to edit YAML manually.

## Success Criteria

- `research-assistant init ~/work/ai-papers` registers the repo in
  `~/.research-assistant/config.yaml` with a name derived from the directory
  (e.g., `ai-papers`)
- If this is the first registered repo, it becomes the default
- If a repo with the same name already exists in config, the command appends a
  suffix or prompts for a different name — no silent overwrite
