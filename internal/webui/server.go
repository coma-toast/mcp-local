package webui

import (
	"bytes"
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
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/agents"
	"github.com/coma-toast/mcp-local/internal/mgr/asttools"
	"github.com/coma-toast/mcp-local/internal/mgr/claudedesktop"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/cursor"
	"github.com/coma-toast/mcp-local/internal/mgr/jsonagent"
	"github.com/coma-toast/mcp-local/internal/mgr/logs"
	"github.com/coma-toast/mcp-local/internal/mgr/opencode"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/mgr/runtime"
	"github.com/coma-toast/mcp-local/internal/mgr/safefile"
)

const maxBodyBytes = 4 << 20 // 4 MiB

type Options struct {
	AllowRemote bool
	Assets      fs.FS
}

type Server struct {
	actions     Actions
	allowRemote bool
	assets      fs.FS
	mu          sync.Mutex
	httpServer  *http.Server
}

func New(actions Actions, opts Options) *Server {
	assets := opts.Assets
	if assets == nil {
		assets = defaultAssets()
	}
	return &Server{
		actions:     actions,
		allowRemote: opts.AllowRemote,
		assets:      assets,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/info", s.handleInfo)
	mux.HandleFunc("GET /api/state", s.handleState)
	mux.HandleFunc("PUT /api/agents", s.handleUpdateAgents)
	mux.HandleFunc("GET /api/agents/configs", s.handleAgentConfigs)

	mux.HandleFunc("POST /api/services", s.handleCreateService)
	mux.HandleFunc("PUT /api/services/{name}", s.handleUpdateService)
	mux.HandleFunc("DELETE /api/services/{name}", s.handleDeleteService)

	mux.HandleFunc("POST /api/services/{name}/start", s.handleStartService)
	mux.HandleFunc("POST /api/services/{name}/stop", s.handleStopService)
	mux.HandleFunc("POST /api/services/{name}/restart", s.handleRestartService)
	mux.HandleFunc("POST /api/services/{name}/rebuild", s.handleRebuildService)
	mux.HandleFunc("POST /api/services/{name}/register", s.handleRegisterService)
	mux.HandleFunc("POST /api/services/{name}/deregister", s.handleDeregisterService)
	mux.HandleFunc("GET /api/services/{name}/health", s.handleServiceHealth)

	mux.HandleFunc("GET /api/services/{name}/tools", s.handleGetTools)
	mux.HandleFunc("POST /api/services/{name}/tools/sync", s.handleSyncTools)
	mux.HandleFunc("PUT /api/services/{name}/tools", s.handleUpdateTools)
	mux.HandleFunc("POST /api/services/{name}/tools/apply", s.handleApplyTools)

	mux.HandleFunc("GET /api/services/{name}/logs", s.handleGetLogs)
	mux.HandleFunc("GET /api/services/{name}/logs/stream", s.handleStreamLogs)

	mux.HandleFunc("POST /api/actions/start-all", s.handleStartAll)
	mux.HandleFunc("POST /api/actions/stop-all", s.handleStopAll)
	mux.HandleFunc("POST /api/actions/restart-all", s.handleRestartAll)
	mux.HandleFunc("POST /api/actions/register-all", s.handleRegisterAll)

	mux.HandleFunc("GET /api/validate", s.handleValidate)
	mux.HandleFunc("POST /api/doctor", s.handleDoctor)

	mux.HandleFunc("GET /api/import/preview", s.handleImportPreview)
	mux.HandleFunc("POST /api/import", s.handleImportApply)

	mux.HandleFunc("GET /api/config/raw", s.handleGetRawConfig)
	mux.HandleFunc("PUT /api/config/raw", s.handlePutRawConfig)

	mux.HandleFunc("GET /api/registry/search", s.handleRegistrySearch)

	mux.Handle("/", s.spaHandler())
	return s.guard(mux)
}

func (s *Server) ListenAndServe(ctx context.Context, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.httpServer = &http.Server{
		Handler:           s.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		if err := s.httpServer.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()
	select {
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.httpServer.Shutdown(shutCtx)
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.allowRemote {
			if !isLocalHost(r.Host) {
				http.Error(w, "forbidden host", http.StatusForbidden)
				return
			}
			if o := r.Header.Get("Origin"); o != "" && !isLocalOrigin(o) {
				http.Error(w, "forbidden origin", http.StatusForbidden)
				return
			}
		}
		if strings.HasPrefix(r.URL.Path, "/api/") && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) {
			ct := r.Header.Get("Content-Type")
			if ct == "" {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
			if mt, _, err := mime.ParseMediaType(ct); err != nil || mt != "application/json" {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

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

func (s *Server) spaHandler() http.Handler {
	fileServer := http.FileServer(http.FS(s.assets))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			writeError(w, http.StatusNotFound, fmt.Errorf("unknown endpoint %s", r.URL.Path))
			return
		}
		clean := strings.TrimPrefix(filepath.Clean(r.URL.Path), "/")
		if clean == "" || clean == "." {
			clean = "index.html"
		}
		if f, err := s.assets.Open(clean); err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}
		data, err := fs.ReadFile(s.assets, "index.html")
		if err != nil {
			http.Error(w, "frontend assets not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func decodeJSON(r *http.Request, dst interface{}) error {
	body := io.LimitReader(r.Body, maxBodyBytes)
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	return json.Unmarshal(data, dst)
}

func loadConfigOrEmpty() (*config.ManagerConfig, bool, error) {
	cfg, err := config.LoadConfig()
	if err == nil {
		return cfg, true, nil
	}
	if os.IsNotExist(err) {
		return &config.ManagerConfig{}, false, nil
	}
	return nil, true, err
}

func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	cfgPath, err := config.GetConfigPath()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	var exists, isSym bool
	var target string
	if fi, err := os.Lstat(cfgPath); err == nil {
		exists = true
		if fi.Mode()&os.ModeSymlink != 0 {
			isSym = true
			target, _ = os.Readlink(cfgPath)
		}
	}
	ocPath := opencode.ConfigPath()
	curPath := cursor.ConfigPath()
	claudePath := claudedesktop.ConfigPath()
	writeJSON(w, http.StatusOK, InfoResponse{
		Version:             s.actions.Version,
		ConfigPath:          cfgPath,
		ConfigExists:        exists,
		ConfigIsSymlink:     isSym,
		ConfigSymlinkTarget: target,
		AgentPaths: map[string]string{
			"opencode": ocPath,
			"cursor":   curPath,
			"claude":   claudePath,
		},
		AgentFilesExist: map[string]bool{
			"opencode": fileExists(ocPath),
			"cursor":   fileExists(curPath),
			"claude":   fileExists(claudePath),
		},
	})
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	cfg, exists, err := loadConfigOrEmpty()
	if err != nil {
		writeJSON(w, http.StatusOK, StateResponse{
			ConfigExists: exists,
			ConfigError:  err.Error(),
			Agents:       AgentsDTO{Managed: false, OpenCode: true, Cursor: true, Claude: true},
			Services:     []ServiceDTO{},
		})
		return
	}
	targets := agents.TargetsFromConfig(*cfg)
	agentsDTO := AgentsDTO{
		Managed:  cfg.Agents != nil,
		OpenCode: targets.OpenCode,
		Cursor:   targets.Cursor,
		Claude:   targets.Claude,
	}
	services := make([]ServiceDTO, 0, len(cfg.Services))
	for _, name := range cfg.SortedNames() {
		svc, err := cfg.ServiceNamed(name)
		if err != nil {
			continue
		}
		services = append(services, toServiceDTO(cfg, svc))
	}
	writeJSON(w, http.StatusOK, StateResponse{
		ConfigExists: exists,
		Agents:       agentsDTO,
		Services:     services,
	})
}

func toServiceDTO(cfg *config.ManagerConfig, svc config.ServiceConfig) ServiceDTO {
	st := runtime.CheckStatus(svc)
	raw, hasRaw := cfg.RawService(svc.Name)
	pick := func(expanded, unexpanded string) string {
		if hasRaw {
			return unexpanded
		}
		return expanded
	}
	deps := svc.Deps
	env := svc.Env
	if hasRaw {
		deps = raw.Deps
		env = raw.Env
	}
	tools := make([]ToolDTO, 0, len(svc.Tools))
	for _, t := range svc.Tools {
		tools = append(tools, ToolDTO{
			Name:          t.Name,
			Enabled:       t.Enabled,
			Tier:          t.Tier,
			EffectiveTier: asttools.EffectiveTierForDisplay(t.Name, t.Tier),
			Description:   t.Description,
			InputSchema:   t.InputSchema,
		})
	}
	return ServiceDTO{
		Name:               svc.Name,
		Type:               svc.MCPType,
		Command:            pick(svc.Command, raw.Command),
		Args:               append([]string(nil), svc.Args...),
		Env:                cloneMap(env),
		Port:               svc.Port,
		ExtraPorts:         append([]int(nil), svc.ExtraPorts...),
		BuildCmd:           svc.BuildCmd,
		BuildCommand:       svc.BuildCommand,
		NoBuildOnStart:     svc.NoBuildOnStart,
		Deps:               append([]string(nil), deps...),
		Description:        svc.Description,
		HealthURL:          pick(svc.HealthURL, raw.HealthURL),
		DashboardURL:       pick(svc.DashboardURL, raw.DashboardURL),
		MCPURL:             pick(svc.MCPURL, raw.MCPURL),
		Headers:            cloneMap(svc.Headers),
		Log:                pick(svc.Log, raw.Log),
		ActiveTier:         svc.ActiveTier,
		CodeMode:           svc.CodeMode,
		NoCodeMode:         svc.NoCodeMode,
		ToolsConfigPath:    pick(svc.ToolsConfigPath, raw.ToolsConfigPath),
		Path:               pick(svc.Path, raw.Path),
		Installer:          svc.Installer,
		Tools:              tools,
		Mode:               process.ModeOf(svc).String(),
		State:              string(st.State),
		Detail:             st.Detail,
		Running:            st.Running,
		PID:                st.PID,
		Uptime:             st.Uptime,
		Stale:              st.Stale,
		OrphanPID:          st.Orphan.PID,
		OrphanPort:         st.Orphan.Port,
		MCPEndpoint:        svc.MCPEndpoint(),
		ResolvedLogPath:    process.LogPath(svc),
		ResolvedToolsPath:  asttools.ToolsConfigPath(svc),
		EffectiveInstaller: svc.EffectiveInstaller(),
		EffectiveBuild:     config.EffectiveBuild(svc),
		OverridesCount:     asttools.OverridesCount(svc.Tools),
		Registered: map[string]bool{
			"opencode": isRegisteredInFile(opencode.ConfigPath(), svc.Name),
			"cursor":   isRegisteredInFile(cursor.ConfigPath(), svc.Name),
			"claude":   isRegisteredInFile(claudedesktop.ConfigPath(), svc.Name),
		},
	}
}

func cloneMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func isRegisteredInFile(path, name string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	m, err := jsonagent.ParseJSONC(data)
	if err != nil {
		return false
	}
	for _, key := range []string{"mcp", "mcpServers"} {
		if block, ok := m[key].(map[string]interface{}); ok {
			if _, exists := block[name]; exists {
				return true
			}
		}
	}
	return false
}

func (s *Server) handleUpdateAgents(w http.ResponseWriter, r *http.Request) {
	var req AgentsDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cfg, _, err := loadConfigOrEmpty()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if !req.Managed {
		cfg.Agents = nil
	} else {
		cfg.Agents = &config.AgentsConfig{
			OpenCode: req.OpenCode,
			Cursor:   req.Cursor,
			Claude:   req.Claude,
		}
	}
	if err := config.SaveConfig(cfg); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	targets := agents.TargetsFromConfig(*cfg)
	writeJSON(w, http.StatusOK, AgentsDTO{
		Managed:  cfg.Agents != nil,
		OpenCode: targets.OpenCode,
		Cursor:   targets.Cursor,
		Claude:   targets.Claude,
	})
}

func applyUpsertFields(dst *config.ServiceConfig, req ServiceUpsertRequest) {
	dst.Name = strings.TrimSpace(req.Name)
	dst.MCPType = strings.TrimSpace(req.Type)
	dst.Command = strings.TrimSpace(req.Command)
	dst.Args = req.Args
	dst.Env = req.Env
	dst.Port = req.Port
	dst.ExtraPorts = req.ExtraPorts
	if dst.BuildCmd != "" {
		dst.BuildCmd = req.BuildCommand
		if req.BuildCommand == "" {
			dst.BuildCommand = ""
		}
	} else {
		dst.BuildCommand = req.BuildCommand
	}
	dst.NoBuildOnStart = req.NoBuildOnStart
	dst.Deps = req.Deps
	dst.Description = req.Description
	dst.HealthURL = strings.TrimSpace(req.HealthURL)
	dst.DashboardURL = strings.TrimSpace(req.DashboardURL)
	dst.MCPURL = strings.TrimSpace(req.MCPURL)
	dst.Headers = req.Headers
	dst.Log = strings.TrimSpace(req.Log)
	dst.ActiveTier = strings.TrimSpace(req.ActiveTier)
	dst.CodeMode = req.CodeMode
	dst.NoCodeMode = req.NoCodeMode
	dst.ToolsConfigPath = strings.TrimSpace(req.ToolsConfigPath)
	dst.Path = strings.TrimSpace(req.Path)
	dst.Installer = strings.TrimSpace(req.Installer)
}

func validateUpsertRequest(req ServiceUpsertRequest) error {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("service name cannot contain slashes")
	}
	t := strings.ToLower(strings.TrimSpace(req.Type))
	switch t {
	case "", "http", "stdio", "filesystem":
	default:
		return fmt.Errorf("invalid type %q (must be http, stdio, filesystem, or empty)", req.Type)
	}
	if t == "filesystem" && strings.TrimSpace(req.Path) == "" {
		return fmt.Errorf("path is required when type is filesystem")
	}
	if t != "filesystem" && strings.TrimSpace(req.Command) == "" && strings.TrimSpace(req.MCPURL) == "" {
		return fmt.Errorf("either command or mcp_url is required")
	}
	if req.Port < 0 || req.Port > 65535 {
		return fmt.Errorf("port %d out of range", req.Port)
	}
	return nil
}

func (s *Server) handleCreateService(w http.ResponseWriter, r *http.Request) {
	var req ServiceUpsertRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateUpsertRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cfg, _, err := loadConfigOrEmpty()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if _, err := cfg.ServiceNamed(strings.TrimSpace(req.Name)); err == nil {
		writeError(w, http.StatusConflict, fmt.Errorf("service %q already exists", req.Name))
		return
	}
	var svc config.ServiceConfig
	applyUpsertFields(&svc, req)
	cfg.UpsertService(svc)
	if err := config.SaveConfig(cfg); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	reloaded, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	saved, _ := reloaded.ServiceNamed(svc.Name)
	writeJSON(w, http.StatusCreated, toServiceDTO(reloaded, saved))
}

func (s *Server) handleUpdateService(w http.ResponseWriter, r *http.Request) {
	target := r.PathValue("name")
	var req ServiceUpsertRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		req.Name = target
	}
	if err := validateUpsertRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	idx := -1
	for i := range cfg.Services {
		if cfg.Services[i].Name == target {
			idx = i
			break
		}
	}
	if idx < 0 {
		writeError(w, http.StatusNotFound, fmt.Errorf("unknown service %q", target))
		return
	}
	newName := strings.TrimSpace(req.Name)
	if newName != target {
		if _, err := cfg.ServiceNamed(newName); err == nil {
			writeError(w, http.StatusConflict, fmt.Errorf("service %q already exists", newName))
			return
		}
	}
	svc := cfg.Services[idx]
	applyUpsertFields(&svc, req)
	cfg.Services[idx] = svc
	if err := config.SaveConfig(cfg); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	reloaded, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	saved, _ := reloaded.ServiceNamed(newName)
	writeJSON(w, http.StatusOK, toServiceDTO(reloaded, saved))
}

func (s *Server) handleDeleteService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	deregister := r.URL.Query().Get("deregister") != "false"
	s.mu.Lock()
	defer s.mu.Unlock()
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	var out bytes.Buffer
	if deregister {
		if err := agents.DeregisterAll(svc, agents.TargetsFromConfig(*cfg)); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("deregister failed: %w", err))
			return
		}
		fmt.Fprintf(&out, "Deregistered %s from enabled agents\n", name)
	}
	remaining := make([]config.ServiceConfig, 0, len(cfg.Services))
	for _, item := range cfg.Services {
		if item.Name != name {
			remaining = append(remaining, item)
		}
	}
	cfg.Services = remaining
	if err := config.SaveConfig(cfg); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	fmt.Fprintf(&out, "Removed %s from config.yaml\n", name)
	writeJSON(w, http.StatusOK, ActionResponse{OK: true, Output: out.String()})
}

