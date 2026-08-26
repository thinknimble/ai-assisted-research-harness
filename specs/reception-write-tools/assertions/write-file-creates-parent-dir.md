---
id: write-file-creates-parent-dir
parent: reception-write-tools
created: 2026-08-26T00:00:00Z
priority: 1
status: done
---

# writeFile creates the target directory if it does not exist

`writeFile()` is sandboxed to `projectRoot/dir` but assumes the directory already exists. If a user deletes `output/` (or never ran `init`), writes fail with "no such file or directory."

## Success Criteria

- `writeFile()` calls `os.MkdirAll` on the target directory before `os.WriteFile`
- `write_text_file` succeeds when `output/` is absent — the directory is created automatically
- `write_spreadsheet_tool` succeeds when `output/` is absent — same behavior
- Existing path-traversal sandboxing is preserved (the `strings.HasPrefix` check still runs)
- No behavior change when the directory already exists
