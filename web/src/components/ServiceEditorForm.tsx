import React, { useEffect, useState } from 'react';
import {
  Button,
  Checkbox,
  GroupBox,
  Select,
  Tab,
  TabBody,
  Tabs,
  TextInput,
} from 'react95';
import type { ServiceDTO, ServiceUpsertRequest } from '../types';

interface KVRow {
  key: string;
  value: string;
}

export function emptyUpsertRequest(): ServiceUpsertRequest {
  return {
    name: '',
    type: '',
    command: '',
    args: [],
    env: {},
    port: 0,
    extra_ports: [],
    build_command: '',
    no_build_on_start: false,
    deps: [],
    description: '',
    health_url: '',
    dashboard_url: '',
    mcp_url: '',
    headers: {},
    log: '',
    active_tier: '',
    code_mode: false,
    no_code_mode: false,
    tools_config_path: '',
    path: '',
    installer: '',
  };
}

export function dtoToUpsertRequest(svc: ServiceDTO): ServiceUpsertRequest {
  return {
    name: svc.name || '',
    type: svc.type || '',
    command: svc.command || '',
    args: svc.args ? [...svc.args] : [],
    env: svc.env ? { ...svc.env } : {},
    port: svc.port || 0,
    extra_ports: svc.extra_ports ? [...svc.extra_ports] : [],
    build_command: svc.effective_build || svc.build_command || svc.build_cmd || '',
    no_build_on_start: !!svc.no_build_on_start,
    deps: svc.deps ? [...svc.deps] : [],
    description: svc.description || '',
    health_url: svc.health_url || '',
    dashboard_url: svc.dashboard_url || '',
    mcp_url: svc.mcp_url || '',
    headers: svc.headers ? { ...svc.headers } : {},
    log: svc.log || '',
    active_tier: svc.active_tier || '',
    code_mode: !!svc.code_mode,
    no_code_mode: !!svc.no_code_mode,
    tools_config_path: svc.tools_config_path || '',
    path: svc.path || '',
    installer: svc.installer || '',
  };
}

function mapToRows(m: Record<string, string> | undefined): KVRow[] {
  if (!m) return [];
  return Object.entries(m).map(([key, value]) => ({ key, value }));
}

function rowsToMap(rows: KVRow[]): Record<string, string> {
  const out: Record<string, string> = {};
  for (const r of rows) {
    const k = r.key.trim();
    if (k) out[k] = r.value;
  }
  return out;
}

const typeOptions = [
  { value: '', label: 'Auto-detect' },
  { value: 'http', label: 'http (Streamable HTTP)' },
  { value: 'stdio', label: 'stdio (Agent-launched)' },
  { value: 'filesystem', label: 'filesystem (Embedded FS)' },
];

const tierOptions = [
  { value: '', label: '(default)' },
  { value: 'core', label: 'core' },
  { value: 'extended', label: 'extended' },
  { value: 'complete', label: 'complete' },
];

const installerOptions = [
  { value: '', label: 'Auto-detect' },
  { value: 'ast-mcp', label: 'ast-mcp (delegated)' },
  { value: 'native', label: 'native (mcp-local)' },
];

interface ServiceEditorFormProps {
  initial: ServiceUpsertRequest;
  isNew?: boolean;
  saving?: boolean;
  submitLabel?: string;
  onSubmit: (req: ServiceUpsertRequest) => void | Promise<void>;
  onCancel?: () => void;
  extraButtons?: React.ReactNode;
}

