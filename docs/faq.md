# Frequently Asked Questions

## General

### What is Research Assistant?

A command-line chat tool that helps you search across a library of research documents and organize new ones. It uses Claude (an AI by Anthropic) to read, summarize, and answer questions about your documents.

### Do I need to know how to code?

No. You just need to be comfortable typing in a terminal. The tool works like a chat — you type questions in plain English and get answers back.

### Does it cost money to use?

The tool itself is free. However, it uses the Anthropic API to power the AI, which has usage-based pricing. Each question costs a small amount (typically fractions of a cent). You'll need an API key from [console.anthropic.com](https://console.anthropic.com/).

---

## Setup

### How do I set up my first library?

Run the `init` command with a folder name:

```bash
research-assistant init my-research
```

It will create the folder structure, ask for your API key, and register the library so you can use it from anywhere.

### Where do I get an API key?

Go to [console.anthropic.com](https://console.anthropic.com/), create an account, and generate an API key. The `init` command will prompt you to enter it.

### I'm getting "ANTHROPIC_API_KEY environment variable is required"

Your `.env` file is missing or doesn't contain the key. Make sure:

1. There's a file called `.env` in your library folder (not `.env.txt` or `env`)
2. It contains a line like: `ANTHROPIC_API_KEY=sk-ant-your-key-here`
3. There are no extra spaces around the `=` sign

If you used `research-assistant init`, this file was created for you.

### Can I use this on Windows?

Yes. Download the Windows binary from the [Releases page](https://github.com/thinknimble/ai-assisted-research-harness/releases), or build from source with `go build -o research-assistant.exe .`. The `.env` file works the same on all platforms.

### I already have a folder with raw/ and formatted/ — do I need to start over?

No. Run `init` on the existing folder:

```bash
research-assistant init /path/to/existing-folder
```

It detects the existing structure and registers it without overwriting your files.

---

## Multiple libraries

### Can I have separate libraries for different projects?

Yes. Create as many as you need:

```bash
research-assistant init ~/research/project-a
research-assistant init ~/research/project-b
```

Use `--repo` to target a specific one:

```bash
research-assistant --mode reception --repo project-b
```

### How do I see all my libraries?

```bash
research-assistant repos
```

This shows all registered libraries and marks the default with `*`.

### How do I change the default library?

Edit `~/.research-assistant/config.yaml` and change the `default` field to the name of the library you want.

---

## Reception mode

### How does the assistant find relevant documents?

It works in two steps:

1. **Scan summaries** — It reads the short summary files in `formatted/`. These are just a few lines each, so this is fast.
2. **Read full documents** — For the documents whose summaries match your question, it reads the complete text from `raw/`.

### What do the messages in square brackets mean?

These are status messages that show what the assistant is doing:

- `[using list_stubs...]` — scanning your document index
- `[using read_raw_file: path...]` — reading a specific document
- `[done]` — that step finished
- `[error: ...]` — something went wrong (the assistant will usually try another approach)

### Why didn't it find information I know is in a document?

Check that the document has a stub in `formatted/`. If the document is only in `raw/` but hasn't been processed by backoffice mode yet, reception mode won't know about it.

### Can I ask follow-up questions?

Yes. The assistant remembers your conversation within a single session. You can say things like "tell me more about that" or "what about the other one?"

### Does it remember my previous sessions?

No. Each time you start the tool, it starts fresh. Your documents are still there, but the conversation history resets.

### What file formats can the assistant write?

The assistant can write text-based files — `.md`, `.csv`, `.json`, and `.txt` — as well as Excel spreadsheets (`.xlsx`). Just ask it to create a file and it will save it to the `output/` directory.

### Where do output files go?

All output files are saved to the `output/` directory inside your research library. This directory is created automatically by `init`. You can open the files with any compatible application — for example, `.xlsx` files open in Excel or Google Sheets.

---

## Backoffice mode

### What does "processing" a document mean?

It means the assistant reads your raw document and creates a short stub file in `formatted/` that contains:

- A title
- A detailed summary
- Tags and metadata
- A pointer back to the full document

This stub is what reception mode uses to find relevant documents quickly.

### Can I edit the stubs it creates?

Yes. The stubs are just text files in the `formatted/` folder. Open them in any text editor and change whatever you like.

### What happens when I push to GitHub?

The assistant runs `git add`, `git commit`, and `git push` for the `formatted/` folder. This uploads your new stubs to the team's shared repository so everyone has access.

### What file formats are supported?

The tool works best with markdown (`.md`) files. Plain text files also work. It reads the file as text, so any text-based format is fine.

---

## Troubleshooting

### The assistant is slow

The first question in a session can take a few seconds because the AI needs to scan your document library. Follow-up questions are usually faster.

If it's consistently slow, you may be using a larger AI model. The default (Claude Sonnet) is a good balance of speed and quality.

### I'm getting API errors

Common causes:

- **Rate limiting** — You're sending too many requests. Wait a moment and try again.
- **Invalid API key** — Double-check your `.env` file.
- **Insufficient credits** — Check your balance at [console.anthropic.com](https://console.anthropic.com/).

### I have a folder of documents but no raw/ or formatted/ — can I use init?

Yes. Run `init` on the folder and it will automatically move your existing files into `raw/` and create the research structure:

```bash
research-assistant init ~/Downloads/my-documents
```

You'll see a message like "Moved 15 existing files into raw/".

### "unknown repo" error

The `--repo` name doesn't match any registered library. Run `research-assistant repos` to see available names.

### The assistant gave a wrong answer

The assistant can only answer based on the documents in your library. If a document contains incorrect information, the assistant will repeat it. You can always ask "which documents did you use?" to check its sources.
