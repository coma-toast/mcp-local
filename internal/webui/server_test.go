package webui

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

func setupTestServer(t *testing.T) (*Server, string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)

	assets := fstest.MapFS{
		"index.html": {Data: []byte("<html><body>mcp-local-95</body></html>")},
	}
	actions := Actions{
		Version: "test-1.0",
		FetchToolsList: func(mcpURL string) ([]config.ToolConfig, error) {
			return []config.ToolConfig{
				{Name: "search_ast", Description: "Search AST", Enabled: true},
				{Name: "execute_code", Description: "Run code", Enabled: true},
			}, nil
		},
		MergeSyncedTools: func(server []config.ToolConfig, existing map[string]config.ToolConfig) []config.ToolConfig {
			var out []config.ToolConfig
			for _, st := range server {
				st.Description = ""
				if prev, ok := existing[st.Name]; ok {
					st.Enabled, st.Tier, st.Description = prev.Enabled, prev.Tier, prev.Description
				}
				out = append(out, st)
			}
			return out
		},
		FindCommand: func(cmd string) error {
			if cmd == "/missing/bin" {
				return os.ErrNotExist
			}
			return nil
		},
	}
	return New(actions, Options{Assets: assets}), home
}

func doReq(t *testing.T, h http.Handler, method, target string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		r = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, target, r)
	req.Host = "127.0.0.1:7831"
	if body != nil || method == http.MethodPost || method == http.MethodPut {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestSecurityGuard(t *testing.T) {
	srv, _ := setupTestServer(t)
	h := srv.Handler()

	// Reject foreign Host (DNS rebinding protection)
	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	req.Host = "evil.example.com:7831"
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for foreign host, got %d", rr.Code)
	}

	// Reject foreign Origin
	rr = doReq(t, h, http.MethodGet, "/api/info", nil, map[string]string{"Origin": "https://evil.example.com"})
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for foreign origin, got %d", rr.Code)
	}

	// Reject non-JSON POST
	req = httptest.NewRequest(http.MethodPost, "/api/services", strings.NewReader("name=x"))
	req.Host = "localhost:7831"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415 for non-json post, got %d", rr.Code)
	}
}

