package embedfs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleMCP(t *testing.T) {
	s := New(t.TempDir())
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantBody   string
	}{
		{"notification is 202", http.MethodPost, `{"jsonrpc":"2.0","method":"notifications/initialized"}`, http.StatusAccepted, ""},
		{"other notification is 202", http.MethodPost, `{"jsonrpc":"2.0","method":"notifications/cancelled","params":{}}`, http.StatusAccepted, ""},
		{"client response is 202", http.MethodPost, `{"jsonrpc":"2.0","id":9,"result":{}}`, http.StatusAccepted, ""},
		{"initialize", http.MethodPost, `{"jsonrpc":"2.0","id":1,"method":"initialize"}`, http.StatusOK, `"protocolVersion"`},
		{"tools/list", http.MethodPost, `{"jsonrpc":"2.0","id":"a","method":"tools/list"}`, http.StatusOK, `"id":"a"`},
		{"GET not allowed", http.MethodGet, "", http.StatusMethodNotAllowed, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			s.handleMCP(rec, httptest.NewRequest(tt.method, "/mcp", strings.NewReader(tt.body)))
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantBody == "" && rec.Body.Len() != 0 {
				t.Errorf("body = %q, want empty", rec.Body.String())
			}
			if tt.wantBody != "" && !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}
