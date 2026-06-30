# First Run

Once you've [installed the tool](installation.md), let's make sure everything is working.

## Start the assistant

Open your terminal and run:

```bash
research-assistant --mode reception
```

You'll see a prompt:

```
research-assistant [reception mode]
Type your message, or 'quit' to exit.

>
```

## Ask a question

Type a question and press ++enter++. For example:

```
> What documents do we have?
```

## What you'll see

The assistant works in steps, and shows you what it's doing:

```
> How do Typeform webhooks handle security?
[using list_stubs...]
[done]
[using read_raw_file: raw/survey-webhooks/typeform-webhooks.md...]
[done]

Typeform signs each webhook payload using HMAC SHA-256...
```

The lines in square brackets show the assistant's thinking process:

- `[using list_stubs...]` — scanning the document index
- `[using read_raw_file: ...]` — reading a full document
- `[done]` — finished that step
- `[error: ...]` — something went wrong (the assistant will usually try a different approach)

The answer streams in as it's generated — you'll see text appearing word by word.

## Try a few things

Here are some good first questions:

- `What documents do we have?`
- `Summarize the Typeform webhooks documentation`
- `How do I verify a webhook signature?`
- `How do Typeform and Qualtrics webhooks compare?`

## Exit

Type `quit` or `exit` to close the assistant. You can also press ++ctrl+c++.
