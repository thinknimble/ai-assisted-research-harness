---
id: write-text-file-tool
parent: reception-write-tools
created: 2026-08-13T00:00:00Z
priority: 1
status: not_started
depends-on: output-directory-exists
---

# Reception Has a write_text_file Tool

Reception mode exposes a `write_text_file` tool for plain-text output formats.

## Success Criteria

- `setupReception()` includes a `write_text_file` tool with parameters `filename` (string, required) and `content` (string, required)
- The tool writes to `output/{filename}` using the existing `writeFile()` helper
- Supported extensions: `.md`, `.csv`, `.json`, `.txt` — the tool rejects filenames with other extensions and returns an error message listing the allowed types
- The tool returns a confirmation string: `"Written output/{filename}"`
- `handleReceptionTool` dispatches `write_text_file` to the write handler
- The reception system prompt is updated to describe the tool and when to use it (e.g., "Use write_text_file to save research summaries, data extractions, or structured output")
