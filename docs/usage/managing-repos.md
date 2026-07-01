# Managing Repos

Research Assistant can manage **multiple document libraries** (repos). Each one is a separate folder with its own `raw/` and `formatted/` directories.

This is useful if you have different projects or research areas that you want to keep separate.

## Creating a new repo

Use the `init` command:

```bash
research-assistant init ~/research/ai-papers
research-assistant init ~/research/api-docs
```

Each one creates a self-contained library with its own folders and API key.

## Registering an existing repo

If you already have a folder with `raw/` and `formatted/` directories — for example, a repo you cloned from a teammate — run `init` on it:

```bash
research-assistant init ~/projects/shared-research
```

The tool detects the existing structure and registers it without overwriting anything:

```
Found existing research directory, registered as "shared-research"
```

Your existing `doc-template.yaml`, `.env`, and all documents are left untouched.

## Listing your repos

See all registered repos:

```bash
research-assistant repos
```

Output:

```
* ai-papers      /Users/you/research/ai-papers
  api-docs        /Users/you/research/api-docs
```

The `*` marks the default repo — the one used when you don't specify `--repo`.

## Using a specific repo

Add `--repo` to target a specific library:

```bash
research-assistant --mode reception --repo api-docs
```

Without `--repo`, the tool uses your default repo.

## How the default repo works

- The **first repo** you create automatically becomes the default
- When you run `research-assistant --mode reception` without `--repo`, it uses the default
- If no repos are registered and no `--repo` flag is given, the tool uses your current directory

## Where is the config stored?

The global config lives at:

```
~/.research-assistant/config.yaml
```

It's a simple file that maps repo names to folder paths:

```yaml
repos:
  ai-papers: /Users/you/research/ai-papers
  api-docs: /Users/you/research/api-docs
default: ai-papers
```

You can edit this file directly if you need to rename a repo, change the default, or remove an entry.
