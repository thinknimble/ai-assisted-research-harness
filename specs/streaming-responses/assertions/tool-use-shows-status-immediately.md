---
id: tool-use-shows-status-immediately
parent: streaming-responses
created: 2026-08-26T00:00:00Z
priority: 1
status: done
---

# Tool-use status line prints as soon as the content block starts streaming

When the LLM composes a large tool call (e.g. writing a multi-page document via `write_text_file`), the `InputJSONDelta` events stream for 30+ seconds before the tool executes. During this time the user sees nothing — dead air between the last text token and the `[using ...]` line.

## Success Criteria

- `sendMessageStreaming` handles `ContentBlockStartEvent` for tool_use blocks and prints a status line to stderr immediately (e.g. `[calling write_text_file...]`)
- The status line appears within ~1 second of the LLM deciding to use a tool — not after the full input JSON is assembled
- The existing `[done]` / `[error]` lines in `runToolLoop` still print after tool execution completes
- Text streaming is unaffected — `TextDelta` events still print to stdout as before
