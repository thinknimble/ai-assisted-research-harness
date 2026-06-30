package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
