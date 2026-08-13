# Research Assistant

A simple chat tool that helps you **search your document library** and **organize new research** — all from the command line.

---

## What does it do?

Research Assistant has two modes:

**Reception** — Ask questions and get answers drawn from your document library. The assistant knows what documents you have, reads the ones that matter, and puts together an answer for you.

**Backoffice** — Add new documents to your library. Drop a raw file into the `raw/` folder, then let the assistant read it, summarize it, and file it properly so it's searchable later.

---

## Who is this for?

Anyone who collects research documents (articles, API docs, notes, papers) and wants a fast way to:

- Find information across many documents at once
- Get answers without opening and reading each file yourself
- Keep documents organized with consistent summaries and metadata

---

## Quick start

Download the latest binary from the [Releases page](https://github.com/thinknimble/ai-assisted-research-harness/releases), then:

```bash
# Set up a new research library
research-assistant init my-research

# Start asking questions
research-assistant --mode reception
```

Already have a folder with research documents? Point `init` at it and your files will be moved into `raw/` automatically:

```bash
research-assistant init /path/to/existing-documents
```

See [Getting Started](getting-started/installation.md) for the full walkthrough.

---

## How it works (the short version)

Your documents live in two folders:

| Folder | What's in it |
|---|---|
| `raw/` | The full documents — articles, notes, API references, anything |
| `formatted/` | Short summaries of each document (just a few lines each) |

When you ask a question, the assistant first scans the short summaries to figure out which documents are relevant. Then it reads the full versions of only those documents to build its answer.

When you add a new document, the assistant reads the full file and creates a short summary automatically.

---

## Managing multiple libraries

You can create separate libraries for different projects. The tool remembers all of them and lets you switch with a simple flag:

```bash
research-assistant --repo my-other-project --mode reception
```

See [Managing Repos](usage/managing-repos.md) for details.
