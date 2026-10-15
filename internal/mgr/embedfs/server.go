package embedfs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var protocolVersion = "2024-11-05"

type Server struct {
	server  *http.Server
	root    string // absolute, symlinks resolved
	rootErr error
}

func New(root string) *Server {
	resolved, err := resolveRoot(root)
	return &Server{root: resolved, rootErr: err}
}

func resolveRoot(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("root %s: %w", root, err)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", fmt.Errorf("root %s: %w", root, err)
	}
	return resolved, nil
}

// Handler is the full HTTP surface: /mcp and /health behind the localhost guard.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", s.handleMCP)
	mux.HandleFunc("/health", s.handleHealth)
	return guard(mux)
}

func (s *Server) Serve(port int, logPath string) error {
	if s.rootErr != nil {
		return s.rootErr
	}
	s.server = &http.Server{Handler: s.Handler(), ReadHeaderTimeout: 10 * time.Second}
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("listen :%d: %w", port, err)
	}

	return s.server.Serve(ln)
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

// guard blocks browsers from driving the server: a cross-site page can send a
// CORS-free text/plain POST, and DNS rebinding makes a foreign name resolve to 127.0.0.1.
func guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isLocalHost(r.Host) {
			http.Error(w, "forbidden host", http.StatusForbidden)
			return
		}
		if o := r.Header.Get("Origin"); o != "" && !isLocalOrigin(o) {
			http.Error(w, "forbidden origin", http.StatusForbidden)
			return
		}
		if r.Method == http.MethodPost {
			if mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err != nil || mt != "application/json" {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// isLocalHost accepts localhost, 127.0.0.1 or [::1], with or without a port.
func isLocalHost(hostport string) bool {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		host = strings.TrimSuffix(strings.TrimPrefix(hostport, "["), "]")
	}
	return strings.EqualFold(host, "localhost") || host == "127.0.0.1" || host == "::1"
}

func isLocalOrigin(origin string) bool {
	u, err := url.Parse(origin)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && isLocalHost(u.Host)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type jsonRPCReq struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonRPCResp struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *Server) handleMCP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// No GET SSE stream or session DELETE: 405 per the Streamable HTTP transport.
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, nil, -32700, "Parse error")
		return
	}

	var req jsonRPCReq
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, nil, -32700, "Parse error")
		return
	}

	// Notifications (no id) and client responses (no method) get 202 with no body.
	if len(req.ID) == 0 || req.Method == "" {
		w.Header().Del("Content-Type")
		w.WriteHeader(http.StatusAccepted)
		return
	}

	var resp jsonRPCResp
	if req.Method == "initialize" {
		resp = s.handleInitialize(req)
	} else {
		resp = s.handleToolCall(req)
	}
	if resp.ID == nil {
		resp.ID = parseID(req.ID)
	}
	json.NewEncoder(w).Encode(resp)
}

func parseID(raw json.RawMessage) interface{} {
	var id interface{}
	json.Unmarshal(raw, &id)
	return id
}

func (s *Server) handleInitialize(req jsonRPCReq) jsonRPCResp {
	return jsonRPCResp{
		JSONRPC: "2.0",
		ID:      parseID(req.ID),
		Result: map[string]interface{}{
			"protocolVersion": protocolVersion,
			"serverInfo": map[string]string{
				"name":    "mcp-local embedfs",
				"version": "0.1.0",
			},
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
		},
	}
}

func (s *Server) handleToolCall(req jsonRPCReq) jsonRPCResp {
	if req.Method == "tools/list" {
		return s.handleToolsList(req)
	}
	if req.Method == "tools/call" {
		return s.handleToolsCall(req)
	}
	return jsonRPCResp{
		JSONRPC: "2.0",
		ID:      parseID(req.ID),
		Error:   &rpcError{Code: -32601, Message: fmt.Sprintf("Method not found: %s", req.Method)},
	}
}

func (s *Server) handleToolsList(req jsonRPCReq) jsonRPCResp {
	return jsonRPCResp{
		JSONRPC: "2.0",
		ID:      parseID(req.ID),
		Result: map[string]interface{}{
			"tools": toolDefinitions(),
		},
	}
}

func (s *Server) handleToolsCall(req jsonRPCReq) jsonRPCResp {
	id := parseID(req.ID)

	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return jsonRPCResp{
			JSONRPC: "2.0",
			ID:      id,
			Error:   &rpcError{Code: -32602, Message: "Invalid params"},
		}
	}

	handler := toolHandler(params.Name)
	if handler == nil {
		return jsonRPCResp{
			JSONRPC: "2.0",
			ID:      id,
			Error:   &rpcError{Code: -32601, Message: fmt.Sprintf("Unknown tool: %s", params.Name)},
		}
	}

	if s.rootErr != nil {
		return jsonRPCResp{JSONRPC: "2.0", ID: id, Error: &rpcError{Code: -32000, Message: s.rootErr.Error()}}
	}
	result, err := handler(s.root, params.Arguments)
	if err != nil {
		return jsonRPCResp{
			JSONRPC: "2.0",
			ID:      id,
			Error:   &rpcError{Code: -32000, Message: err.Error()},
		}
	}

	return jsonRPCResp{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": result},
			},
		},
	}
}

func writeError(w http.ResponseWriter, id interface{}, code int, msg string) {
	resp := jsonRPCResp{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &rpcError{Code: code, Message: msg},
	}
	json.NewEncoder(w).Encode(resp)
}

// safeJoin maps sub onto root (already symlink-resolved) and resolves symlinks
// in the result, so a link inside root cannot carry a read or write outside it.
func safeJoin(root, sub string) (string, error) {
	joined := filepath.Join(root, sub)
	if !isWithin(joined, root) {
		return "", fmt.Errorf("path escapes root: %s", sub)
	}
	target, err := resolvePath(joined)
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", sub, err)
	}
	if !isWithin(target, root) {
		return "", fmt.Errorf("path escapes root: %s", sub)
	}
	return target, nil
}

// resolvePath resolves symlinks in the deepest existing prefix of the clean
// absolute path p and re-joins the not-yet-existing rest.
func resolvePath(p string) (string, error) {
	cur, rest := p, ""
	for {
		_, err := os.Lstat(cur)
		if err == nil {
			// Fails for a dangling symlink, which a write would otherwise follow.
			resolved, err := filepath.EvalSymlinks(cur)
			if err != nil {
				return "", err
			}
			return filepath.Join(resolved, rest), nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return p, nil
		}
		rest = filepath.Join(filepath.Base(cur), rest)
		cur = parent
	}
}

func isWithin(target, root string) bool {
	rel, err := filepath.Rel(root, target)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// LogTo appends a timestamped line to logPath. The log may hold paths and errors from
// the served tree, so it is private and never written through a planted symlink.
func LogTo(logPath, line string) {
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND|syscall.O_NOFOLLOW, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), line))
}
