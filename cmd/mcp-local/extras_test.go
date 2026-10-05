package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFetchToolsList_StreamableHTTP(t *testing.T) {
	var notified bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			return
		}
		if r.Header.Get("Accept") != "application/json, text/event-stream" {
			t.Errorf("Accept = %q", r.Header.Get("Accept"))
		}
		body, _ := io.ReadAll(r.Body)
		var raw map[string]json.RawMessage
		_ = json.Unmarshal(body, &raw)
		var method string
		_ = json.Unmarshal(raw["method"], &method)
		if method != "initialize" && (r.Header.Get("Mcp-Session-Id") != "s1" || r.Header.Get("MCP-Protocol-Version") != "2025-06-18") {
			t.Errorf("%s: session=%q version=%q", method, r.Header.Get("Mcp-Session-Id"), r.Header.Get("MCP-Protocol-Version"))
		}
		switch method {
		case "initialize":
			w.Header().Set("Mcp-Session-Id", "s1")
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":{"protocolVersion":"2025-06-18"}}`, raw["id"])
		case "notifications/initialized":
			if _, has := raw["id"]; has {
				t.Errorf("notification carries an id: %s", body)
			}
			notified = true
			w.WriteHeader(http.StatusAccepted)
		case "tools/list":
			w.Header().Set("Content-Type", "text/event-stream")
			fmt.Fprintf(w, "event: message\ndata: {\"jsonrpc\":\"2.0\",\"id\":%s,\"result\":{\"tools\":[{\"name\":\"handoff\",\"description\":\"d\"}]}}\n\n", raw["id"])
		}
	}))
	defer ts.Close()
	tools, err := fetchToolsList(ts.URL)
	if err != nil {
		t.Fatalf("fetchToolsList: %v", err)
	}
	if !notified || len(tools) != 1 || tools[0].Name != "handoff" || !tools[0].Enabled {
		t.Errorf("notified=%v tools=%+v", notified, tools)
	}
}

func TestIsRegisteredInFile_JSONC(t *testing.T) {
	path := filepath.Join(t.TempDir(), "opencode.jsonc")
	jsonc := "{\n  // comment\n  \"mcp\": {\n    \"svc\": {\"type\": \"remote\"}, /* block */\n  },\n}\n"
	if err := os.WriteFile(path, []byte(jsonc), 0o644); err != nil {
		t.Fatal(err)
	}
	if !isRegisteredInFile(path, "svc") || isRegisteredInFile(path, "other") {
		t.Error("JSONC config not parsed")
	}
}

func TestBridgeURL(t *testing.T) {
	entry := map[string]interface{}{"command": "/usr/local/bin/mcp-local", "args": []interface{}{"bridge", "http://x/mcp"}}
	if url, ok := bridgeURL(entry); !ok || url != "http://x/mcp" {
		t.Errorf("bridgeURL = %q, %v", url, ok)
	}
	if _, ok := bridgeURL(map[string]interface{}{"command": "npx", "args": []interface{}{"bridge", "x"}}); ok {
		t.Error("non-mcp-local command matched")
	}
}
