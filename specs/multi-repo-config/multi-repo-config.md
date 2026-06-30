---
id: multi-repo-config
created: 2026-06-30T00:00:00Z
priority: 2
---

# Multi-Repo Config

Users may have multiple research directories — one per project, client, or topic.
Today the tool requires `cd`-ing into the right directory since it uses `os.Getwd()`
as the project root. Non-technical users shouldn't have to think about working
directories.

After this spec is complete, the tool maintains a global config file that tracks
registered research repos. Users reference repos by name and the tool resolves
the path automatically.

## Context

- `docs.go` sets `projectRoot` via `os.Getwd()` — this is the single point that
  needs to respect the config
- `init` (from scaffold-init spec) already creates research directories — it should
  auto-register new repos in the config
- The config lives at `~/.research-assistant/config.yaml` (follows XDG-ish conventions)
- The API key lives in each repo's `.env`, not in the global config (different repos
  could use different keys)