func (s *Server) handleStartService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var req struct {
		Build *bool `json:"build"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	build := true
	if req.Build != nil {
		build = *req.Build
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	var out, errOut bytes.Buffer
	started, failed := s.actions.StartServices(cfg, []string{name}, build, &out, &errOut)
	if len(started) > 0 && s.actions.RegisterStarted != nil {
		s.actions.RegisterStarted(cfg, started, &out)
	}
	combined := strings.TrimSpace(out.String() + "\n" + errOut.String())
	dto := toServiceDTO(cfg, svc)
	if failed > 0 {
		writeJSON(w, http.StatusOK, ActionResponse{OK: false, Output: combined, Error: fmt.Sprintf("failed to start %s", name), Status: &dto})
		return
	}
	writeJSON(w, http.StatusOK, ActionResponse{OK: true, Output: combined, Status: &dto})
}

func (s *Server) handleStopService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var req struct {
		Deregister *bool `json:"deregister"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	deregister := true
	if req.Deregister != nil {
		deregister = *req.Deregister
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	var out, errOut bytes.Buffer
	stopErr := s.actions.StopService(cfg, name, deregister, &out, &errOut)
	combined := strings.TrimSpace(out.String() + "\n" + errOut.String())
	dto := toServiceDTO(cfg, svc)
	if stopErr != nil {
		writeJSON(w, http.StatusOK, ActionResponse{OK: false, Output: combined, Error: stopErr.Error(), Status: &dto})
		return
	}
	writeJSON(w, http.StatusOK, ActionResponse{OK: true, Output: combined, Status: &dto})
}

func (s *Server) handleRestartService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	var out, errOut bytes.Buffer
	started, failed := s.actions.RestartServices(cfg, []string{name}, &out, &errOut)
	if len(started) > 0 && s.actions.RegisterStarted != nil {
		s.actions.RegisterStarted(cfg, started, &out)
	}
	combined := strings.TrimSpace(out.String() + "\n" + errOut.String())
	dto := toServiceDTO(cfg, svc)
	if failed > 0 {
		writeJSON(w, http.StatusOK, ActionResponse{OK: false, Output: combined, Error: fmt.Sprintf("failed to restart %s", name), Status: &dto})
		return
	}
	writeJSON(w, http.StatusOK, ActionResponse{OK: true, Output: combined, Status: &dto})
}

