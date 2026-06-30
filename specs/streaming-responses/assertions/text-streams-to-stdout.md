---
id: text-streams-to-stdout
parent: streaming-responses
created: 2026-06-30T00:00:00Z
priority: 1
status: not_started
---

# Text Streams Token-by-Token to Stdout

API responses use the streaming endpoint instead of the blocking `Messages.New` call.
Text content prints to stdout incrementally as delta events arrive — not buffered until
the response is complete.

## Success Criteria

- `sendMessage` (or its replacement) uses the SDK's streaming API
- Each text delta event prints its content to stdout immediately
- The full streamed text is still captured and appended to the conversation history
  so multi-turn context is preserved
- Both backoffice and reception modes stream identically (they share the same loop)
