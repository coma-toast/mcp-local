// Package mcpserver implements a newline-delimited JSON-RPC 2.0 MCP server over stdio.
package mcpserver

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

const (
	// ProtocolVersion is the MCP protocol version negotiated during initialize.
	ProtocolVersion = "2024-11-05"
	maxLineBytes    = 10 << 20 // 10 MiB
)

// ToolDef describes an MCP tool exposed by the server.
type ToolDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

// ToolHandler executes a tool call and returns text output or an error.
type ToolHandler func(ctx context.Context, args json.RawMessage) (string, error)

type registeredTool struct {
	def     ToolDef
	handler ToolHandler
}

// Server serves registered MCP tools over a newline-delimited JSON-RPC stream.
type Server struct {
	name    string
	version string
	tools   []registeredTool
	byName  map[string]ToolHandler
	writeMu sync.Mutex
}

// New creates a new stdio MCP server with the given serverInfo name and version.
func New(name, version string) *Server {
	return &Server{
		name:    name,
		version: version,
		byName:  make(map[string]ToolHandler),
	}
}

// Register adds a tool to the server.
func (s *Server) Register(def ToolDef, handler ToolHandler) {
	if _, exists := s.byName[def.Name]; !exists {
		s.tools = append(s.tools, registeredTool{def: def, handler: handler})
	}
	s.byName[def.Name] = handler
}

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type response struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Serve reads newline-delimited JSON-RPC messages from in and writes responses to out.
func (s *Server) Serve(ctx context.Context, in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 64*1024), maxLineBytes)
	for sc.Scan() {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var req request
		if err := json.Unmarshal(line, &req); err != nil {
			if err := s.writeJSON(out, response{
				JSONRPC: "2.0",
				Error:   &rpcError{Code: -32700, Message: "Parse error"},
			}); err != nil {
				return err
			}
			continue
		}
		// Notifications have no id; client responses have no method.
		if len(req.ID) == 0 || string(req.ID) == "null" || req.Method == "" {
			continue
		}
		resp := s.dispatch(ctx, req)
		if err := s.writeJSON(out, resp); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (s *Server) dispatch(ctx context.Context, req request) response {
	id := parseID(req.ID)
	switch req.Method {
	case "initialize":
		return response{
			JSONRPC: "2.0",
			ID:      id,
			Result: map[string]any{
				"protocolVersion": ProtocolVersion,
				"serverInfo": map[string]string{
					"name":    s.name,
					"version": s.version,
				},
				"capabilities": map[string]any{
					"tools": map[string]any{},
				},
			},
		}
	case "ping":
		return response{JSONRPC: "2.0", ID: id, Result: map[string]any{}}
	case "tools/list":
		defs := make([]ToolDef, 0, len(s.tools))
		for _, t := range s.tools {
			defs = append(defs, t.def)
		}
		return response{
			JSONRPC: "2.0",
			ID:      id,
			Result:  map[string]any{"tools": defs},
		}
	case "tools/call":
		return s.handleToolCall(ctx, id, req.Params)
	default:
		return response{
			JSONRPC: "2.0",
			ID:      id,
			Error:   &rpcError{Code: -32601, Message: fmt.Sprintf("Method not found: %s", req.Method)},
		}
	}
}

func (s *Server) handleToolCall(ctx context.Context, id any, rawParams json.RawMessage) response {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(rawParams, &params); err != nil {
		return response{
			JSONRPC: "2.0",
			ID:      id,
			Error:   &rpcError{Code: -32602, Message: "Invalid params"},
		}
	}
	handler, ok := s.byName[params.Name]
	if !ok {
		return response{
			JSONRPC: "2.0",
			ID:      id,
			Error:   &rpcError{Code: -32601, Message: fmt.Sprintf("Unknown tool: %s", params.Name)},
		}
	}
	text, err := handler(ctx, params.Arguments)
	if err != nil {
		return response{
			JSONRPC: "2.0",
			ID:      id,
			Result: map[string]any{
				"content": []map[string]any{
					{"type": "text", "text": err.Error()},
				},
				"isError": true,
			},
		}
	}
	return response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]any{
			"content": []map[string]any{
				{"type": "text", "text": text},
			},
		},
	}
}

func (s *Server) writeJSON(w io.Writer, v any) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(append(b, '\n'))
	return err
}

func parseID(raw json.RawMessage) any {
	var id any
	_ = json.Unmarshal(raw, &id)
	return id
}