func (s *Server) handleRebuildService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if config.EffectiveBuild(svc) == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("no build_command configured for %s", name))
		return
	}
	if err := process.Build(svc); err != nil {
		writeJSON(w, http.StatusOK, ActionResponse{OK: false, Error: err.Error(), Output: err.Error()})
		return
	}
	var out bytes.Buffer
	fmt.Fprintf(&out, "✅ %s: rebuild complete\n", name)
	if err := agents.RegisterAll([]config.ServiceConfig{svc}, agents.TargetsFromConfig(*cfg)); err != nil {
		fmt.Fprintf(&out, "⚠️ Registration failed: %v\n", err)
	} else {
		fmt.Fprintf(&out, "✅ Agent registration refreshed\n")
	}
	dto := toServiceDTO(cfg, svc)
	writeJSON(w, http.StatusOK, ActionResponse{OK: true, Output: out.String(), Status: &dto})
}

func (s *Server) handleRegisterService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	targets := agents.TargetsFromConfig(*cfg)
	if err := agents.RegisterAll([]config.ServiceConfig{svc}, targets); err != nil {
		dto := toServiceDTO(cfg, svc)
		writeJSON(w, http.StatusOK, ActionResponse{OK: false, Error: err.Error(), Output: err.Error(), Status: &dto})
		return
	}
	dto := toServiceDTO(cfg, svc)
	writeJSON(w, http.StatusOK, ActionResponse{
		OK:     true,
		Output: fmt.Sprintf("Registered %s with %v", name, targets.Hosts()),
		Status: &dto,
	})
}

