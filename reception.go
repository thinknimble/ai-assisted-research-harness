package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

const receptionSystemPrompt = `You are a research assistant with access to a library of document references.

## How to answer questions

1. First, call list_stubs to scan the available document index.
   Each stub has a summary describing what the full document contains.
2. Identify which documents are relevant based on their summaries.
3. If you need specific details, use read_raw_file to pull in the full document content.
4. Synthesize your answer from the documents you have read.

## Guidelines

- Always ground your answers in the available documents.
- If no documents are relevant, say so honestly.
- Cite which documents you used in your answer.
- You can read multiple raw files if needed to answer comprehensively.`

func setupReception() ModeSetup {
	tools := []anthropic.ToolUnionParam{
		{
			OfTool: &anthropic.ToolParam{
				Name:        "list_stubs",
				Description: anthropic.String("List all formatted document stubs with their full frontmatter (title, summary, tags, path, etc). Use this to scan available knowledge and find relevant documents."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "read_raw_file",
				Description: anthropic.String("Read the full content of a raw document file. Use the 'path' field from a stub's frontmatter to get the file path."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"path": map[string]any{
							"type":        "string",
							"description": "Path to the raw file (from the stub's 'path' field)",
						},
					},
					Required: []string{"path"},
				},
			},
		},
	}

	return ModeSetup{
		SystemPrompt: receptionSystemPrompt,
		Tools:        tools,
		HandleTool:   handleReceptionTool,
	}
}

func handleReceptionTool(name string, input json.RawMessage) (string, error) {
	switch name {
	case "list_stubs":
		files, err := listDir("formatted")
		if err != nil {
			return "", fmt.Errorf("failed to list formatted/: %w", err)
		}
		if len(files) == 0 {
			return "No formatted stubs found.", nil
		}
		var sb strings.Builder
		for _, f := range files {
			content, err := readFileContent(f)
			if err != nil {
				sb.WriteString(fmt.Sprintf("--- %s: error reading: %s\n", f, err))
				continue
			}
			sb.WriteString(fmt.Sprintf("=== %s ===\n%s\n\n", f, content))
		}
		return sb.String(), nil

	case "read_raw_file":
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

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}
