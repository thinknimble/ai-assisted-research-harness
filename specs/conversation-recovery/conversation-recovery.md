---
id: conversation-recovery
created: 2026-08-26T00:00:00Z
priority: 1
---

# Conversation Recovery

The tool loop in `agent.go` can leave orphaned `tool_use` blocks in the message history — a `tool_use` without a matching `tool_result` in the next message. Once this happens, every subsequent API call fails with a 400 "tool_use ids without tool_result" error and the session is permanently bricked.

The most common trigger: the LLM hits `max_tokens` (4096) while composing a large tool input (e.g. a multi-page document for `write_text_file`). The response contains a partial tool_use block with `stop_reason: "max_tokens"`. `runToolLoop` appends the message but only processes tools when `stop_reason == "tool_use"`, leaving the orphan.

After this spec is complete, orphaned tool_use blocks never corrupt the conversation — they're either prevented at the source or healed before the next API call.

## Context

- `agent.go:runToolLoop` — the core loop that sends messages and processes tool results
- `agent.go:sendMessageStreaming` — streams the response, accumulates the full `Message`
- Both backoffice and reception modes share `runToolLoop`
- The Anthropic API requires every `tool_use` block to have a corresponding `tool_result` in the immediately following message
