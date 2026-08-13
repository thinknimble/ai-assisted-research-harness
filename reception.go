package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/xuri/excelize/v2"
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
- You can read multiple raw files if needed to answer comprehensively.
- Use write_spreadsheet to export tabular data as an Excel file — pass column headers and row data as JSON arrays
- Use write_text_file to save research summaries, data extractions, or structured output`

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
				Name:        "write_spreadsheet",
				Description: anthropic.String("Create an Excel (.xlsx) spreadsheet file with the given headers and row data."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"filename": map[string]any{
							"type":        "string",
							"description": "Output filename (must end in .xlsx)",
						},
						"headers": map[string]any{
							"type":        "array",
							"items":       map[string]any{"type": "string"},
							"description": "Column header labels for the first row",
						},
						"rows": map[string]any{
							"type":        "array",
							"description": "Data rows — each inner array is one row of cell values",
							"items": map[string]any{
								"type":  "array",
								"items": map[string]any{"type": "string"},
							},
						},
					},
					Required: []string{"filename", "headers", "rows"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "write_text_file",
				Description: anthropic.String("Write a plain-text file (Markdown, CSV, JSON, or TXT) to the output directory."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"filename": map[string]any{
							"type":        "string",
							"description": "Output filename (must end in .md, .csv, .json, or .txt)",
						},
						"content": map[string]any{
							"type":        "string",
							"description": "Text content to write to the file",
						},
					},
					Required: []string{"filename", "content"},
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

	case "write_text_file":
		return handleWriteTextFile(input)

	case "write_spreadsheet":
		return handleWriteSpreadsheet(input)

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

var allowedTextExtensions = map[string]bool{
	".md":   true,
	".csv":  true,
	".json": true,
	".txt":  true,
}

func handleWriteTextFile(input json.RawMessage) (string, error) {
	var params struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(params.Filename))
	if !allowedTextExtensions[ext] {
		return fmt.Sprintf("Unsupported extension %q — allowed types: .md, .csv, .json, .txt", ext), nil
	}

	if err := writeFile("output", params.Filename, params.Content); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Written output/%s", filepath.Base(params.Filename)), nil
}

func handleWriteSpreadsheet(input json.RawMessage) (string, error) {
	var params struct {
		Filename string     `json:"filename"`
		Headers  []string   `json:"headers"`
		Rows     [][]string `json:"rows"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	if !strings.HasSuffix(params.Filename, ".xlsx") {
		return "", fmt.Errorf("filename must end in .xlsx")
	}

	f := excelize.NewFile()
	defer f.Close()

	// Write headers to row 1
	for i, h := range params.Headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, h)
	}

	// Write data rows starting at row 2
	for rowIdx, row := range params.Rows {
		for colIdx, val := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue("Sheet1", cell, val)
		}
	}

	dest := filepath.Join(projectRoot, "output", filepath.Base(params.Filename))
	if err := f.SaveAs(dest); err != nil {
		return "", fmt.Errorf("failed to save xlsx: %w", err)
	}

	return fmt.Sprintf("Written output/%s (%d rows)", filepath.Base(params.Filename), len(params.Rows)), nil
}