func (s *Server) handleDeregisterService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	targets := agents.TargetsFromConfig(*cfg)
	if err := agents.DeregisterAll(svc, targets); err != nil {
		dto := toServiceDTO(cfg, svc)
		writeJSON(w, http.StatusOK, ActionResponse{OK: false, Error: err.Error(), Output: err.Error(), Status: &dto})
		return
	}
	dto := toServiceDTO(cfg, svc)
	writeJSON(w, http.StatusOK, ActionResponse{
		OK:     true,
		Output: fmt.Sprintf("Deregistered %s from %v", name, targets.Hosts()),
		Status: &dto,
	})
}

func (s *Server) handleStartAll(w http.ResponseWriter, r *http.Request) {
	cfg, _, err := loadConfigOrEmpty()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	names := cfg.SortedNames()
	var out, errOut bytes.Buffer
	started, failed := s.actions.StartServices(cfg, names, true, &out, &errOut)
	if len(started) > 0 && s.actions.RegisterStarted != nil {
		s.actions.RegisterStarted(cfg, started, &out)
	}
	combined := strings.TrimSpace(out.String() + "\n" + errOut.String())
	if failed > 0 {
		writeJSON(w, http.StatusOK, ActionResponse{OK: false, Output: combined, Error: fmt.Sprintf("%d service(s) failed to start", failed)})
		return
	}
	writeJSON(w, http.StatusOK, ActionResponse{OK: true, Output: combined})
}