func TestServiceCRUDToolsAndRawConfig(t *testing.T) {
	srv, home := setupTestServer(t)
	h := srv.Handler()

	// Initial empty state
	rr := doReq(t, h, http.MethodGet, "/api/state", nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /api/state: %d %s", rr.Code, rr.Body.String())
	}
	var st StateResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &st); err != nil {
		t.Fatal(err)
	}
	if st.ConfigExists || len(st.Services) != 0 {
		t.Fatalf("expected empty nonexistent config, got %+v", st)
	}

	// Create a service with ${HOME}
	rr = doReq(t, h, http.MethodPost, "/api/services", ServiceUpsertRequest{
		Name:    "ast-context-cache",
		Type:    "http",
		Command: "${HOME}/bin/ast-mcp",
		Port:    7821,
		MCPURL:  "http://localhost:7821/mcp",
		Env:     map[string]string{"AST_HOME": "${HOME}/.astcache"},
	}, nil)
	if rr.Code != http.StatusCreated {
		t.Fatalf("POST /api/services: %d %s", rr.Code, rr.Body.String())
	}

	// Update agents config
	rr = doReq(t, h, http.MethodPut, "/api/agents", AgentsDTO{
		Managed:  true,
		OpenCode: true,
		Cursor:   true,
		Claude:   false,
	}, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("PUT /api/agents: %d %s", rr.Code, rr.Body.String())
	}

	// Register service -> should write OpenCode and Cursor files
	rr = doReq(t, h, http.MethodPost, "/api/services/ast-context-cache/register", map[string]any{}, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("register: %d %s", rr.Code, rr.Body.String())
	}

	// Sync tools using mock FetchToolsList
	rr = doReq(t, h, http.MethodPost, "/api/services/ast-context-cache/tools/sync", map[string]any{}, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("sync tools: %d %s", rr.Code, rr.Body.String())
	}
	var toolsResp ToolsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &toolsResp); err != nil {
		t.Fatal(err)
	}
	if len(toolsResp.Tools) != 2 {
		t.Fatalf("expected 2 synced tools, got %d", len(toolsResp.Tools))
	}

	// Disable execute_code and set tier on search_ast, apply to tools.json
	tier := "core"
	toolsResp.Tools[0].Tier = "core"
	toolsResp.Tools[1].Enabled = false
	rr = doReq(t, h, http.MethodPut, "/api/services/ast-context-cache/tools", ToolsUpdateRequest{
		ActiveTier: &tier,
		Tools:      toolsResp.Tools,
	}, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("update tools: %d %s", rr.Code, rr.Body.String())
	}
	toolsJSONPath := filepath.Join(home, ".astcache", "tools.json")
	rawTools, err := os.ReadFile(toolsJSONPath)
	if err != nil {
		t.Fatalf("expected tools.json at %s: %v", toolsJSONPath, err)
	}
	if !strings.Contains(string(rawTools), `"execute_code"`) || !strings.Contains(string(rawTools), `"enabled": false`) {
		t.Fatalf("unexpected tools.json contents:\n%s", rawTools)
	}

	// Update service port and confirm ${HOME} and tools are preserved
	rr = doReq(t, h, http.MethodPut, "/api/services/ast-context-cache", ServiceUpsertRequest{
		Name:    "ast-context-cache",
		Type:    "http",
		Command: "${HOME}/bin/ast-mcp",
		Port:    7829,
		MCPURL:  "http://localhost:7829/mcp",
		Env:     map[string]string{"AST_HOME": "${HOME}/.astcache"},
	}, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("PUT /api/services/ast-context-cache: %d %s", rr.Code, rr.Body.String())
	}
	var updated ServiceDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &updated); err != nil {
		t.Fatal(err)
	}
	if updated.Port != 7829 || len(updated.Tools) != 2 || updated.Command != "${HOME}/bin/ast-mcp" {
		t.Fatalf("unexpected updated service: %+v", updated)
	}

	// Invalid raw YAML should return 400
	rr = doReq(t, h, http.MethodPut, "/api/config/raw", map[string]string{"content": "services:\n  - [broken"}, nil)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid YAML, got %d", rr.Code)
	}

	// SPA fallback serves index.html
	req := httptest.NewRequest(http.MethodGet, "/some/spa/route", nil)
	req.Host = "127.0.0.1:7831"
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "mcp-local-95") {
		t.Fatalf("expected SPA index.html, got %d %s", rr.Code, rr.Body.String())
	}

	// Delete service
	rr = doReq(t, h, http.MethodDelete, "/api/services/ast-context-cache?deregister=true", nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("DELETE /api/services/ast-context-cache: %d %s", rr.Code, rr.Body.String())
	}
}

func TestRegistryServerToPreset(t *testing.T) {
	raw := []byte(`{
		"name": "io.github.upstash/context7",
		"title": "Context7",
		"description": "Up-to-date docs for AI agents",
		"packages": [{
			"registryType": "npm",
			"identifier": "@upstash/context7-mcp",
			"transport": {"type": "stdio"},
			"environmentVariables": [{"name": "CONTEXT7_API_KEY", "description": "API key"}]
		}]
	}`)
	var srv mcpRegistryServer
	if err := json.Unmarshal(raw, &srv); err != nil {
		t.Fatal(err)
	}
	preset, ok := registryServerToPreset(srv)
	if !ok {
		t.Fatal("expected preset to be converted")
	}
	if preset.Service.Name != "context7" || preset.Service.Command != "npx" || len(preset.Service.Args) != 2 {
		t.Fatalf("unexpected preset service: %+v", preset.Service)
	}
	if _, hasKey := preset.Service.Env["CONTEXT7_API_KEY"]; !hasKey {
		t.Fatalf("expected CONTEXT7_API_KEY in env: %+v", preset.Service.Env)
	}
}
