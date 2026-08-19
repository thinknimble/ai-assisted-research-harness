# Installation

## What you need

- **A computer** running macOS, Linux, or Windows
- **An Anthropic API key** — this is what lets the assistant use Claude (the AI) to read and answer questions. You can get one at [console.anthropic.com](https://console.anthropic.com/)

## Quick install (macOS / Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/thinknimble/ai-assisted-research-harness/main/install.sh | sh
```

The script detects your platform, downloads the latest release, and installs to `~/.local/bin` (override with `RESEARCH_INSTALL_DIR`). If `~/.local/bin` is not on your `PATH`, the script prints the exact line to add to your shell config. Verify:

```bash
research-assistant --help
```

### Windows (PowerShell)

```powershell
$arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "arm64" }
$url = "https://github.com/thinknimble/ai-assisted-research-harness/releases/latest/download/research-assistant-windows-${arch}.exe"
Invoke-WebRequest -Uri $url -OutFile "$env:LOCALAPPDATA\Microsoft\WindowsApps\research-assistant.exe"
```

??? tip "Manual download"
    Download the binary for your platform from the [Releases page](https://github.com/thinknimble/ai-assisted-research-harness/releases/latest):

    | Platform              | Binary name                              |
    |-----------------------|------------------------------------------|
    | macOS (Apple Silicon) | `research-assistant-darwin-arm64`         |
    | macOS (Intel)         | `research-assistant-darwin-amd64`         |
    | Linux (x86_64)        | `research-assistant-linux-amd64`          |
    | Linux (ARM)           | `research-assistant-linux-arm64`          |
    | Windows (x86_64)      | `research-assistant-windows-amd64.exe`    |

??? tip "Building from source"
    If you have Go installed, you can build it yourself:
    ```bash
    go build -o research-assistant .
    ```

## Create a research library

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

The tool lists the files it will move and asks for confirmation:

```
Files to move into raw/:
  report.pdf
  meeting-notes.txt
  data.xlsx
Tip: add a .researchignore file to exclude files from adoption.
Proceed? [y/N] y
Moved 3 existing files into raw/
Created research directory at /Users/you/Downloads/my-documents
  raw/
  formatted/
  doc-template.yaml
  .env
```

Answering anything other than `y` or `Y` aborts without moving files.

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

## Verify it works

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

## Updating

```bash
research-assistant update          # install latest
research-assistant update --check  # preview without installing
```

`update` replaces the binary in place, so it needs write access to the install directory. The default install location (`~/.local/bin`) is user-owned, so updates just work. If the binary lives in a root-owned directory like `/usr/local/bin`, run `sudo research-assistant update`. To switch to sudo-free updates, reinstall to the default:

```bash
curl -fsSL https://raw.githubusercontent.com/thinknimble/ai-assisted-research-harness/main/install.sh | sh
```

## Troubleshooting

**"ANTHROPIC_API_KEY environment variable is required"**
:   Your `.env` file is missing or the key is not set correctly. Make sure it contains `ANTHROPIC_API_KEY=sk-ant-...`

**"command not found: research-assistant"**
:   `~/.local/bin` is not on your PATH. The install script prints the exact line to add — check its output above.

**"permission denied"**
:   On macOS or Linux, make the file executable:
    ```bash
    chmod +x research-assistant
    ```

**"failed to move ... to raw/"**
:   The tool tried to adopt your existing files into `raw/` but couldn't move one of them. Check file permissions and try again.