func (s *Server) handleStopAll(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Deregister *bool `json:"deregister"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	deregister := true
	if req.Deregister != nil {
		deregister = *req.Deregister
	}
	cfg, _, err := loadConfigOrEmpty()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	var out, errOut bytes.Buffer
	failed := 0
	for _, name := range cfg.SortedNames() {
		if err := s.actions.StopService(cfg, name, deregister, &out, &errOut); err != nil {
			fmt.Fprintf(&errOut, "  ❌ %s: %v\n", name, err)
			failed++
		}
	}
	combined := strings.TrimSpace(out.String() + "\n" + errOut.String())
	if failed > 0 {
		writeJSON(w, http.StatusOK, ActionResponse{OK: false, Output: combined, Error: fmt.Sprintf("%d service(s) failed to stop", failed)})
		return
	}
	writeJSON(w, http.StatusOK, ActionResponse{OK: true, Output: combined})
}

func (s *Server) handleRestartAll(w http.ResponseWriter, r *http.Request) {
	cfg, _, err := loadConfigOrEmpty()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	var out, errOut bytes.Buffer
	started, failed := s.actions.RestartServices(cfg, cfg.SortedNames(), &out, &errOut)
	if len(started) > 0 && s.actions.RegisterStarted != nil {
		s.actions.RegisterStarted(cfg, started, &out)
	}
	combined := strings.TrimSpace(out.String() + "\n" + errOut.String())
	if failed > 0 {
		writeJSON(w, http.StatusOK, ActionResponse{OK: false, Output: combined, Error: fmt.Sprintf("%d service(s) failed to restart", failed)})
		return
	}
	writeJSON(w, http.StatusOK, ActionResponse{OK: true, Output: combined})
}

func (s *Server) handleRegisterAll(w http.ResponseWriter, r *http.Request) {
	cfg, _, err := loadConfigOrEmpty()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	var toReg []config.ServiceConfig
	for _, name := range cfg.SortedNames() {
		svc, _ := cfg.ServiceNamed(name)
		if svc.MCPURL == "" && svc.Command == "" {
			continue
		}
		toReg = append(toReg, svc)
	}
	targets := agents.TargetsFromConfig(*cfg)
	if err := agents.RegisterAll(toReg, targets); err != nil {
		writeJSON(w, http.StatusOK, ActionResponse{OK: false, Error: err.Error(), Output: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, ActionResponse{
		OK:     true,
		Output: fmt.Sprintf("Registered %d service(s) with %v", len(toReg), targets.Hosts()),
	})
}

func (s *Server) handleServiceHealth(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if strings.TrimSpace(svc.HealthURL) == "" {
		writeJSON(w, http.StatusOK, HealthResponse{OK: false, Error: "no health_url configured"})
		return
	}
	start := time.Now()
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(svc.HealthURL)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		writeJSON(w, http.StatusOK, HealthResponse{OK: false, URL: svc.HealthURL, LatencyMs: latency, Error: err.Error()})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 256<<10))
	pretty := strings.TrimSpace(string(body))
	var obj interface{}
	if json.Unmarshal(body, &obj) == nil {
		if b, err := json.MarshalIndent(obj, "", "  "); err == nil {
			pretty = string(b)
		}
	}
	writeJSON(w, http.StatusOK, HealthResponse{
		OK:         resp.StatusCode >= 200 && resp.StatusCode < 300,
		URL:        svc.HealthURL,
		StatusCode: resp.StatusCode,
		Body:       pretty,
		LatencyMs:  latency,
	})
}

func buildToolsResponse(svc config.ServiceConfig) ToolsResponse {
	tools := make([]ToolDTO, 0, len(svc.Tools))
	for _, t := range svc.Tools {
		tools = append(tools, ToolDTO{
			Name:          t.Name,
			Enabled:       t.Enabled,
			Tier:          t.Tier,
			EffectiveTier: asttools.EffectiveTierForDisplay(t.Name, t.Tier),
			Description:   t.Description,
			InputSchema:   t.InputSchema,
		})
	}
	b, _ := asttools.MarshalOverrides(svc.Tools)
	return ToolsResponse{
		Name:            svc.Name,
		ActiveTier:      svc.ActiveTier,
		CodeMode:        svc.CodeMode,
		NoCodeMode:      svc.NoCodeMode,
		ToolsConfigPath: asttools.ToolsConfigPath(svc),
		Tools:           tools,
		OverridesJSON:   string(b),
		OverridesCount:  asttools.OverridesCount(svc.Tools),
	}
}

func (s *Server) handleGetTools(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, buildToolsResponse(svc))
}

func (s *Server) handleSyncTools(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	mcpURL := strings.TrimSpace(svc.MCPEndpoint())
	if mcpURL == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("no mcp_url or port configured for %s", name))
		return
	}
	liveTools, err := s.actions.FetchToolsList(mcpURL)
	if err != nil {
		writeError(w, http.StatusBadGateway, fmt.Errorf("failed to fetch tools: %w", err))
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cfg, err = config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	svc, err = cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	existing := make(map[string]config.ToolConfig, len(svc.Tools))
	for _, t := range svc.Tools {
		existing[t.Name] = t
	}
	merged := s.actions.MergeSyncedTools(liveTools, existing)
	for i := range cfg.Services {
		if cfg.Services[i].Name == name {
			cfg.Services[i].Tools = merged
			svc = cfg.Services[i]
			break
		}
	}
	if err := config.SaveConfig(cfg); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, buildToolsResponse(svc))
}

func (s *Server) handleUpdateTools(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var req ToolsUpdateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	idx := -1
	for i := range cfg.Services {
		if cfg.Services[i].Name == name {
			idx = i
			break
		}
	}
	if idx < 0 {
		writeError(w, http.StatusNotFound, fmt.Errorf("unknown service %q", name))
		return
	}
	svc := &cfg.Services[idx]
	if req.ActiveTier != nil {
		svc.ActiveTier = strings.TrimSpace(*req.ActiveTier)
	}
	if req.CodeMode != nil {
		svc.CodeMode = *req.CodeMode
	}
	if req.NoCodeMode != nil {
		svc.NoCodeMode = *req.NoCodeMode
	}
	if req.ToolsConfigPath != nil {
		svc.ToolsConfigPath = strings.TrimSpace(*req.ToolsConfigPath)
	}
	if req.Tools != nil {
		updated := make([]config.ToolConfig, 0, len(req.Tools))
		for _, t := range req.Tools {
			if strings.TrimSpace(t.Name) == "" {
				continue
			}
			updated = append(updated, config.ToolConfig{
				Name:        t.Name,
				Enabled:     t.Enabled,
				Tier:        strings.TrimSpace(t.Tier),
				Description: t.Description,
				InputSchema: t.InputSchema,
			})
		}
		svc.Tools = updated
	}
	if err := config.SaveConfig(cfg); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	apply := true
	if req.Apply != nil {
		apply = *req.Apply
	}
	if apply && len(svc.Tools) > 0 {
		if err := asttools.WriteToolsJSON(asttools.ToolsConfigPath(*svc), svc.Tools); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}
	writeJSON(w, http.StatusOK, buildToolsResponse(*svc))
}

func (s *Server) handleApplyTools(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	path := asttools.ToolsConfigPath(svc)
	if err := asttools.WriteToolsJSON(path, svc.Tools); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, ActionResponse{
		OK:     true,
		Output: fmt.Sprintf("Wrote %s (%d tool overrides). Restart %s to apply.", path, asttools.OverridesCount(svc.Tools), name),
	})
}

func (s *Server) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	n := 200
	if q := r.URL.Query().Get("lines"); q != "" {
		if v, err := strconv.Atoi(q); err == nil && v > 0 && v <= 2000 {
			n = v
		}
	}
	logPath := process.LogPath(svc)
	lines, err := tailFileLines(logPath, n)
	exists := true
	if os.IsNotExist(err) {
		exists = false
		lines = nil
		err = nil
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"path":   logPath,
		"exists": exists,
		"lines":  lines,
	})
}

func tailFileLines(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	off := fi.Size() - 128<<10
	if off < 0 {
		off = 0
	}
	buf := make([]byte, fi.Size()-off)
	if _, err := f.ReadAt(buf, off); err != nil && err != io.EOF {
		return nil, err
	}
	raw := strings.TrimRight(string(buf), "\n")
	if raw == "" {
		return []string{}, nil
	}
	lines := strings.Split(raw, "\n")
	if off > 0 && len(lines) > 1 {
		lines = lines[1:]
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines, nil
}

func (s *Server) handleStreamLogs(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	cfg, err := config.LoadConfig()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	svc, err := cfg.ServiceNamed(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("streaming unsupported"))
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	logPath := process.LogPath(svc)
	var mu sync.Mutex
	sendEvent := func(event, payload string) {
		mu.Lock()
		defer mu.Unlock()
		if event != "" {
			fmt.Fprintf(w, "event: %s\n", event)
		}
		for _, line := range strings.Split(payload, "\n") {
			fmt.Fprintf(w, "data: %s\n", line)
		}
		fmt.Fprint(w, "\n")
		flusher.Flush()
	}

	if lines, err := tailFileLines(logPath, 80); err == nil && len(lines) > 0 {
		for _, line := range lines {
			sendEvent("log", line)
		}
	}

	ctx := r.Context()
	done := make(chan struct{})
	go func() {
		defer close(done)
		t := time.NewTicker(15 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				mu.Lock()
				fmt.Fprint(w, ": keepalive\n\n")
				flusher.Flush()
				mu.Unlock()
			}
		}
	}()

	sw := &sseLineWriter{onLine: func(line string) { sendEvent("log", line) }}
	errW := &sseLineWriter{onLine: func(line string) { sendEvent("info", line) }}
	_ = logs.Follower{Path: logPath, Out: sw, Err: errW}.Run(ctx)
	<-done
}

type sseLineWriter struct {
	buf    bytes.Buffer
	onLine func(string)
}

func (w *sseLineWriter) Write(p []byte) (int, error) {
	w.buf.Write(p)
	for {
		s := w.buf.String()
		idx := strings.IndexByte(s, '\n')
		if idx < 0 {
			break
		}
		line := strings.TrimRight(s[:idx], "\r")
		w.buf.Next(idx + 1)
		w.onLine(line)
	}
	return len(p), nil
}

func (s *Server) handleValidate(w http.ResponseWriter, r *http.Request) {
	cfg, exists, err := loadConfigOrEmpty()
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"valid":  false,
			"issues": []ValidationIssue{{Service: "config.yaml", Level: "error", Message: err.Error()}},
		})
		return
	}
	if !exists {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"valid":  true,
			"issues": []ValidationIssue{},
		})
		return
	}
	issues := make([]ValidationIssue, 0)
	errorsCount := 0
	for _, svc := range cfg.Services {
		switch {
		case svc.Command == "" && !svc.IsHTTP():
			issues = append(issues, ValidationIssue{Service: svc.Name, Level: "error", Message: "stdio service has no command"})
			errorsCount++
		case svc.Command == "" && svc.MCPEndpoint() == "":
			issues = append(issues, ValidationIssue{Service: svc.Name, Level: "error", Message: "no command or mcp_url configured"})
			errorsCount++
		case svc.Command != "" && s.actions.FindCommand != nil:
			if err := s.actions.FindCommand(config.ExpandPath(svc.Command)); err != nil {
				issues = append(issues, ValidationIssue{Service: svc.Name, Level: "error", Message: fmt.Sprintf("command not found: %v", err)})
				errorsCount++
			}
		}
		for _, dep := range svc.Deps {
			expanded := config.ExpandPath(dep)
			if _, err := os.Stat(expanded); os.IsNotExist(err) {
				issues = append(issues, ValidationIssue{Service: svc.Name, Level: "error", Message: fmt.Sprintf("dependency missing: %s", expanded)})
				errorsCount++
			}
		}
		if svc.Port > 0 && svc.Port < 1024 && svc.Port != 80 && svc.Port != 443 {
			issues = append(issues, ValidationIssue{Service: svc.Name, Level: "warning", Message: fmt.Sprintf("port %d may require root", svc.Port)})
		}
		if svc.MCPURL != "" && !strings.HasPrefix(svc.MCPURL, "http://") && !strings.HasPrefix(svc.MCPURL, "https://") {
			issues = append(issues, ValidationIssue{Service: svc.Name, Level: "warning", Message: "mcp_url should start with http:// or https://"})
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"valid":  errorsCount == 0,
		"issues": issues,
	})
}

func (s *Server) handleDoctor(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Service string `json:"service"`
		Fix     bool   `json:"fix"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	cfg, _, err := loadConfigOrEmpty()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	names := cfg.SortedNames()
	if req.Service != "" {
		if _, err := cfg.ServiceNamed(req.Service); err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		names = []string{req.Service}
	}
	var out, errOut bytes.Buffer
	results := make([]DoctorServiceResult, 0, len(names))
	remaining := 0
	for _, name := range names {
		svc, _ := cfg.ServiceNamed(name)
		res := s.actions.DoctorService(cfg, svc, req.Fix, &out, &errOut)
		results = append(results, res)
		remaining += res.Remaining
	}
	combined := strings.TrimSpace(out.String() + "\n" + errOut.String())
	writeJSON(w, http.StatusOK, DoctorResponse{
		Services:  results,
		Remaining: remaining,
		Output:    combined,
	})
}

