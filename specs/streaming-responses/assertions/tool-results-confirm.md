---
id: tool-results-confirm
parent: streaming-responses
created: 2026-06-30T00:00:00Z
priority: 1
status: done
depends-on: tool-calls-print-status
---

# Tool Results Print Confirmation Before Resuming

**Tests:** agent_test.go

After a tool finishes executing, a brief confirmation prints to stderr before the
next API call resumes streaming. This closes the feedback loop — the user sees
the tool was called AND that it completed.

## Success Criteria

- After each tool execution, a confirmation line prints to stderr
  (e.g., `[done]` or `[done — 2 files found]`)
- If the tool errored, the line indicates failure (e.g., `[error: file not found]`)
- The confirmation appears before the next streaming API call begins
