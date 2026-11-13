package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func callControlMCP(t *testing.T, lines ...string) []map[string]any {
	t.Helper()
	srv := newControlMCPServer()
	var out bytes.Buffer
	in := strings.NewReader(strings.Join(lines, "\n") + "\n")
	if err := srv.Serve(context.Background(), in, &out); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	var resps []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("unmarshal %q: %v", line, err)
		}
		resps = append(resps, m)
	}
	return resps
}

func TestControlMCP_EndToEnd(t *testing.T) {
	tempHome(t, "agents:\n  opencode: false\n  cursor: true\n  claude: false\nservices: []\n")

	resps := callControlMCP(t,
		`{"jsonrpc":"2.0","id":1,"method":"initialize"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"add_service","arguments":{"name":"ast","type":"http","mcp_url":"http://127.0.0.1:7821/mcp","port":7821}}}`,
		`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"list_services","arguments":{}}}`,
		`{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"configure_tools","arguments":{"name":"ast","active_tier":"complete","code_mode":true,"tool":"execute_code","enabled":true,"tier":"complete"}}}`,
		`{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"list_tools","arguments":{"name":"ast"}}}`,
		`{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"register_service","arguments":{"name":"ast"}}}`,
		`{"jsonrpc":"2.0","id":8,"method":"tools/call","params":{"name":"service_status","arguments":{"name":"ast"}}}`,
		`{"jsonrpc":"2.0","id":9,"method":"tools/call","params":{"name":"remove_service","arguments":{"name":"ast"}}}`,
	)

	if len(resps) != 9 {
		t.Fatalf("got %d responses, want 9", len(resps))
	}
	for i, r := range resps {
		if r["error"] != nil {
			t.Fatalf("resp %d returned error: %v", i+1, r["error"])
		}
		if res, ok := r["result"].(map[string]any); ok {
			if isErr, _ := res["isError"].(bool); isErr {
				t.Fatalf("resp %d returned tool isError: %v", i+1, res)
			}
		}
	}

	toolText := func(idx int) string {
		res := resps[idx]["result"].(map[string]any)
		content := res["content"].([]any)[0].(map[string]any)
		return content["text"].(string)
	}

	if !strings.Contains(toolText(3), `"name": "ast"`) {
		t.Errorf("list_services output = %s", toolText(3))
	}
	if !strings.Contains(toolText(5), `"active_tier": "complete"`) || !strings.Contains(toolText(5), `"execute_code"`) {
		t.Errorf("list_tools output = %s", toolText(5))
	}
	if !strings.Contains(toolText(7), `"cursor": true`) {
		t.Errorf("service_status output = %s", toolText(7))
	}
	if !strings.Contains(toolText(8), `"removed": "ast"`) {
		t.Errorf("remove_service output = %s", toolText(8))
	}
}