func (s *Server) scanImportCandidates(cfg *config.ManagerConfig) []ImportCandidate {
	existing := make(map[string]bool, len(cfg.Services))
	for _, svc := range cfg.Services {
		existing[svc.Name] = true
	}
	sources := []struct {
		name     string
		path     string
		blockKey string
	}{
		{"opencode", opencode.ConfigPath(), "mcp"},
		{"cursor", cursor.ConfigPath(), "mcpServers"},
		{"claude", claudedesktop.ConfigPath(), "mcpServers"},
	}
	var candidates []ImportCandidate
	for _, src := range sources {
		data, err := os.ReadFile(src.path)
		if err != nil {
			continue
		}
		m, err := jsonagent.ParseJSONC(data)
		if err != nil {
			continue
		}
		block, _ := m[src.blockKey].(map[string]interface{})
		if block == nil {
			continue
		}
		names := make([]string, 0, len(block))
		for k := range block {
			names = append(names, k)
		}
		sort.Strings(names)
		for _, name := range names {
			entry, _ := block[name].(map[string]interface{})
			if entry == nil {
				continue
			}
			svc, ok := s.actions.ServiceFromEntry(name, entry)
			if !ok {
				continue
			}
			candidates = append(candidates, ImportCandidate{
				Name:    name,
				Source:  src.name,
				Exists:  existing[name],
				Mode:    process.ModeOf(svc).String(),
				Command: svc.Command,
				Args:    svc.Args,
				MCPURL:  svc.MCPURL,
				Service: toServiceDTO(cfg, svc),
			})
		}
	}
	return candidates
}

