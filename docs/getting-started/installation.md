# Installation

## What you need

- **A computer** running macOS, Linux, or Windows
- **An Anthropic API key** — this is what lets the assistant use Claude (the AI) to read and answer questions. You can get one at [console.anthropic.com](https://console.anthropic.com/)

## Step 1: Download the tool

Ask your team lead for the latest `research-assistant` file, or build it yourself if you have Go installed:

```bash
go build -o research-assistant .
```

This creates a file called `research-assistant` in your project folder.

## Step 2: Create a research library

Run the `init` command to set up a new research library. This creates the folder structure and prompts you for your API key:

```bash
research-assistant init my-research
```

You'll see:

```
Enter your Anthropic API key: ********
Created research directory at /Users/you/my-research
  raw/
  formatted/
  doc-template.yaml
  .env
  registered as "my-research" in global config

Next: run:
  research-assistant --mode backoffice
```

!!! tip "What `init` creates"
    - `raw/` — where you'll put your full documents
    - `formatted/` — where the assistant stores document summaries
    - `doc-template.yaml` — the template that defines how summaries are structured
    - `.env` — your API key (kept private, never shared)

!!! warning "Keep your key private"
    The `.env` file contains your personal API key. Never share it or commit it to GitHub.

## Step 3: Verify it works

```bash
research-assistant --mode reception
```

You should see:

```
research-assistant [reception mode]
Type your message, or 'quit' to exit.

>
```

Type `quit` to exit. You're all set!

## Troubleshooting

**"ANTHROPIC_API_KEY environment variable is required"**
:   Your `.env` file is missing or the key is not set correctly. Make sure it contains `ANTHROPIC_API_KEY=sk-ant-...`

**"command not found: research-assistant"**
:   You need to use `./research-assistant` (with the `./` prefix), or add the tool to a directory in your PATH.

**"permission denied"**
:   On macOS or Linux, make the file executable:
    ```bash
    chmod +x research-assistant
    ```

**"directory already has content — aborting"**
:   The `init` command won't overwrite an existing folder. Either choose a new name or use an empty directory.
