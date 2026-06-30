---
id: tool-calls-print-status
parent: streaming-responses
created: 2026-06-30T00:00:00Z
priority: 1
status: not_started
depends-on: text-streams-to-stdout
---

# Tool Invocations Print a Status Line

When the agent invokes a tool, a visible status line prints to stderr so the user
knows what's happening. This covers the gap between "agent decided to use a tool"
and "tool result came back."

## Success Criteria

- When a tool_use block is received, a line prints to stderr indicating the tool name
  (e.g., `[using list_stubs...]` or `[reading raw/survey-webhooks/typeform.md...]`)
- Status lines go to stderr so they don't pollute stdout if output is piped
- Tool name is always shown; tool input is shown when it adds useful context
  (e.g., file paths, search queries) but large payloads are not dumped
