# Ignoring Files

You can exclude files from your research library using a `.researchignore` file. This works like `.gitignore` — any file matching a pattern is hidden from both reception and backoffice modes.

## Setup

Create a `.researchignore` file in your research project root (next to `raw/` and `formatted/`):

```bash
touch .researchignore
```

## Pattern syntax

Uses the same syntax as `.gitignore`:

```gitignore
# Ignore all PDFs in the drafts folder
raw/drafts/*.pdf

# Ignore anything in a temp directory
**/temp/

# Ignore a specific file
raw/notes/scratch.md
```

## Examples

**Hide work-in-progress files from the assistant:**

```gitignore
raw/wip/**
raw/**/*.draft.md
```

**Exclude large binary files that can't be read anyway:**

```gitignore
**/*.zip
**/*.tar.gz
```

## How it works

- The `.researchignore` file is loaded every time the tool lists files in your library
- Matched files are silently skipped — they won't appear in file listings and the assistant won't try to read or process them
- The files remain on disk; they're just hidden from the tool
- If the file is missing or empty, nothing is filtered

## Init respects it too

If you place a `.researchignore` in a folder **before** running `research-assistant init`, matching files won't be moved into `raw/`. They stay where they are.

This is useful when you're adopting an existing folder but want to keep certain files (logs, temp files, binaries) out of the research library from the start.
