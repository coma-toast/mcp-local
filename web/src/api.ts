import type {
  ActionResponse,
  AgentConfigView,
  AgentsDTO,
  DoctorResponse,
  HealthResponse,
  ImportCandidate,
  InfoResponse,
  RawConfigResponse,
  ServiceDTO,
  ServiceUpsertRequest,
  StateResponse,
  ToolsResponse,
  ToolsUpdateRequest,
  ValidateResponse,
} from './types';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const headers: Record<string, string> = {
    ...(init?.headers as Record<string, string> | undefined),
  };
  if (init?.method && ['POST', 'PUT', 'PATCH'].includes(init.method.toUpperCase())) {
    headers['Content-Type'] = 'application/json';
  }
  const resp = await fetch(path, { ...init, headers });
  const text = await resp.text();
  let data: unknown = undefined;
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      if (!resp.ok) {
        throw new Error(text || `HTTP ${resp.status}`);
      }
    }
  }
  if (!resp.ok) {
    const errMsg =
      typeof data === 'object' && data !== null && 'error' in data
        ? String((data as { error: unknown }).error)
        : text || `HTTP ${resp.status}`;
    throw new Error(errMsg);
  }
  return data as T;
}

const enc = encodeURIComponent;

export const api = {
  info: () => request<InfoResponse>('/api/info'),
  state: () => request<StateResponse>('/api/state'),
  updateAgents: (agents: AgentsDTO) =>
    request<AgentsDTO>('/api/agents', { method: 'PUT', body: JSON.stringify(agents) }),
  agentConfigs: () => request<{ configs: AgentConfigView[] }>('/api/agents/configs'),

  createService: (svc: ServiceUpsertRequest) =>
    request<ServiceDTO>('/api/services', { method: 'POST', body: JSON.stringify(svc) }),
  updateService: (name: string, svc: ServiceUpsertRequest) =>
    request<ServiceDTO>(`/api/services/${enc(name)}`, { method: 'PUT', body: JSON.stringify(svc) }),
  deleteService: (name: string, deregister = true) =>
    request<ActionResponse>(`/api/services/${enc(name)}?deregister=${deregister ? 'true' : 'false'}`, {
      method: 'DELETE',
    }),

  startService: (name: string, build = true) =>
    request<ActionResponse>(`/api/services/${enc(name)}/start`, {
      method: 'POST',
      body: JSON.stringify({ build }),
    }),
  stopService: (name: string, deregister = true) =>
    request<ActionResponse>(`/api/services/${enc(name)}/stop`, {
      method: 'POST',
      body: JSON.stringify({ deregister }),
    }),
  restartService: (name: string) =>
    request<ActionResponse>(`/api/services/${enc(name)}/restart`, { method: 'POST', body: '{}' }),
  rebuildService: (name: string) =>
    request<ActionResponse>(`/api/services/${enc(name)}/rebuild`, { method: 'POST', body: '{}' }),
  registerService: (name: string) =>
    request<ActionResponse>(`/api/services/${enc(name)}/register`, { method: 'POST', body: '{}' }),
  deregisterService: (name: string) =>
    request<ActionResponse>(`/api/services/${enc(name)}/deregister`, { method: 'POST', body: '{}' }),
  healthService: (name: string) => request<HealthResponse>(`/api/services/${enc(name)}/health`),

  getTools: (name: string) => request<ToolsResponse>(`/api/services/${enc(name)}/tools`),
  syncTools: (name: string) =>
    request<ToolsResponse>(`/api/services/${enc(name)}/tools/sync`, { method: 'POST', body: '{}' }),
  updateTools: (name: string, body: ToolsUpdateRequest) =>
    request<ToolsResponse>(`/api/services/${enc(name)}/tools`, {
      method: 'PUT',
      body: JSON.stringify(body),
    }),
  applyTools: (name: string) =>
    request<ActionResponse>(`/api/services/${enc(name)}/tools/apply`, { method: 'POST', body: '{}' }),

  getLogs: (name: string, lines = 200) =>
    request<{ path: string; exists: boolean; lines: string[] }>(
      `/api/services/${enc(name)}/logs?lines=${lines}`,
    ),

  startAll: () => request<ActionResponse>('/api/actions/start-all', { method: 'POST', body: '{}' }),
  stopAll: (deregister = true) =>
    request<ActionResponse>('/api/actions/stop-all', {
      method: 'POST',
      body: JSON.stringify({ deregister }),
    }),
  restartAll: () => request<ActionResponse>('/api/actions/restart-all', { method: 'POST', body: '{}' }),
  registerAll: () => request<ActionResponse>('/api/actions/register-all', { method: 'POST', body: '{}' }),

  validate: () => request<ValidateResponse>('/api/validate'),
  doctor: (service = '', fix = false) =>
    request<DoctorResponse>('/api/doctor', {
      method: 'POST',
      body: JSON.stringify({ service, fix }),
    }),

  importPreview: () => request<{ candidates: ImportCandidate[] }>('/api/import/preview'),
  importApply: (items: { source: string; name: string }[], overwrite = false) =>
    request<{ imported: number }>('/api/import', {
      method: 'POST',
      body: JSON.stringify({ items, overwrite }),
    }),

  getRawConfig: () => request<RawConfigResponse>('/api/config/raw'),
  putRawConfig: (content: string) =>
    request<RawConfigResponse>('/api/config/raw', {
      method: 'PUT',
      body: JSON.stringify({ content }),
    }),

  searchRegistry: (q: string) =>
    request<{ presets: import('./types').ServicePreset[]; error?: string }>(
      `/api/registry/search?q=${enc(q)}`,
    ),
};