func (s *Server) handleImportPreview(w http.ResponseWriter, r *http.Request) {
	cfg, _, err := loadConfigOrEmpty()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	candidates := s.scanImportCandidates(cfg)
	if candidates == nil {
		candidates = []ImportCandidate{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"candidates": candidates})
}

func (s *Server) handleImportApply(w http.ResponseWriter, r *http.Request) {
	var req ImportApplyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cfg, _, err := loadConfigOrEmpty()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	candidates := s.scanImportCandidates(cfg)
	byKey := make(map[string]ImportCandidate, len(candidates))
	for _, c := range candidates {
		byKey[c.Source+":"+c.Name] = c
	}
	imported := 0
	for _, item := range req.Items {
		c, ok := byKey[item.Source+":"+item.Name]
		if !ok {
			continue
		}
		if _, err := cfg.ServiceNamed(c.Name); err == nil && !req.Overwrite {
			continue
		}
		svc := config.ServiceConfig{
			Name:    c.Service.Name,
			MCPType: c.Service.Type,
			Command: c.Service.Command,
			Args:    c.Service.Args,
			Env:     c.Service.Env,
			MCPURL:  c.Service.MCPURL,
			Headers: c.Service.Headers,
		}
		cfg.UpsertService(svc)
		imported++
	}
	if imported > 0 {
		if err := config.SaveConfig(cfg); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"imported": imported})
}

