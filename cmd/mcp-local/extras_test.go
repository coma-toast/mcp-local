package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
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

func TestValidateDashboardURL(t *testing.T) {
	valid := []string{
		"http://localhost:7830",
		"https://example.com/dashboard",
		"http://127.0.0.1:8080/ui?tab=1",
	}
	for _, u := range valid {
		if err := validateDashboardURL(u); err != nil {
			t.Errorf("validateDashboardURL(%q) = %v, want nil", u, err)
		}
	}
	invalid := []string{
		"",
		"file:///etc/passwd",
		"javascript:alert(1)",
		"http://",
		"://missing-scheme",
	}
	for _, u := range invalid {
		if err := validateDashboardURL(u); err == nil {
			t.Errorf("validateDashboardURL(%q) = nil, want error", u)
		}
	}
}

func TestEditorCommand(t *testing.T) {
	tests := []struct {
		editor string
		want   []string
	}{
		{"", []string{"nano", "/tmp/config.yaml"}},
		{"   ", []string{"nano", "/tmp/config.yaml"}},
		{"vim", []string{"vim", "/tmp/config.yaml"}},
		{"code --wait", []string{"code", "--wait", "/tmp/config.yaml"}},
		{"nvim -u NONE", []string{"nvim", "-u", "NONE", "/tmp/config.yaml"}},
	}
	for _, tt := range tests {
		cmd := editorCommand(tt.editor, "/tmp/config.yaml")
		if !reflect.DeepEqual(cmd.Args, tt.want) {
			t.Errorf("editorCommand(%q) args = %v, want %v", tt.editor, cmd.Args, tt.want)
		}
	}
}

func TestCmdRebuild_UsesCommandDir(t *testing.T) {
	svcDir := filepath.Join(t.TempDir(), "svc")
	if err := os.MkdirAll(svcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	binPath := filepath.Join(svcDir, "my-server")
	if err := os.WriteFile(binPath, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	cfgYAML := fmt.Sprintf("agents:\n  opencode: false\n  cursor: false\n  claude: false\nservices:\n  - name: s1\n    command: %s\n    build_command: pwd > built.cwd\n", binPath)
	tempHome(t, cfgYAML)

	if _, err := runCmd(t, cmdRebuild(), "s1"); err != nil {
		t.Fatalf("rebuild failed: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(svcDir, "built.cwd"))
	if err != nil {
		t.Fatalf("expected built.cwd in %s: %v", svcDir, err)
	}
	if filepath.Clean(string(got[:len(got)-1])) != filepath.Clean(svcDir) {
		t.Errorf("build ran in %q, want %q", string(got), svcDir)
	}
}
