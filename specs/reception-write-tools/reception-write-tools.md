---
id: reception-write-tools
created: 2026-08-13T00:00:00Z
priority: 1
---

# Reception Write Tools

Reception mode is currently read-only — it can list stubs and read raw files, but it cannot produce any output files. Users asking the AI to synthesize research, extract tables, or generate reports have to copy-paste from the terminal.

After this spec is complete, reception can write output files in both text-based formats (markdown, CSV, JSON, plain text) and binary formats (xlsx). The LLM provides structured data through tool parameters, and the Go tool handlers assemble the final files — including binary formats the LLM cannot produce directly.

## Context

- `reception.go` defines `setupReception()` with the current tool set (`list_stubs`, `read_raw_file`)
- `docs.go` has `writeFile(dir, filename, content)` which handles path resolution and sandboxing
- `excelize` is already a dependency and supports both reading and writing `.xlsx`
- `init.go` creates the directory scaffold — needs to include `output/` alongside `raw/` and `formatted/`
- All output files go to `output/` within the research repo