export const ServiceEditorForm: React.FC<ServiceEditorFormProps> = ({
  initial,
  isNew = false,
  saving = false,
  submitLabel = 'Save Configuration',
  onSubmit,
  onCancel,
  extraButtons,
}) => {
  const [form, setForm] = useState<ServiceUpsertRequest>(initial);
  const [portText, setPortText] = useState(initial.port ? String(initial.port) : '');
  const [extraPortsText, setExtraPortsText] = useState((initial.extra_ports || []).join(', '));
  const [argsRows, setArgsRows] = useState<string[]>(initial.args || []);
  const [depsRows, setDepsRows] = useState<string[]>(initial.deps || []);
  const [envRows, setEnvRows] = useState<KVRow[]>(() => mapToRows(initial.env));
  const [headerRows, setHeaderRows] = useState<KVRow[]>(() => mapToRows(initial.headers));
  const [showSecrets, setShowSecrets] = useState(false);
  const [subTab, setSubTab] = useState(0);
  const [localErr, setLocalErr] = useState('');

  useEffect(() => {
    setForm(initial);
    setPortText(initial.port ? String(initial.port) : '');
    setExtraPortsText((initial.extra_ports || []).join(', '));
    setArgsRows(initial.args || []);
    setDepsRows(initial.deps || []);
    setEnvRows(mapToRows(initial.env));
    setHeaderRows(mapToRows(initial.headers));
    setLocalErr('');
  }, [initial]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setLocalErr('');
    const portNum = portText.trim() === '' ? 0 : Number(portText.trim());
    if (Number.isNaN(portNum) || portNum < 0 || portNum > 65535) {
      setLocalErr('Port must be a number between 0 and 65535.');
      return;
    }
    const extraPorts: number[] = [];
    if (extraPortsText.trim() !== '') {
      for (const part of extraPortsText.split(',')) {
        const p = Number(part.trim());
        if (Number.isNaN(p) || p <= 0 || p > 65535) {
          setLocalErr(`Invalid extra port: "${part.trim()}"`);
          return;
        }
        extraPorts.push(p);
      }
    }
    const payload: ServiceUpsertRequest = {
      ...form,
      name: form.name.trim(),
      port: portNum,
      extra_ports: extraPorts,
      args: argsRows.filter((a) => a !== ''),
      deps: depsRows.map((d) => d.trim()).filter(Boolean),
      env: rowsToMap(envRows),
      headers: rowsToMap(headerRows),
    };
    if (!payload.name) {
      setLocalErr('Service name is required.');
      return;
    }
    await onSubmit(payload);
  };

  return (
    <form onSubmit={handleSave} style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
      <Tabs value={subTab} onChange={(val) => setSubTab(val)}>
        <Tab value={0}>General & Command</Tab>
        <Tab value={1}>Network & URLs</Tab>
        <Tab value={2}>Env & Headers ({envRows.length + headerRows.length})</Tab>
        <Tab value={3}>AST & Advanced</Tab>
      </Tabs>
      <TabBody style={{ padding: 12 }}>
        {subTab === 0 && (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(230px, 1fr))', gap: 10 }}>
            <div>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                Service Name *
              </label>
              <TextInput
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="e.g. ast-context-cache"
                fullWidth
              />
            </div>

            <div>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                Transport / Type
              </label>
              <Select
                options={typeOptions}
                value={form.type}
                onChange={(opt) => setForm({ ...form, type: opt.value })}
                width="100%"
              />
            </div>

            <div style={{ gridColumn: '1 / -1' }}>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                Description
              </label>
              <TextInput
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
                placeholder="Optional human-readable summary"
                fullWidth
              />
            </div>

            {form.type === 'filesystem' ? (
              <div style={{ gridColumn: '1 / -1' }}>
                <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                  Root Directory Path * (supports ~/ and ${'{HOME}'})
                </label>
                <TextInput
                  value={form.path}
                  onChange={(e) => setForm({ ...form, path: e.target.value })}
                  placeholder="~/projects"
                  fullWidth
                />
              </div>
            ) : (
              <>
                <div style={{ gridColumn: '1 / -1' }}>
                  <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                    Command (binary path or npx/uvx; blank for remote-only)
                  </label>
                  <TextInput
                    value={form.command}
                    onChange={(e) => setForm({ ...form, command: e.target.value })}
                    placeholder="${HOME}/git/ast-context-cache/ast-mcp or npx"
                    fullWidth
                  />
                </div>

                <div style={{ gridColumn: '1 / -1' }}>
                  <GroupBox label="Arguments (args)">
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                      {argsRows.length === 0 && (
                        <div style={{ fontSize: 12, opacity: 0.75 }}>No arguments configured.</div>
                      )}
                      {argsRows.map((arg, idx) => (
                        <div key={idx} style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                          <TextInput
                            value={arg}
                            onChange={(e) => {
                              const next = [...argsRows];
                              next[idx] = e.target.value;
                              setArgsRows(next);
                            }}
                            placeholder={`arg[${idx}]`}
                            fullWidth
                          />
                          <Button
                            type="button"
                            size="sm"
                            onClick={() => setArgsRows(argsRows.filter((_, i) => i !== idx))}
                          >
                            Remove
                          </Button>
                        </div>
                      ))}
                      <div>
                        <Button
                          type="button"
                          size="sm"
                          onClick={() => setArgsRows([...argsRows, ''])}
                        >
                          + Add Arg
                        </Button>
                      </div>
                    </div>
                  </GroupBox>
                </div>
              </>
            )}

            <div style={{ gridColumn: '1 / -1', display: 'flex', gap: 10, alignItems: 'flex-end', flexWrap: 'wrap' }}>
              <div style={{ flex: 1, minWidth: 220 }}>
                <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                  Build Command
                </label>
                <TextInput
                  value={form.build_command}
                  onChange={(e) => setForm({ ...form, build_command: e.target.value })}
                  placeholder="make build"
                  fullWidth
                />
              </div>
              <Checkbox
                checked={form.no_build_on_start}
                onChange={(e) => setForm({ ...form, no_build_on_start: e.target.checked })}
                label="Skip build on start (no_build_on_start)"
              />
            </div>
          </div>
        )}

        {subTab === 1 && (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(230px, 1fr))', gap: 10 }}>
            <div>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                Primary Port (0 for stdio)
              </label>
              <TextInput
                value={portText}
                onChange={(e) => setPortText(e.target.value)}
                placeholder="7821"
                fullWidth
              />
            </div>
            <div>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                Extra Ports (comma-separated)
              </label>
              <TextInput
                value={extraPortsText}
                onChange={(e) => setExtraPortsText(e.target.value)}
                placeholder="7830"
                fullWidth
              />
            </div>
            <div style={{ gridColumn: '1 / -1' }}>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                MCP Endpoint URL (mcp_url)
              </label>
              <TextInput
                value={form.mcp_url}
                onChange={(e) => setForm({ ...form, mcp_url: e.target.value })}
                placeholder="http://localhost:7821/mcp"
                fullWidth
              />
            </div>
            <div>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                Health URL (health_url)
              </label>
              <TextInput
                value={form.health_url}
                onChange={(e) => setForm({ ...form, health_url: e.target.value })}
                placeholder="http://localhost:7821/health"
                fullWidth
              />
            </div>
            <div>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                Dashboard URL (dashboard_url)
              </label>
              <TextInput
                value={form.dashboard_url}
                onChange={(e) => setForm({ ...form, dashboard_url: e.target.value })}
                placeholder="http://localhost:7830"
                fullWidth
              />
            </div>
          </div>
        )}

        {subTab === 2 && (
          <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
            <GroupBox label="Environment Variables (env)">
              <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Checkbox
                    checked={showSecrets}
                    onChange={(e) => setShowSecrets(e.target.checked)}
                    label="Reveal secret values"
                  />
                  <Button
                    type="button"
                    size="sm"
                    onClick={() => setEnvRows([...envRows, { key: '', value: '' }])}
                  >
                    + Add Env Var
                  </Button>
                </div>
                {envRows.length === 0 && (
                  <div style={{ fontSize: 12, opacity: 0.75 }}>No environment variables configured.</div>
                )}
                {envRows.map((row, idx) => {
                  const looksSecret = /token|secret|key|pass|auth/i.test(row.key);
                  return (
                    <div key={idx} style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                      <TextInput
                        value={row.key}
                        onChange={(e) => {
                          const next = [...envRows];
                          next[idx] = { ...row, key: e.target.value };
                          setEnvRows(next);
                        }}
                        placeholder="KEY"
                        style={{ flex: 1 }}
                      />
                      <TextInput
                        type={looksSecret && !showSecrets ? 'password' : 'text'}
                        value={row.value}
                        onChange={(e) => {
                          const next = [...envRows];
                          next[idx] = { ...row, value: e.target.value };
                          setEnvRows(next);
                        }}
                        placeholder="VALUE"
                        style={{ flex: 2 }}
                      />
                      <Button
                        type="button"
                        size="sm"
                        onClick={() => setEnvRows(envRows.filter((_, i) => i !== idx))}
                      >
                        ×
                      </Button>
                    </div>
                  );
                })}
              </div>
            </GroupBox>

            <GroupBox label="HTTP Headers (headers)">
              <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                {headerRows.length === 0 && (
                  <div style={{ fontSize: 12, opacity: 0.75 }}>No custom HTTP headers configured.</div>
                )}
                {headerRows.map((row, idx) => (
                  <div key={idx} style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                    <TextInput
                      value={row.key}
                      onChange={(e) => {
                        const next = [...headerRows];
                        next[idx] = { ...row, key: e.target.value };
                        setHeaderRows(next);
                      }}
                      placeholder="Authorization"
                      style={{ flex: 1 }}
                    />
                    <TextInput
                      type={showSecrets ? 'text' : 'password'}
                      value={row.value}
                      onChange={(e) => {
                        const next = [...headerRows];
                        next[idx] = { ...row, value: e.target.value };
                        setHeaderRows(next);
                      }}
                      placeholder="Bearer ..."
                      style={{ flex: 2 }}
                    />
                    <Button
                      type="button"
                      size="sm"
                      onClick={() => setHeaderRows(headerRows.filter((_, i) => i !== idx))}
                    >
                      ×
                    </Button>
                  </div>
                ))}
                <div>
                  <Button
                    type="button"
                    size="sm"
                    onClick={() => setHeaderRows([...headerRows, { key: '', value: '' }])}
                  >
                    + Add Header
                  </Button>
                </div>
              </div>
            </GroupBox>
          </div>
        )}

        {subTab === 3 && (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(230px, 1fr))', gap: 10 }}>
            <div>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                Log File Path (blank = ~/.mcp-local/logs/&lt;name&gt;.log)
              </label>
              <TextInput
                value={form.log}
                onChange={(e) => setForm({ ...form, log: e.target.value })}
                placeholder="~/logs/service.log"
                fullWidth
              />
            </div>

            <div>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                Installer
              </label>
              <Select
                options={installerOptions}
                value={form.installer}
                onChange={(opt) => setForm({ ...form, installer: opt.value })}
                width="100%"
              />
            </div>

            <div>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                AST Active Tier (active_tier)
              </label>
              <Select
                options={tierOptions}
                value={form.active_tier}
                onChange={(opt) => setForm({ ...form, active_tier: opt.value })}
                width="100%"
              />
            </div>

            <div>
              <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
                Tools Config Path (tools_config_path)
              </label>
              <TextInput
                value={form.tools_config_path}
                onChange={(e) => setForm({ ...form, tools_config_path: e.target.value })}
                placeholder="~/.astcache/tools.json"
                fullWidth
              />
            </div>

            <div style={{ gridColumn: '1 / -1', display: 'flex', gap: 16, flexWrap: 'wrap' }}>
              <Checkbox
                checked={form.code_mode}
                onChange={(e) =>
                  setForm({
                    ...form,
                    code_mode: e.target.checked,
                    no_code_mode: e.target.checked ? false : form.no_code_mode,
                  })
                }
                label="Enable Code Mode (AST_MCP_CODE_MODE=true)"
              />
              <Checkbox
                checked={form.no_code_mode}
                onChange={(e) =>
                  setForm({
                    ...form,
                    no_code_mode: e.target.checked,
                    code_mode: e.target.checked ? false : form.code_mode,
                  })
                }
                label="Disable Code Mode (no_code_mode)"
              />
            </div>

            <div style={{ gridColumn: '1 / -1' }}>
              <GroupBox label="Required File Dependencies (deps)">
                <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                  {depsRows.length === 0 && (
                    <div style={{ fontSize: 12, opacity: 0.75 }}>No dependency paths checked before start.</div>
                  )}
                  {depsRows.map((dep, idx) => (
                    <div key={idx} style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                      <TextInput
                        value={dep}
                        onChange={(e) => {
                          const next = [...depsRows];
                          next[idx] = e.target.value;
                          setDepsRows(next);
                        }}
                        placeholder="/path/to/required/file"
                        fullWidth
                      />
                      <Button
                        type="button"
                        size="sm"
                        onClick={() => setDepsRows(depsRows.filter((_, i) => i !== idx))}
                      >
                        Remove
                      </Button>
                    </div>
                  ))}
                  <div>
                    <Button
                      type="button"
                      size="sm"
                      onClick={() => setDepsRows([...depsRows, ''])}
                    >
                      + Add Dependency
                    </Button>
                  </div>
                </div>
              </GroupBox>
            </div>
          </div>
        )}
      </TabBody>

      {localErr && (
        <div style={{ color: '#c00000', fontWeight: 'bold', fontSize: 12 }}>{localErr}</div>
      )}

      <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8, flexWrap: 'wrap' }}>
        {extraButtons}
        {onCancel && (
          <Button type="button" onClick={onCancel} disabled={saving}>
            Cancel
          </Button>
        )}
        <Button type="submit" primary disabled={saving}>
          {saving ? 'Saving...' : isNew ? 'Create Service' : submitLabel}
        </Button>
      </div>
    </form>
  );
};
