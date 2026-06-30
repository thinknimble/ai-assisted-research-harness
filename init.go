package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//go:embed doc-template.yaml
var embeddedDocTemplate string

func runInit(args []string, stdout io.Writer, readKey func() (string, error)) error {
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

	// Prompt for API key
	fmt.Fprint(stdout, "Enter your Anthropic API key: ")
	apiKey, err := readKey()
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}
	fmt.Fprintln(stdout) // newline after masked input

	apiKey = strings.TrimSpace(apiKey)
	envPath := filepath.Join(absTarget, ".env")

	if apiKey == "" {
		placeholder := "# Get your API key at https://console.anthropic.com/settings/keys\n# Uncomment the line below and paste your key:\n# ANTHROPIC_API_KEY=your-key-here\n"
		if err := os.WriteFile(envPath, []byte(placeholder), 0644); err != nil {
			return fmt.Errorf("failed to write .env: %w", err)
		}
	} else {
		content := "ANTHROPIC_API_KEY=" + apiKey + "\n"
		if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write .env: %w", err)
		}
	}

	// Register repo in global config
	cfg, err := LoadGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	repoName := RegisterRepo(cfg, absTarget)
	if err := SaveGlobalConfig(cfg); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}

	fmt.Fprintf(stdout, "Created research directory at %s\n", absTarget)
	fmt.Fprintln(stdout, "  raw/")
	fmt.Fprintln(stdout, "  formatted/")
	fmt.Fprintln(stdout, "  doc-template.yaml")
	fmt.Fprintln(stdout, "  .env")
	fmt.Fprintf(stdout, "  registered as %q in global config\n", repoName)
	fmt.Fprintln(stdout, "")
	if apiKey == "" {
		fmt.Fprintln(stdout, "Next: add your ANTHROPIC_API_KEY to .env, then run:")
	} else {
		fmt.Fprintln(stdout, "Next: run:")
	}
	fmt.Fprintf(stdout, "  research-assistant --mode backoffice\n")

	return nil
}
