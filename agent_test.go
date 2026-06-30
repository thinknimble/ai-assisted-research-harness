package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// sseResponse builds a valid SSE byte stream that simulates a streaming
// text response with the given chunks, ending with stop_reason "end_turn".
func sseResponse(chunks []string) []byte {
	var buf bytes.Buffer

	buf.WriteString("event: message_start\n")
	buf.WriteString(`data: {"type":"message_start","message":{"id":"msg_test","type":"message","role":"assistant","content":[],"model":"claude-sonnet-4-6-20250514","stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":10,"output_tokens":0}}}`)
	buf.WriteString("\n\n")

	buf.WriteString("event: content_block_start\n")
	buf.WriteString(`data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`)
	buf.WriteString("\n\n")

	for _, chunk := range chunks {
		buf.WriteString("event: content_block_delta\n")
		fmt.Fprintf(&buf, `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":%q}}`, chunk)
		buf.WriteString("\n\n")
	}

	buf.WriteString("event: content_block_stop\n")
	buf.WriteString(`data: {"type":"content_block_stop","index":0}`)
	buf.WriteString("\n\n")

	buf.WriteString("event: message_delta\n")
	buf.WriteString(`data: {"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null},"usage":{"output_tokens":5}}`)
	buf.WriteString("\n\n")

	buf.WriteString("event: message_stop\n")
	buf.WriteString(`data: {"type":"message_stop"}`)
	buf.WriteString("\n\n")

	return buf.Bytes()
}

func TestStreamingPrintsTextToStdout(t *testing.T) {
	chunks := []string{"Hello", ", ", "world", "!"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
		w.Write(sseResponse(chunks))
	}))
	defer server.Close()

	client := anthropic.NewClient(
		option.WithAPIKey("test-key"),
		option.WithBaseURL(server.URL),
	)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	setup := ModeSetup{
		SystemPrompt: "test",
		Tools:        nil,
		HandleTool:   nil,
	}
	messages := []anthropic.MessageParam{
		{
			Role: anthropic.MessageParamRoleUser,
			Content: []anthropic.ContentBlockParamUnion{
				{OfText: &anthropic.TextBlockParam{Text: "hi"}},
			},
		},
	}

	text, err := runToolLoop(client, "claude-sonnet-4-6-20250514", setup, &messages)

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify returned text matches full concatenation
	expected := "Hello, world!"
	if text != expected {
		t.Errorf("returned text = %q, want %q", text, expected)
	}

	// Verify text was printed to stdout (streamed)
	var captured bytes.Buffer
	captured.ReadFrom(r)
	if captured.String() != expected {
		t.Errorf("stdout output = %q, want %q", captured.String(), expected)
	}

	// Verify message was appended to conversation history
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages in history, got %d", len(messages))
	}
}

// sseToolUseResponse builds an SSE stream that simulates the assistant invoking
// a tool (stop_reason "tool_use").
func sseToolUseResponse(toolName, toolID string, toolInput json.RawMessage) []byte {
	var buf bytes.Buffer

	buf.WriteString("event: message_start\n")
	buf.WriteString(`data: {"type":"message_start","message":{"id":"msg_test2","type":"message","role":"assistant","content":[],"model":"claude-sonnet-4-6-20250514","stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":10,"output_tokens":0}}}`)
	buf.WriteString("\n\n")

	// tool_use content block
	buf.WriteString("event: content_block_start\n")
	fmt.Fprintf(&buf, `data: {"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":%q,"name":%q,"input":{}}}`, toolID, toolName)
	buf.WriteString("\n\n")

	// Stream the input JSON via input_json_delta
	inputStr := string(toolInput)
	buf.WriteString("event: content_block_delta\n")
	fmt.Fprintf(&buf, `data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":%q}}`, inputStr)
	buf.WriteString("\n\n")

	buf.WriteString("event: content_block_stop\n")
	buf.WriteString(`data: {"type":"content_block_stop","index":0}`)
	buf.WriteString("\n\n")

	buf.WriteString("event: message_delta\n")
	buf.WriteString(`data: {"type":"message_delta","delta":{"stop_reason":"tool_use","stop_sequence":null},"usage":{"output_tokens":5}}`)
	buf.WriteString("\n\n")

	buf.WriteString("event: message_stop\n")
	buf.WriteString(`data: {"type":"message_stop"}`)
	buf.WriteString("\n\n")

	return buf.Bytes()
}

func TestToolCallPrintsStatusToStderr(t *testing.T) {
	toolInput := json.RawMessage(`{"path":"raw/survey-webhooks/typeform.md"}`)
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
		callCount++
		if callCount == 1 {
			// First call: assistant invokes a tool
			w.Write(sseToolUseResponse("read_file", "toolu_01", toolInput))
		} else {
			// Second call: assistant responds with text after tool result
			w.Write(sseResponse([]string{"Done."}))
		}
	}))
	defer server.Close()

	client := anthropic.NewClient(
		option.WithAPIKey("test-key"),
		option.WithBaseURL(server.URL),
	)

	// Capture stderr
	oldStderr := os.Stderr
	stderrR, stderrW, _ := os.Pipe()
	os.Stderr = stderrW

	// Also capture stdout so it doesn't leak
	oldStdout := os.Stdout
	_, stdoutW, _ := os.Pipe()
	os.Stdout = stdoutW

	setup := ModeSetup{
		SystemPrompt: "test",
		Tools: []anthropic.ToolUnionParam{
			{OfTool: &anthropic.ToolParam{
				Name:        "read_file",
				Description: anthropic.String("Read a file"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"path": map[string]any{"type": "string"},
					},
				},
			}},
		},
		HandleTool: func(name string, input json.RawMessage) (string, error) {
			return "file contents here", nil
		},
	}
	messages := []anthropic.MessageParam{
		{
			Role: anthropic.MessageParamRoleUser,
			Content: []anthropic.ContentBlockParamUnion{
				{OfText: &anthropic.TextBlockParam{Text: "read the file"}},
			},
		},
	}

	_, err := runToolLoop(client, "claude-sonnet-4-6-20250514", setup, &messages)

	stderrW.Close()
	stdoutW.Close()
	os.Stderr = oldStderr
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var stderrBuf bytes.Buffer
	stderrBuf.ReadFrom(stderrR)
	stderrOutput := stderrBuf.String()

	// Must contain the tool name and path context
	if !strings.Contains(stderrOutput, "read_file") {
		t.Errorf("stderr should contain tool name 'read_file', got: %q", stderrOutput)
	}
	if !strings.Contains(stderrOutput, "raw/survey-webhooks/typeform.md") {
		t.Errorf("stderr should contain path context, got: %q", stderrOutput)
	}
}

func TestToolStatusLabel(t *testing.T) {
	tests := []struct {
		name     string
		tool     string
		input    string
		wantSub  string // substring that must appear
	}{
		{"no input", "list_stubs", `{}`, "using list_stubs..."},
		{"with path", "read_file", `{"path":"raw/doc.md"}`, "using read_file: raw/doc.md..."},
		{"with query", "search", `{"query":"webhooks"}`, "using search: webhooks..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toolStatusLabel(tt.tool, json.RawMessage(tt.input))
			if got != tt.wantSub {
				t.Errorf("toolStatusLabel(%q, %s) = %q, want %q", tt.tool, tt.input, got, tt.wantSub)
			}
		})
	}
}
