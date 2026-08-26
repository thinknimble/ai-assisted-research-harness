---
id: orphaned-tool-use-prevented
parent: conversation-recovery
created: 2026-08-26T00:00:00Z
priority: 1
status: done
---

# Responses with tool_use blocks but non-tool_use stop_reason never corrupt the conversation

When `stop_reason` is `max_tokens` (or any value other than `tool_use`) and the response contains tool_use blocks, those blocks get synthetic tool_result responses before the next API call — preventing orphaned tool_use in the message history.

## Success Criteria

- If a response has `stop_reason != "tool_use"` but contains one or more tool_use content blocks, `runToolLoop` appends synthetic `tool_result` blocks (with `is_error: true` and a message like `"Tool call interrupted — response was truncated"`) immediately after appending the response message
- The conversation continues normally after the synthetic results — the LLM can retry or respond with text
- A test: a simulated `max_tokens` response containing a tool_use block does not cause a 400 error on the next API call

**Tests:** agent_test.go (TestOrphanedToolUsePreventedOnTruncation)
