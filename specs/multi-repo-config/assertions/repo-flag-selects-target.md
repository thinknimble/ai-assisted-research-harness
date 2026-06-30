---
id: repo-flag-selects-target
parent: multi-repo-config
created: 2026-06-30T00:00:00Z
priority: 1
status: done
depends-on: global-config-file
---

# --repo Flag Selects the Target Research Directory

Users can target any registered repo by name without `cd`-ing into it. When no
flag is given, the tool uses the default repo from config. When no config exists,
falls back to the current working directory (backwards compatible).

## Success Criteria

- `research-assistant --repo ai-papers --mode reception` sets `projectRoot` to
  the path registered under `ai-papers` in config
- `research-assistant --mode backoffice` (no `--repo` flag) uses the `default`
  repo from config
- If no config file exists and no `--repo` flag is given, the tool uses
  `os.Getwd()` as before (no breaking change)
- If `--repo` references a name not in config, the tool prints an error listing
  available repos
