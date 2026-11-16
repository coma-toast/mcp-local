export interface InfoResponse {
  version: string;
  config_path: string;
  config_exists: boolean;
  config_is_symlink: boolean;
  config_symlink_target?: string;
  agent_paths: Record<string, string>;
  agent_files_exist: Record<string, boolean>;
}

export interface AgentsDTO {
  managed: boolean;
  opencode: boolean;
  cursor: boolean;
  claude: boolean;
}

export interface ToolDTO {
  name: string;
  enabled: boolean;
  tier?: string;
  effective_tier?: string;
  description?: string;
  input_schema?: Record<string, unknown>;
}

export interface ServiceDTO {
  name: string;
  type?: string;
  command?: string;
  args?: string[];
  env?: Record<string, string>;
  port?: number;
  extra_ports?: number[];
  build_cmd?: string;
  build_command?: string;
  no_build_on_start?: boolean;
  deps?: string[];
  description?: string;
  health_url?: string;
  dashboard_url?: string;
  mcp_url?: string;
  headers?: Record<string, string>;
  log?: string;
  active_tier?: string;
  code_mode?: boolean;
  no_code_mode?: boolean;
  tools_config_path?: string;
  path?: string;
  installer?: string;
  tools?: ToolDTO[];

  mode: 'daemon' | 'embedded' | 'remote' | 'stdio' | string;
  state: 'running' | 'degraded' | 'orphaned' | 'port conflict' | 'stopped' | 'stale' | string;
  detail: string;
  running: boolean;
  pid?: number;
  uptime?: string;
  stale?: boolean;
  orphan_pid?: number;
  orphan_port?: number;
  mcp_endpoint?: string;
  resolved_log_path: string;
  resolved_tools_path: string;
  effective_installer: string;
  effective_build?: string;
  overrides_count: number;
  registered: Record<string, boolean>;
}

export interface StateResponse {
  config_exists: boolean;
  config_error?: string;
  agents: AgentsDTO;
  services: ServiceDTO[];
}

export interface ServiceUpsertRequest {
  name: string;
  type: string;
  command: string;
  args: string[];
  env: Record<string, string>;
  port: number;
  extra_ports: number[];
  build_command: string;
  no_build_on_start: boolean;
  deps: string[];
  description: string;
  health_url: string;
  dashboard_url: string;
  mcp_url: string;
  headers: Record<string, string>;
  log: string;
  active_tier: string;
  code_mode: boolean;
  no_code_mode: boolean;
  tools_config_path: string;
  path: string;
  installer: string;
}

export interface ActionResponse {
  ok: boolean;
  output: string;
  error?: string;
  status?: ServiceDTO;
}

export interface HealthResponse {
  ok: boolean;
  url?: string;
  status_code?: number;
  body?: string;
  latency_ms: number;
  error?: string;
}

export interface ToolsResponse {
  name: string;
  active_tier?: string;
  code_mode: boolean;
  no_code_mode: boolean;
  tools_config_path: string;
  tools: ToolDTO[];
  overrides_json: string;
  overrides_count: number;
}

export interface ToolsUpdateRequest {
  active_tier?: string;
  code_mode?: boolean;
  no_code_mode?: boolean;
  tools_config_path?: string;
  tools?: ToolDTO[];
  apply?: boolean;
}

export interface ValidationIssue {
  service: string;
  level: 'error' | 'warning';
  message: string;
}

export interface ValidateResponse {
  valid: boolean;
  issues: ValidationIssue[];
}

export interface DoctorFindingDTO {
  message: string;
  fixable: boolean;
  fixed: boolean;
  error?: string;
}

export interface DoctorServiceResult {
  name: string;
  state: string;
  detail: string;
  findings: DoctorFindingDTO[];
  remaining: number;
}

export interface DoctorResponse {
  services: DoctorServiceResult[];
  remaining: number;
  output: string;
}

export interface ImportCandidate {
  name: string;
  source: string;
  exists: boolean;
  mode: string;
  command?: string;
  args?: string[];
  mcp_url?: string;
  service: ServiceDTO;
}

export interface AgentConfigView {
  host: string;
  path: string;
  exists: boolean;
  enabled: boolean;
  content: string;
  error?: string;
}

export interface BackupInfo {
  timestamp: string;
  path: string;
  size: number;
}

export interface RawConfigResponse {
  path: string;
  exists: boolean;
  content: string;
  backups: BackupInfo[];
}

export type PresetSource = 'github.com/mcp' | 'mcpservers.org' | 'zapier.com/mcp' | 'template' | string;

export interface ServicePreset {
  id: string;
  title: string;
  category: string;
  source: PresetSource;
  source_url?: string;
  description: string;
  tags?: string[];
  setup_hint?: string;
  popular?: boolean;
  service: ServiceUpsertRequest;
}

