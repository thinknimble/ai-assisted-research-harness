package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"golang.org/x/term"
)

func loadEnv() {
	path := filepath.Join(projectRoot, ".env")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"'`)
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

func main() {
	// Handle subcommands before flag parsing
	if len(os.Args) >= 2 && os.Args[1] == "repos" {
		if err := runRepos(os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	if len(os.Args) >= 2 && os.Args[1] == "init" {
		readKey := func() (string, error) {
			b, err := term.ReadPassword(int(os.Stdin.Fd()))
			return string(b), err
		}
		if err := runInit(os.Args[2:], os.Stdout, readKey); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	mode := flag.String("mode", "", "Mode: 'backoffice' or 'reception'")
	model := flag.String("model", string(anthropic.ModelClaudeSonnet4_6), "Claude model to use")
	repo := flag.String("repo", "", "Target repo name from config")
	flag.Parse()

	if *mode != "backoffice" && *mode != "reception" {
		fmt.Fprintln(os.Stderr, `Usage: research-assistant <command> [options]

Commands:
  init [path]       Initialize a new research directory (default: current dir)
  repos             List registered research repos

Options:
  --mode <backoffice|reception>   Run in the specified mode (required)
  --model <model>                 Claude model to use (default: claude-sonnet-4-6)
  --repo <name>                   Target repo name from config`)
		os.Exit(1)
	}

	root, err := ResolveProjectRoot(*repo)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	projectRoot = root

	// Show which repo is active and how to switch
	repoName := filepath.Base(projectRoot)
	fmt.Printf("Using repo: %s (%s)\n", repoName, projectRoot)
	fmt.Println("  Switch with: research-assistant --repo <name> --mode ...")
	fmt.Println("  List repos:  research-assistant repos")
	fmt.Println()

	loadEnv()

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "ANTHROPIC_API_KEY environment variable is required (set in .env or environment)")
		os.Exit(1)
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	var setup ModeSetup
	switch *mode {
	case "backoffice":
		setup = setupBackoffice()
	case "reception":
		setup = setupReception()
	}

	// Backoffice: process unformatted files one at a time, then exit
	if *mode == "backoffice" {
		unprocessed, err := unprocessedRawFiles()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning files: %s\n", err)
			os.Exit(1)
		}
		if len(unprocessed) == 0 {
			fmt.Println("No new files to process.")
			return
		}

		fmt.Printf("Processing %d new file(s)...\n\n", len(unprocessed))
		failed := 0
		for i, file := range unprocessed {
			fmt.Printf("[%d/%d] %s\n", i+1, len(unprocessed), file)
			prompt := fmt.Sprintf("Read the raw file at %q and create a formatted stub for it.", file)
			messages := []anthropic.MessageParam{{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					{OfText: &anthropic.TextBlockParam{Text: prompt}},
				},
			}}
			_, err := runToolLoop(client, *model, setup, &messages)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", file, err)
				failed++
				continue
			}
			fmt.Println()
		}
		if failed > 0 {
			fmt.Fprintf(os.Stderr, "\n%d file(s) failed to process.\n", failed)
		}
		fmt.Printf("\nDone. Processed %d/%d file(s).\n", len(unprocessed)-failed, len(unprocessed))
		return
	}

	// Reception: interactive chat loop
	fmt.Printf("research-assistant [%s mode]\n", *mode)
	fmt.Println("Type your message, or 'quit' to exit.")
	fmt.Println()

	messages := []anthropic.MessageParam{}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "quit" || input == "exit" {
			break
		}

		messages = append(messages, anthropic.MessageParam{
			Role: anthropic.MessageParamRoleUser,
			Content: []anthropic.ContentBlockParamUnion{
				{OfText: &anthropic.TextBlockParam{Text: input}},
			},
		})

		_, err := runToolLoop(client, *model, setup, &messages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			continue
		}

		fmt.Println()
		fmt.Println()
	}
}
