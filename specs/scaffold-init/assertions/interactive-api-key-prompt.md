---
id: interactive-api-key-prompt
parent: scaffold-init
created: 2026-06-30T00:00:00Z
priority: 1
status: not_started
depends-on: init-creates-directory-structure
---

# Init Prompts for API Key Interactively

Non-technical users should not have to manually create or edit a `.env` file. The
init command asks for their API key and writes it for them.

## Success Criteria

- During `init`, the command prompts "Enter your Anthropic API key:"
- The entered key is written to `.env` as `ANTHROPIC_API_KEY=<key>`
- If the user skips (empty input), `.env` is created with a placeholder
  comment explaining where to get a key and how to set it
- The key input is not echoed to the terminal (masked for security)
