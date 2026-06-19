package embedfs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var protocolVersion = "2024-11-05"

type Server struct {
	server *http.Server
	root   string
}

func New(root string) *Server {
	return &Server{root: root}
}

func (s *Server) Serve(port int, logPath string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", s.handleMCP)
	mux.HandleFunc("/health", s.handleHealth)

	s.server = &http.Server{
		Handler: mux,
	}

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
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		writeError(w, nil, -32700, "Method not allowed")
		return
	}

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

	var resp jsonRPCResp
	switch req.Method {
	case "initialize":
		resp = s.handleInitialize(req)
	case "notifications/initialized":
		resp = jsonRPCResp{
			JSONRPC: "2.0",
		}
	default:
		if req.ID == nil {
			resp = jsonRPCResp{JSONRPC: "2.0"}
		} else {
			resp = s.handleToolCall(req)
		}
	}

	if resp.ID == nil && req.ID != nil {
		resp.ID = parseID(req.ID)
	}

	if resp.ID == nil && resp.Error == nil && resp.Result == nil {
		return
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

func safeJoin(root, sub string) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	target := filepath.Join(rootAbs, sub)
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	if !isWithin(targetAbs, rootAbs) {
		return "", fmt.Errorf("path escapes root: %s", sub)
	}
	return targetAbs, nil
}

func isWithin(target, root string) bool {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}

func LogTo(logPath, line string) {
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), line))
}
