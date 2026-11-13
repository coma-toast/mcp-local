package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestServer_ProtocolAndTools(t *testing.T) {
	srv := New("test-srv", "1.0.0")
	srv.Register(ToolDef{
		Name:        "echo",
		Description: "Echo input",
		InputSchema: map[string]any{"type": "object"},
	}, func(ctx context.Context, args json.RawMessage) (string, error) {
		var p struct {
			Msg string `json:"msg"`
		}
		if err := json.Unmarshal(args, &p); err != nil {
			return "", err
		}
		if p.Msg == "fail" {
			return "", errors.New("boom")
		}
		return "echo:" + p.Msg, nil
	})

	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize"}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"ping"}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"echo","arguments":{"msg":"hi"}}}`,
		`{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"echo","arguments":{"msg":"fail"}}}`,
		`{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"missing"}}`,
		`{"jsonrpc":"2.0","id":7,"method":"nope"}`,
		`{bad json`,
	}, "\n") + "\n"

	var out bytes.Buffer
	if err := srv.Serve(context.Background(), strings.NewReader(input), &out); err != nil {
		t.Fatalf("Serve: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 8 {
		t.Fatalf("got %d responses, want 8:\n%s", len(lines), out.String())
	}
	checks := []string{
		`"protocolVersion":"2024-11-05"`,
		`"result":{}`,
		`"name":"echo"`,
		`"text":"echo:hi"`,
		`"isError":true`,
		`Unknown tool: missing`,
		`Method not found: nope`,
		`-32700`,
	}
	for i, want := range checks {
		if !strings.Contains(lines[i], want) {
			t.Errorf("line %d = %s, want substring %q", i, lines[i], want)
		}
	}
}
