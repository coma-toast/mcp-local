package webui

import (
	"io"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

// Actions supplies lifecycle, MCP inspection, and diagnostic hooks from the CLI package.
type Actions struct {
	Version          string
	StartServices    func(cfg *config.ManagerConfig, names []string, build bool, out, errOut io.Writer) ([]config.ServiceConfig, int)
	StopService      func(cfg *config.ManagerConfig, name string, deregister bool, out, errOut io.Writer) error
	RestartServices  func(cfg *config.ManagerConfig, names []string, out, errOut io.Writer) ([]config.ServiceConfig, int)
	RegisterStarted  func(cfg *config.ManagerConfig, started []config.ServiceConfig, out io.Writer)
	FetchToolsList   func(mcpURL string) ([]config.ToolConfig, error)
	MergeSyncedTools func(server []config.ToolConfig, existing map[string]config.ToolConfig) []config.ToolConfig
	FindCommand      func(cmd string) error
	ServiceFromEntry func(name string, entry map[string]interface{}) (config.ServiceConfig, bool)
	DoctorService    func(cfg *config.ManagerConfig, svc config.ServiceConfig, fix bool, out, errOut io.Writer) DoctorServiceResult
}

type InfoResponse struct {
	Version             string            `json:"version"`
	ConfigPath          string            `json:"config_path"`
	ConfigExists        bool              `json:"config_exists"`
	ConfigIsSymlink     bool              `json:"config_is_symlink"`
	ConfigSymlinkTarget string            `json:"config_symlink_target,omitempty"`
	AgentPaths          map[string]string `json:"agent_paths"`
	AgentFilesExist     map[string]bool   `json:"agent_files_exist"`
}

type AgentsDTO struct {
	Managed  bool `json:"managed"`
	OpenCode bool `json:"opencode"`
	Cursor   bool `json:"cursor"`
	Claude   bool `json:"claude"`
}

type ToolDTO struct {
	Name          string                 `json:"name"`
	Enabled       bool                   `json:"enabled"`
	Tier          string                 `json:"tier,omitempty"`
	EffectiveTier string                 `json:"effective_tier,omitempty"`
	Description   string                 `json:"description,omitempty"`
	InputSchema   map[string]interface{} `json:"input_schema,omitempty"`
}

type ServiceDTO struct {
	Name            string            `json:"name"`
	Type            string            `json:"type,omitempty"`
	Command         string            `json:"command,omitempty"`
	Args            []string          `json:"args,omitempty"`
	Env             map[string]string `json:"env,omitempty"`
	Port            int               `json:"port,omitempty"`
	ExtraPorts      []int             `json:"extra_ports,omitempty"`
	BuildCmd        string            `json:"build_cmd,omitempty"`
	BuildCommand    string            `json:"build_command,omitempty"`
	NoBuildOnStart  bool              `json:"no_build_on_start,omitempty"`
	Deps            []string          `json:"deps,omitempty"`
	Description     string            `json:"description,omitempty"`
	HealthURL       string            `json:"health_url,omitempty"`
	DashboardURL    string            `json:"dashboard_url,omitempty"`
	MCPURL          string            `json:"mcp_url,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	Log             string            `json:"log,omitempty"`
	ActiveTier      string            `json:"active_tier,omitempty"`
	CodeMode        bool              `json:"code_mode,omitempty"`
	NoCodeMode      bool              `json:"no_code_mode,omitempty"`
	ToolsConfigPath string            `json:"tools_config_path,omitempty"`
	Path            string            `json:"path,omitempty"`
	Installer       string            `json:"installer,omitempty"`
	Tools           []ToolDTO         `json:"tools,omitempty"`

	// Computed / live fields
	Mode               string          `json:"mode"`
	State              string          `json:"state"`
	Detail             string          `json:"detail"`
	Running            bool            `json:"running"`
	PID                int             `json:"pid,omitempty"`
	Uptime             string          `json:"uptime,omitempty"`
	Stale              bool            `json:"stale,omitempty"`
	OrphanPID          int             `json:"orphan_pid,omitempty"`
	OrphanPort         int             `json:"orphan_port,omitempty"`
	MCPEndpoint        string          `json:"mcp_endpoint,omitempty"`
	ResolvedLogPath    string          `json:"resolved_log_path"`
	ResolvedToolsPath  string          `json:"resolved_tools_path"`
	EffectiveInstaller string          `json:"effective_installer"`
	EffectiveBuild     string          `json:"effective_build,omitempty"`
	OverridesCount     int             `json:"overrides_count"`
	Registered         map[string]bool `json:"registered"`
}

type StateResponse struct {
	ConfigExists bool         `json:"config_exists"`
	ConfigError  string       `json:"config_error,omitempty"`
	Agents       AgentsDTO    `json:"agents"`
	Services     []ServiceDTO `json:"services"`
}

type ServiceUpsertRequest struct {
	Name            string            `json:"name"`
	Type            string            `json:"type"`
	Command         string            `json:"command"`
	Args            []string          `json:"args"`
	Env             map[string]string `json:"env"`
	Port            int               `json:"port"`
	ExtraPorts      []int             `json:"extra_ports"`
	BuildCommand    string            `json:"build_command"`
	NoBuildOnStart  bool              `json:"no_build_on_start"`
	Deps            []string          `json:"deps"`
	Description     string            `json:"description"`
	HealthURL       string            `json:"health_url"`
	DashboardURL    string            `json:"dashboard_url"`
	MCPURL          string            `json:"mcp_url"`
	Headers         map[string]string `json:"headers"`
	Log             string            `json:"log"`
	ActiveTier      string            `json:"active_tier"`
	CodeMode        bool              `json:"code_mode"`
	NoCodeMode      bool              `json:"no_code_mode"`
	ToolsConfigPath string            `json:"tools_config_path"`
	Path            string            `json:"path"`
	Installer       string            `json:"installer"`
}

type ActionResponse struct {
	OK     bool        `json:"ok"`
	Output string      `json:"output"`
	Error  string      `json:"error,omitempty"`
	Status *ServiceDTO `json:"status,omitempty"`
}

type HealthResponse struct {
	OK         bool   `json:"ok"`
	URL        string `json:"url"`
	StatusCode int    `json:"status_code,omitempty"`
	Body       string `json:"body,omitempty"`
	LatencyMs  int64  `json:"latency_ms"`
	Error      string `json:"error,omitempty"`
}

type ToolsResponse struct {
	Name            string    `json:"name"`
	ActiveTier      string    `json:"active_tier,omitempty"`
	CodeMode        bool      `json:"code_mode"`
	NoCodeMode      bool      `json:"no_code_mode"`
	ToolsConfigPath string    `json:"tools_config_path"`
	Tools           []ToolDTO `json:"tools"`
	OverridesJSON   string    `json:"overrides_json"`
	OverridesCount  int       `json:"overrides_count"`
}

type ToolsUpdateRequest struct {
	ActiveTier      *string   `json:"active_tier"`
	CodeMode        *bool     `json:"code_mode"`
	NoCodeMode      *bool     `json:"no_code_mode"`
	ToolsConfigPath *string   `json:"tools_config_path"`
	Tools           []ToolDTO `json:"tools"`
	Apply           *bool     `json:"apply"`
}

type ValidationIssue struct {
	Service string `json:"service"`
	Level   string `json:"level"` // "error" | "warning"
	Message string `json:"message"`
}

type DoctorFindingDTO struct {
	Message string `json:"message"`
	Fixable bool   `json:"fixable"`
	Fixed   bool   `json:"fixed"`
	Error   string `json:"error,omitempty"`
}

type DoctorServiceResult struct {
	Name      string             `json:"name"`
	State     string             `json:"state"`
	Detail    string             `json:"detail"`
	Findings  []DoctorFindingDTO `json:"findings"`
	Remaining int                `json:"remaining"`
}

type DoctorResponse struct {
	Services  []DoctorServiceResult `json:"services"`
	Remaining int                   `json:"remaining"`
	Output    string                `json:"output"`
}

type ImportCandidate struct {
	Name    string     `json:"name"`
	Source  string     `json:"source"`
	Exists  bool       `json:"exists"`
	Mode    string     `json:"mode"`
	Command string     `json:"command,omitempty"`
	Args    []string   `json:"args,omitempty"`
	MCPURL  string     `json:"mcp_url,omitempty"`
	Service ServiceDTO `json:"service"`
}

type ImportApplyRequest struct {
	Items []struct {
		Source string `json:"source"`
		Name   string `json:"name"`
	} `json:"items"`
	Overwrite bool `json:"overwrite"`
}

type AgentConfigView struct {
	Host    string `json:"host"`
	Path    string `json:"path"`
	Exists  bool   `json:"exists"`
	Enabled bool   `json:"enabled"`
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

type RawConfigResponse struct {
	Path    string       `json:"path"`
	Exists  bool         `json:"exists"`
	Content string       `json:"content"`
	Backups []BackupInfo `json:"backups"`
}

type BackupInfo struct {
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
}

type PresetDTO struct {
	ID          string               `json:"id"`
	Title       string               `json:"title"`
	Category    string               `json:"category"`
	Source      string               `json:"source"`
	SourceURL   string               `json:"source_url,omitempty"`
	Description string               `json:"description"`
	Tags        []string             `json:"tags,omitempty"`
	SetupHint   string               `json:"setup_hint,omitempty"`
	Popular     bool                 `json:"popular,omitempty"`
	Service     ServiceUpsertRequest `json:"service"`
}
