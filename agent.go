package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
)

type ToolHandler func(name string, input json.RawMessage) (string, error)

type ModeSetup struct {
	SystemPrompt string
	Tools        []anthropic.ToolUnionParam
	HandleTool   ToolHandler
}

// sendMessageStreaming calls the streaming API and prints text deltas to stdout
// as they arrive. It returns the fully accumulated Message for conversation history.
func sendMessageStreaming(client anthropic.Client, model string, setup ModeSetup, messages []anthropic.MessageParam) (*anthropic.Message, error) {
	stream := client.Messages.NewStreaming(context.Background(), anthropic.MessageNewParams{
		Model:     model,
		MaxTokens: 4096,
		System: []anthropic.TextBlockParam{
			{Text: setup.SystemPrompt},
		},
		Messages: messages,
		Tools:    setup.Tools,
	})
	defer stream.Close()

	var message anthropic.Message
	for stream.Next() {
		event := stream.Current()
		message.Accumulate(event)

		switch evt := event.AsAny().(type) {
		case anthropic.ContentBlockDeltaEvent:
			switch delta := evt.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				fmt.Fprint(os.Stdout, delta.Text)
			}
		}
	}

	if err := stream.Err(); err != nil {
		return nil, err
	}

	return &message, nil
}

// toolStatusLabel returns a human-readable label for a tool invocation.
// It always includes the tool name; for tools with a short, useful input
// (e.g. a file path or search query) it appends that context.
func toolStatusLabel(name string, input json.RawMessage) string {
	var params map[string]any
	if err := json.Unmarshal(input, &params); err == nil {
		for _, key := range []string{"path", "filename", "query"} {
			if v, ok := params[key]; ok {
				if s, ok := v.(string); ok && len(s) <= 120 {
					return fmt.Sprintf("using %s: %s...", name, s)
				}
			}
		}
	}
	return fmt.Sprintf("using %s...", name)
}

func runToolLoop(client anthropic.Client, model string, setup ModeSetup, messages *[]anthropic.MessageParam) (string, error) {
	for {
		resp, err := sendMessageStreaming(client, model, setup, *messages)
		if err != nil {
			return "", fmt.Errorf("api error: %w", err)
		}

		*messages = append(*messages, resp.ToParam())

		if resp.StopReason != "tool_use" {
			var text string
			for _, block := range resp.Content {
				if block.Type == "text" {
					text += block.AsText().Text
				}
			}
			return text, nil
		}

		var toolResults []anthropic.ContentBlockParamUnion
		for _, block := range resp.Content {
			if block.Type != "tool_use" {
				continue
			}
			tu := block.AsToolUse()
			fmt.Fprintf(os.Stderr, "[%s]\n", toolStatusLabel(tu.Name, tu.Input))
			result, toolErr := setup.HandleTool(tu.Name, tu.Input)

			if toolErr != nil {
				fmt.Fprintf(os.Stderr, "[error: %s]\n", toolErr.Error())
				toolResults = append(toolResults, anthropic.ContentBlockParamUnion{
					OfToolResult: &anthropic.ToolResultBlockParam{
						ToolUseID: tu.ID,
						IsError:   anthropic.Bool(true),
						Content: []anthropic.ToolResultBlockParamContentUnion{
							{OfText: &anthropic.TextBlockParam{Text: toolErr.Error()}},
						},
					},
				})
			} else {
				fmt.Fprintf(os.Stderr, "[done]\n")
				toolResults = append(toolResults, anthropic.ContentBlockParamUnion{
					OfToolResult: &anthropic.ToolResultBlockParam{
						ToolUseID: tu.ID,
						Content: []anthropic.ToolResultBlockParamContentUnion{
							{OfText: &anthropic.TextBlockParam{Text: result}},
						},
					},
				})
			}
		}

		*messages = append(*messages, anthropic.MessageParam{
			Role:    anthropic.MessageParamRoleUser,
			Content: toolResults,
		})
	}
}
