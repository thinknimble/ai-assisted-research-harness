---
id: message-history-self-heals
parent: conversation-recovery
created: 2026-08-26T00:00:00Z
priority: 1
status: done
---

# Message history self-heals orphaned tool_use blocks before each API call

A defensive check runs before every `sendMessageStreaming` call. If the message history contains any assistant message with tool_use blocks that lack a corresponding tool_result in the following user message, synthetic tool_result blocks are injected to repair the history.

## Success Criteria

- Before each API call in `runToolLoop`, the message slice is scanned for orphaned tool_use blocks
- Any orphaned tool_use gets a synthetic `tool_result` (with `is_error: true`, message `"Tool call was not completed"`) inserted into the history
- A test: a pre-corrupted message history with an orphaned tool_use does not cause a 400 error — the self-heal patches it and the API call succeeds
- The healing is idempotent — running it on an already-valid history is a no-op
