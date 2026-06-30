---
id: streaming-responses
created: 2026-06-30T00:00:00Z
priority: 1
---

# Streaming Responses

The CLI currently blocks silently during API calls and tool execution. Users have no
visibility into what the agent is doing — no text output, no tool-call indicators, nothing
until the full response is ready.

After this spec is complete, the CLI feels alive the entire time: text appears as it's
generated, and tool activity is visible between streaming chunks.

## Context

- `agent.go` contains `runToolLoop` and `sendMessage` — the core response handling
- `sendMessage` uses the blocking `client.Messages.New` API
- The Anthropic Go SDK (`v1.53.0`) supports streaming via `client.Messages.Stream`
- Both modes (backoffice and reception) use the same `runToolLoop`, so streaming benefits both
