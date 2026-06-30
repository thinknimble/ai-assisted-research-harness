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
	loadEnv()

	mode := flag.String("mode", "", "Mode: 'backoffice' or 'reception'")
	model := flag.String("model", string(anthropic.ModelClaudeSonnet4_6), "Claude model to use")
	flag.Parse()

	if *mode != "backoffice" && *mode != "reception" {
		fmt.Fprintln(os.Stderr, "Usage: research-assistant --mode <backoffice|reception>")
		os.Exit(1)
	}

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