func (s *Server) handleAgentConfigs(w http.ResponseWriter, r *http.Request) {
	cfg, _, err := loadConfigOrEmpty()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	targets := agents.TargetsFromConfig(*cfg)
	specs := []struct {
		host    string
		path    string
		enabled bool
	}{
		{"opencode", opencode.ConfigPath(), targets.OpenCode},
		{"cursor", cursor.ConfigPath(), targets.Cursor},
		{"claude", claudedesktop.ConfigPath(), targets.Claude},
	}
	views := make([]AgentConfigView, 0, len(specs))
	for _, sp := range specs {
		v := AgentConfigView{Host: sp.host, Path: sp.path, Enabled: sp.enabled}
		data, err := os.ReadFile(sp.path)
		if err != nil {
			if !os.IsNotExist(err) {
				v.Error = err.Error()
			}
			views = append(views, v)
			continue
		}
		v.Exists = true
		v.Content = string(data)
		views = append(views, v)
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"configs": views})
}

func (s *Server) handleGetRawConfig(w http.ResponseWriter, r *http.Request) {
	path, err := config.GetConfigPath()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	resp := RawConfigResponse{Path: path, Backups: listBackups(path)}
	data, err := os.ReadFile(path)
	if err == nil {
		resp.Exists = true
		resp.Content = string(data)
	} else if !os.IsNotExist(err) {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handlePutRawConfig(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if _, err := config.ParseConfig([]byte(req.Content)); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid config.yaml: %w", err))
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	path, err := config.GetConfigPath()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	content := req.Content
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	if err := safefile.Replace(path, []byte(content)); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, RawConfigResponse{
		Path:    path,
		Exists:  true,
		Content: content,
		Backups: listBackups(path),
	})
}

func listBackups(cfgPath string) []BackupInfo {
	bdir := safefile.BackupDir()
	bname := safefile.BackupName(cfgPath)
	entries, err := os.ReadDir(bdir)
	if err != nil {
		return []BackupInfo{}
	}
	var out []BackupInfo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		p := filepath.Join(bdir, e.Name(), bname)
		if fi, err := os.Stat(p); err == nil {
			out = append(out, BackupInfo{
				Timestamp: e.Name(),
				Path:      p,
				Size:      fi.Size(),
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Timestamp > out[j].Timestamp })
	return out
}
