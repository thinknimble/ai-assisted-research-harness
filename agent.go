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

// healOrphanedToolUse scans the message history for assistant messages that
// contain tool_use blocks without a corresponding tool_result in the next user
// message. For each orphan, a synthetic tool_result (is_error: true) is
// injected so the API call won't fail with a 400.
func healOrphanedToolUse(messages *[]anthropic.MessageParam) {
	msgs := *messages
	for i := 0; i < len(msgs); i++ {
		if msgs[i].Role != anthropic.MessageParamRoleAssistant {
			continue
		}

		// Collect tool_use IDs from this assistant message.
		var toolUseIDs []string
		for _, block := range msgs[i].Content {
			if block.OfToolUse != nil {
				toolUseIDs = append(toolUseIDs, block.OfToolUse.ID)
			}
		}
		if len(toolUseIDs) == 0 {
			continue
		}

		// Check which IDs are answered in the next user message.
		answered := make(map[string]bool)
		if i+1 < len(msgs) && msgs[i+1].Role == anthropic.MessageParamRoleUser {
			for _, block := range msgs[i+1].Content {
				if block.OfToolResult != nil {
					answered[block.OfToolResult.ToolUseID] = true
				}
			}
		}

		// Build synthetic results for orphans.
		var orphanResults []anthropic.ContentBlockParamUnion
		for _, id := range toolUseIDs {
			if !answered[id] {
				orphanResults = append(orphanResults, anthropic.ContentBlockParamUnion{
					OfToolResult: &anthropic.ToolResultBlockParam{
						ToolUseID: id,
						IsError:   anthropic.Bool(true),
						Content: []anthropic.ToolResultBlockParamContentUnion{
							{OfText: &anthropic.TextBlockParam{Text: "Tool call was not completed"}},
						},
					},
				})
			}
		}
		if len(orphanResults) == 0 {
			continue
		}

		// Inject: either append to existing next-user-message or insert a new one.
		if i+1 < len(msgs) && msgs[i+1].Role == anthropic.MessageParamRoleUser {
			msgs[i+1].Content = append(msgs[i+1].Content, orphanResults...)
		} else {
			// Insert a new user message after the assistant message.
			newMsg := anthropic.MessageParam{
				Role:    anthropic.MessageParamRoleUser,
				Content: orphanResults,
			}
			msgs = append(msgs[:i+1], append([]anthropic.MessageParam{newMsg}, msgs[i+1:]...)...)
			*messages = msgs
		}
	}
}

func runToolLoop(client anthropic.Client, model string, setup ModeSetup, messages *[]anthropic.MessageParam) (string, error) {
	for {
		healOrphanedToolUse(messages)
		resp, err := sendMessageStreaming(client, model, setup, *messages)
		if err != nil {
			return "", fmt.Errorf("api error: %w", err)
		}

		*messages = append(*messages, resp.ToParam())

		if resp.StopReason != "tool_use" {
			// If the response contains tool_use blocks despite a non-tool_use
			// stop_reason (e.g. max_tokens), inject synthetic error results so
			// the next API call doesn't get a 400 for orphaned tool_use.
			var syntheticResults []anthropic.ContentBlockParamUnion
			for _, block := range resp.Content {
				if block.Type == "tool_use" {
					tu := block.AsToolUse()
					syntheticResults = append(syntheticResults, anthropic.ContentBlockParamUnion{
						OfToolResult: &anthropic.ToolResultBlockParam{
							ToolUseID: tu.ID,
							IsError:   anthropic.Bool(true),
							Content: []anthropic.ToolResultBlockParamContentUnion{
								{OfText: &anthropic.TextBlockParam{Text: "Tool call interrupted — response was truncated"}},
							},
						},
					})
				}
			}
			if len(syntheticResults) > 0 {
				*messages = append(*messages, anthropic.MessageParam{
					Role:    anthropic.MessageParamRoleUser,
					Content: syntheticResults,
				})
				continue // loop back so the LLM can retry or respond
			}

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
