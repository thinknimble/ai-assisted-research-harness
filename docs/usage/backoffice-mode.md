# Backoffice Mode

Backoffice mode helps you **add new documents** to your library. It reads raw files, creates summaries, and can push everything to GitHub.

## How to start

```bash
research-assistant --mode backoffice
```

Or target a specific library:

```bash
research-assistant --mode backoffice --repo my-project
```

## Adding a new document

### Step 1: Drop the file in `raw/`

Save your document (a markdown file, plain text, notes, etc.) somewhere inside the `raw/` folder. You can organize it into subfolders:

```
raw/
  survey-webhooks/
    typeform-webhooks.md
    qualtrics-webhooks.md
  papers/
    your-new-paper.md        <-- put it here
```

### Step 2: Ask the assistant to process it

```
> Process the new paper in raw/papers/your-new-paper.md
```

The assistant will:

1. Read the full document
2. Figure out what type of document it is (article, reference, notes, etc.)
3. Write a summary that captures the key points
4. Create a formatted stub in the `formatted/` folder

You'll see its progress as it works:

```
[using read_file: doc-template.yaml...]
[done]
[using read_file: raw/papers/your-new-paper.md...]
[done]
[using write_formatted_stub: your-new-paper.md...]
[done]

Created formatted/your-new-paper.md with a summary of the paper.
```

### Step 3: Review and push

The assistant will show you what it created. If you're happy with it:

```
> Looks good, push it to GitHub
```

It will commit the new stub and push it to the repository:

```
[using git_add_commit_push...]
[done]

Committed and pushed to GitHub.
```

## Processing multiple documents at once

You can ask the assistant to handle several files in one go:

```
> Process all unformatted documents
```

It will compare what's in `raw/` against what's already in `formatted/`, find the gaps, and process each one.

## What the assistant creates

For each raw document, the assistant creates a small "stub" file in `formatted/`. Here's what one looks like:

```yaml
---
title: "Typeform Webhooks API Reference"
type: reference
status: active
created: 2026-06-30
updated: 2026-06-30
author: pari

source: web
source_url: "https://www.typeform.com/developers/webhooks/"
path: raw/survey-webhooks/typeform-webhooks.md

tags: [typeform, webhooks, api, surveys]
related:
  - formatted/qualtrics-webhooks.md
summary: "Complete Typeform Webhooks API reference — CRUD endpoints,
  payload schema, signature verification, and security."
---
```

The key parts are:

- **`path`** — points back to the full document in `raw/`
- **`summary`** — a detailed description so the assistant can decide if this document is relevant to a question without reading the whole thing
- **`tags`** — keywords for categorization
- **`related`** — links to other documents on similar topics

## Editing a stub

If you want to change a summary or fix tags, you can either:

- Edit the file in `formatted/` directly
- Ask the assistant: `Update the summary for the Typeform webhooks stub`
