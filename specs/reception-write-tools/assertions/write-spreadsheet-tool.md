---
id: write-spreadsheet-tool
parent: reception-write-tools
created: 2026-08-13T00:00:00Z
priority: 1
status: not_started
depends-on: output-directory-exists
---

# Reception Has a write_spreadsheet Tool

Reception mode exposes a `write_spreadsheet` tool that accepts structured row data and produces an `.xlsx` file. The LLM provides data as JSON; the Go handler uses `excelize` to assemble the binary file.

## Success Criteria

- `setupReception()` includes a `write_spreadsheet` tool with parameters:
  - `filename` (string, required) — must end in `.xlsx`; the tool rejects other extensions
  - `headers` (array of strings, required) — column header labels for the first row
  - `rows` (array of arrays of strings, required) — each inner array is one row of cell values
- The tool creates an `.xlsx` file at `output/{filename}` using `excelize.NewFile()`
- Headers are written to row 1 starting at cell A1, data rows follow starting at row 2
- The tool returns a confirmation string: `"Written output/{filename} (N rows)"`
- `handleReceptionTool` dispatches `write_spreadsheet` to the handler
- The reception system prompt describes the tool: "Use write_spreadsheet to export tabular data as an Excel file — pass column headers and row data as JSON arrays"

**Note:** Cell values are all strings. The LLM cannot reliably distinguish numeric vs text intent, and `excelize` string cells display correctly in Excel/Sheets. Numeric formatting can be added later if needed.
