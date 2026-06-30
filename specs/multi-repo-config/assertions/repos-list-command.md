---
id: repos-list-command
parent: multi-repo-config
created: 2026-06-30T00:00:00Z
priority: 2
status: not_started
depends-on: global-config-file
---

# repos Command Lists Registered Repos

Users can see all their registered research directories and which one is the default.

## Success Criteria

- `research-assistant repos` prints all registered repos with name, path, and
  a marker indicating which is the default
- If no repos are registered, prints a message suggesting `research-assistant init`
