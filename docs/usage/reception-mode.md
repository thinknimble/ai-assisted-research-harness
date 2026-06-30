# Reception Mode

Reception mode is your **research assistant**. Ask it questions and it will search your document library for answers.

## How to start

```bash
research-assistant --mode reception
```

Or target a specific library:

```bash
research-assistant --mode reception --repo my-project
```

## Asking questions

Just type your question in plain English and press ++enter++. The assistant will:

1. Scan all your document summaries to find which ones are relevant
2. Read the full text of the relevant documents
3. Give you an answer based on what it found

### Example conversation

```
> How do Typeform webhooks handle security?
[using list_stubs...]
[done]
[using read_raw_file: raw/survey-webhooks/typeform-webhooks.md...]
[done]

Typeform signs each webhook payload using HMAC SHA-256. When you set a secret
on your webhook, Typeform includes a Typeform-Signature header with each
delivery. To verify authenticity, you compute the HMAC SHA-256 hash of the
raw payload using your secret, base64-encode it, and compare it against the
header value.

Sources: raw/survey-webhooks/typeform-webhooks.md
```

### Understanding the status messages

While the assistant works, you'll see status messages in square brackets:

| Message | Meaning |
|---|---|
| `[using list_stubs...]` | Scanning the document index for relevant documents |
| `[using read_raw_file: path...]` | Reading a specific document for details |
| `[done]` | The step completed successfully |
| `[error: ...]` | Something went wrong — the assistant will try to recover |

The final answer streams in word by word as it's generated.

## Tips for good questions

!!! tip "Be specific"
    Instead of "tell me about webhooks", try "how do I verify a Typeform webhook signature in Python?"

!!! tip "Ask follow-ups"
    The assistant remembers your conversation. You can ask "what about Qualtrics?" and it will understand you're comparing webhook approaches.

!!! tip "Ask for comparisons"
    "How do Typeform and Qualtrics webhook payloads differ?" works well when you have docs on both.

## What it can do

- Answer questions about any document in your library
- Compare information across multiple documents
- Summarize long documents
- Find specific details (API endpoints, field names, code examples)

## What it cannot do

- Answer questions about topics not covered in your documents
- Access the internet or look up new information
- Remember conversations after you close the tool
