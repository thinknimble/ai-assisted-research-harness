---
id: global-config-file
parent: multi-repo-config
created: 2026-06-30T00:00:00Z
priority: 1
status: done
---

# Global Config File Exists at ~/.research-assistant/config.yaml

The tool reads and writes a YAML config file that tracks all registered research
repos and which one is the default.

## Success Criteria

- Config file location is `~/.research-assistant/config.yaml`
- Config schema contains a `repos` map (name → path) and a `default` field
- If the config file does not exist, the tool creates it on first use
- If the config directory does not exist, the tool creates it
- The config is only read when `--repo` is used or when resolving the default —
  running from a research directory with no flags still works via `os.Getwd()`

**Tests:** config_test.go
