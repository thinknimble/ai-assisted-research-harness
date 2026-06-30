package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//go:embed doc-template.yaml
var embeddedDocTemplate string

func runInit(args []string, stdout io.Writer) error {
	target := "."
	if len(args) > 0 && args[0] != "" {
		target = args[0]
	}

	absTarget, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check if directory exists and has content
	entries, err := os.ReadDir(absTarget)
	if err == nil && len(entries) > 0 {
		return fmt.Errorf("directory %s already has content — aborting to avoid overwriting", absTarget)
	}

	// Create directory structure
	if err := os.MkdirAll(filepath.Join(absTarget, "raw"), 0755); err != nil {
		return fmt.Errorf("failed to create raw/: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(absTarget, "formatted"), 0755); err != nil {
		return fmt.Errorf("failed to create formatted/: %w", err)
	}

	// Write doc-template.yaml
	templatePath := filepath.Join(absTarget, "doc-template.yaml")
	if err := os.WriteFile(templatePath, []byte(embeddedDocTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write doc-template.yaml: %w", err)
	}

	fmt.Fprintf(stdout, "Created research directory at %s\n", absTarget)
	fmt.Fprintln(stdout, "  raw/")
	fmt.Fprintln(stdout, "  formatted/")
	fmt.Fprintln(stdout, "  doc-template.yaml")
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "Next: set ANTHROPIC_API_KEY in your environment, then run:")
	fmt.Fprintf(stdout, "  research-assistant --mode backoffice\n")

	return nil
}
