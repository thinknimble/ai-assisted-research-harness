package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

const backofficeSystemPrompt = `You are a document formatting assistant for a research knowledge base.

Your job is to process raw research documents into formatted reference stubs.

## Schema

The frontmatter schema is defined in doc-template.yaml. Use the read_file tool to read it if you need a refresher. Key rules:
- Every stub must have: title, type, status, created, updated, author, source, source_url, path, tags, related, summary
- The summary field is critical — it must be detailed enough that an LLM can decide whether to read the full document based on the summary alone
- The path field must point to the raw file's relative location (e.g. raw/survey-webhooks/typeform-webhooks.md)
- The formatted stub contains ONLY frontmatter (between --- delimiters), no body content

## Workflow

1. Use list_raw_files to see available raw documents
2. Use list_formatted_files to see what has already been processed
3. Read raw files to understand their content
4. Create formatted stubs with write_formatted_stub
5. Use git_add_commit_push to commit and push when the user is satisfied

## Important

- Write rich, detailed summaries. Include key entities, concepts, and data points. A sparse summary defeats the purpose.
- Set the author to "pari" unless told otherwise.
- Use today's date for created/updated fields.
- Infer the type from the content (reference, paper, article, code, notes, etc.)`

func setupBackoffice() ModeSetup {
	tools := []anthropic.ToolUnionParam{
		{
			OfTool: &anthropic.ToolParam{
				Name:        "list_raw_files",
				Description: anthropic.String("List all files in the raw/ directory recursively."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "read_file",
				Description: anthropic.String("Read the contents of a file. Path must be relative to the project root."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"path": map[string]any{
							"type":        "string",
							"description": "File path to read (relative to project root)",
						},
					},
					Required: []string{"path"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "list_formatted_files",
				Description: anthropic.String("List all files in the formatted/ directory."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "write_formatted_stub",
				Description: anthropic.String("Write a formatted document stub to the formatted/ directory. Content should be YAML frontmatter only (between --- delimiters)."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"filename": map[string]any{
							"type":        "string",
							"description": "Filename (e.g. 'my-doc.md'). Written to formatted/{filename}.",
						},
						"content": map[string]any{
							"type":        "string",
							"description": "Full file content (frontmatter block)",
						},
					},
					Required: []string{"filename", "content"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "git_add_commit_push",
				Description: anthropic.String("Stage all changes in formatted/, commit with the given message, and push to origin."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"message": map[string]any{
							"type":        "string",
							"description": "Git commit message",
						},
					},
					Required: []string{"message"},
				},
			},
		},
	}

	return ModeSetup{
		SystemPrompt: backofficeSystemPrompt,
		Tools:        tools,
		HandleTool:   handleBackofficeTool,
	}
}

func handleBackofficeTool(name string, input json.RawMessage) (string, error) {
	switch name {
	case "list_raw_files":
		files, err := listDir("raw")
		if err != nil {
			return "", fmt.Errorf("failed to list raw/: %w", err)
		}
		if len(files) == 0 {
			return "No raw files found.", nil
		}
		return strings.Join(files, "\n"), nil

	case "read_file":
		var params struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal(input, &params); err != nil {
			return "", fmt.Errorf("invalid input: %w", err)
		}
		content, err := readFileContent(params.Path)
		if err != nil {
			return "", fmt.Errorf("failed to read %s: %w", params.Path, err)
		}
		return content, nil

	case "list_formatted_files":
		files, err := listDir("formatted")
		if err != nil {
			return "", fmt.Errorf("failed to list formatted/: %w", err)
		}
		if len(files) == 0 {
			return "No formatted files found.", nil
		}
		return strings.Join(files, "\n"), nil

	case "write_formatted_stub":
		var params struct {
			Filename string `json:"filename"`
			Content  string `json:"content"`
		}
		if err := json.Unmarshal(input, &params); err != nil {
			return "", fmt.Errorf("invalid input: %w", err)
		}
		if err := writeFile("formatted", params.Filename, params.Content); err != nil {
			return "", fmt.Errorf("failed to write formatted/%s: %w", params.Filename, err)
		}
		return fmt.Sprintf("Written formatted/%s", params.Filename), nil

	case "git_add_commit_push":
		var params struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(input, &params); err != nil {
			return "", fmt.Errorf("invalid input: %w", err)
		}
		return gitAddCommitPush(params.Message)

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func gitAddCommitPush(message string) (string, error) {
	cmds := []struct {
		name string
		args []string
	}{
		{"git", []string{"add", "formatted/"}},
		{"git", []string{"commit", "-m", message}},
		{"git", []string{"push"}},
	}

	var output strings.Builder
	for _, c := range cmds {
		cmd := exec.Command(c.name, c.args...)
		cmd.Dir = projectRoot
		out, err := cmd.CombinedOutput()
		output.WriteString(fmt.Sprintf("$ %s %s\n%s\n", c.name, strings.Join(c.args, " "), string(out)))
		if err != nil {
			return output.String(), fmt.Errorf("%s failed: %w", c.name, err)
		}
	}
	return output.String(), nil
}
