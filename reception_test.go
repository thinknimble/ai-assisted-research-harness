package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestWriteSpreadsheetTool(t *testing.T) {
	// Set up a temp project root with output/
	tmp := t.TempDir()
	origRoot := projectRoot
	projectRoot = tmp
	t.Cleanup(func() { projectRoot = origRoot })

	if err := os.MkdirAll(filepath.Join(tmp, "output"), 0755); err != nil {
		t.Fatal(err)
	}

	t.Run("happy path creates xlsx with headers and rows", func(t *testing.T) {
		input := map[string]any{
			"filename": "report.xlsx",
			"headers":  []string{"Name", "Age", "City"},
			"rows": [][]string{
				{"Alice", "30", "NYC"},
				{"Bob", "25", "LA"},
			},
		}
		raw, _ := json.Marshal(input)

		result, err := handleReceptionTool("write_spreadsheet", raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != "Written output/report.xlsx (2 rows)" {
			t.Errorf("got %q, want %q", result, "Written output/report.xlsx (2 rows)")
		}

		// Verify the file is valid xlsx
		f, err := excelize.OpenFile(filepath.Join(tmp, "output", "report.xlsx"))
		if err != nil {
			t.Fatalf("failed to open created xlsx: %v", err)
		}
		defer f.Close()

		// Check headers in row 1
		for i, hdr := range []string{"Name", "Age", "City"} {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			val, _ := f.GetCellValue("Sheet1", cell)
			if val != hdr {
				t.Errorf("header cell %s = %q, want %q", cell, val, hdr)
			}
		}

		// Check data row 1
		val, _ := f.GetCellValue("Sheet1", "A2")
		if val != "Alice" {
			t.Errorf("A2 = %q, want Alice", val)
		}
		val, _ = f.GetCellValue("Sheet1", "B3")
		if val != "25" {
			t.Errorf("B3 = %q, want 25", val)
		}
	})

	t.Run("rejects non-xlsx extension", func(t *testing.T) {
		input := map[string]any{
			"filename": "report.csv",
			"headers":  []string{"A"},
			"rows":     [][]string{{"1"}},
		}
		raw, _ := json.Marshal(input)

		_, err := handleReceptionTool("write_spreadsheet", raw)
		if err == nil {
			t.Fatal("expected error for non-xlsx extension")
		}
		if !strings.Contains(err.Error(), ".xlsx") {
			t.Errorf("error should mention .xlsx, got: %v", err)
		}
	})

	t.Run("zero rows produces valid file", func(t *testing.T) {
		input := map[string]any{
			"filename": "empty.xlsx",
			"headers":  []string{"Col1"},
			"rows":     [][]string{},
		}
		raw, _ := json.Marshal(input)

		result, err := handleReceptionTool("write_spreadsheet", raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "Written output/empty.xlsx (0 rows)" {
			t.Errorf("got %q", result)
		}
	})
}

func TestWriteTextFileTool(t *testing.T) {
	tmp := t.TempDir()
	origRoot := projectRoot
	projectRoot = tmp
	t.Cleanup(func() { projectRoot = origRoot })

	if err := os.MkdirAll(filepath.Join(tmp, "output"), 0755); err != nil {
		t.Fatal(err)
	}

	t.Run("writes markdown file and returns confirmation", func(t *testing.T) {
		input := map[string]any{
			"filename": "summary.md",
			"content":  "# Hello\nWorld",
		}
		raw, _ := json.Marshal(input)

		result, err := handleReceptionTool("write_text_file", raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "Written output/summary.md" {
			t.Errorf("got %q, want %q", result, "Written output/summary.md")
		}

		data, err := os.ReadFile(filepath.Join(tmp, "output", "summary.md"))
		if err != nil {
			t.Fatalf("file not created: %v", err)
		}
		if string(data) != "# Hello\nWorld" {
			t.Errorf("content = %q", string(data))
		}
	})

	t.Run("accepts all allowed extensions", func(t *testing.T) {
		for _, ext := range []string{".csv", ".json", ".txt"} {
			input := map[string]any{
				"filename": "test" + ext,
				"content":  "data",
			}
			raw, _ := json.Marshal(input)
			_, err := handleReceptionTool("write_text_file", raw)
			if err != nil {
				t.Errorf("extension %s should be allowed, got error: %v", ext, err)
			}
		}
	})

	t.Run("rejects disallowed extension with error message", func(t *testing.T) {
		input := map[string]any{
			"filename": "hack.exe",
			"content":  "bad",
		}
		raw, _ := json.Marshal(input)

		result, err := handleReceptionTool("write_text_file", raw)
		if err != nil {
			t.Fatalf("should return error message, not Go error: %v", err)
		}
		if !strings.Contains(result, ".md") || !strings.Contains(result, ".csv") ||
			!strings.Contains(result, ".json") || !strings.Contains(result, ".txt") {
			t.Errorf("rejection message should list allowed types, got: %q", result)
		}
	})
}

func TestSetupReceptionIncludesWriteTextFile(t *testing.T) {
	setup := setupReception()

	found := false
	for _, tool := range setup.Tools {
		if tool.OfTool != nil && tool.OfTool.Name == "write_text_file" {
			found = true
			props, ok := tool.OfTool.InputSchema.Properties.(map[string]any)
			if !ok {
				t.Fatal("InputSchema.Properties is not map[string]any")
			}
			for _, param := range []string{"filename", "content"} {
				if _, exists := props[param]; !exists {
					t.Errorf("missing required parameter: %s", param)
				}
			}
			break
		}
	}
	if !found {
		t.Error("write_text_file tool not found in setupReception()")
	}
}

func TestReceptionSystemPromptMentionsWriteTextFile(t *testing.T) {
	if !strings.Contains(receptionSystemPrompt, "write_text_file") {
		t.Error("system prompt should mention write_text_file")
	}
}

func TestSetupReceptionIncludesWriteSpreadsheet(t *testing.T) {
	setup := setupReception()

	found := false
	for _, tool := range setup.Tools {
		if tool.OfTool != nil && tool.OfTool.Name == "write_spreadsheet" {
			found = true
			// Verify required parameters exist
			props, ok := tool.OfTool.InputSchema.Properties.(map[string]any)
			if !ok {
				t.Fatal("InputSchema.Properties is not map[string]any")
			}
			for _, param := range []string{"filename", "headers", "rows"} {
				if _, exists := props[param]; !exists {
					t.Errorf("missing required parameter: %s", param)
				}
			}
			break
		}
	}
	if !found {
		t.Error("write_spreadsheet tool not found in setupReception()")
	}
}

func TestReceptionSystemPromptMentionsWriteSpreadsheet(t *testing.T) {
	if !strings.Contains(receptionSystemPrompt, "write_spreadsheet") {
		t.Error("system prompt should mention write_spreadsheet")
	}
}
