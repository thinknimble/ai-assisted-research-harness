# Installation

## What you need

- **A computer** running macOS, Linux, or Windows
- **An Anthropic API key** — this is what lets the assistant use Claude (the AI) to read and answer questions. You can get one at [console.anthropic.com](https://console.anthropic.com/)

## Step 1: Download the tool

Go to the [Releases page](https://github.com/thinknimble/ai-assisted-research-harness/releases) and download the right file for your system:

| System | File |
|---|---|
| Mac (Apple Silicon) | `research-assistant-darwin-arm64` |
| Mac (Intel) | `research-assistant-darwin-amd64` |
| Linux | `research-assistant-linux-amd64` |
| Windows | `research-assistant-windows-amd64.exe` |

Rename the downloaded file to `research-assistant` (or `research-assistant.exe` on Windows) and put it somewhere convenient.

On macOS or Linux, make it executable:

```bash
chmod +x research-assistant
```

??? tip "Building from source"
    If you have Go installed, you can build it yourself:
    ```bash
    go build -o research-assistant .
    ```

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
  output/
  doc-template.yaml
  .env
  registered as "my-research" in global config

Next: run:
  research-assistant --mode backoffice
```

!!! tip "What `init` creates"
    - `raw/` — where you'll put your full documents
    - `formatted/` — where the assistant stores document summaries
    - `output/` — where reception saves generated files
    - `doc-template.yaml` — the template that defines how summaries are structured
    - `.env` — your API key (kept private, never shared)

!!! warning "Keep your key private"
    The `.env` file contains your personal API key. Never share it or commit it to GitHub.

### Using an existing folder of documents

If you already have a folder of documents (emails, PDFs, notes), you can point `init` at it. The tool will move your files into `raw/` and set up the research structure around them.

!!! tip "Exclude files from adoption"
    If you add a [`.researchignore`](../usage/research-ignore.md) file to the folder before running `init`, matching files won't be moved into `raw/`.

```bash
research-assistant init ~/Downloads/my-documents
```

```
Moved 15 existing files into raw/
Created research directory at /Users/you/Downloads/my-documents
  raw/
  formatted/
  doc-template.yaml
  .env
```

### Registering an existing library

If you already have a folder with `raw/` and `formatted/` directories (for example, you cloned a teammate's repo), you can register it without creating new files:

```bash
research-assistant init /path/to/existing-project
```

The tool detects the existing structure and just registers it:

```
Found existing research directory, registered as "existing-project"
```

It won't overwrite your `doc-template.yaml` or `.env` if they already exist.

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
:   You need to use `./research-assistant` (with the `./` prefix), or move the tool to a directory in your PATH.

**"permission denied"**
:   On macOS or Linux, make the file executable:
    ```bash
    chmod +x research-assistant
    ```

**"failed to move ... to raw/"**
:   The tool tried to adopt your existing files into `raw/` but couldn't move one of them. Check file permissions and try again.
