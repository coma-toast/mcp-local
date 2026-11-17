import React, { useCallback, useEffect, useMemo, useState } from 'react';
import {
  Button,
  Checkbox,
  Frame,
  GroupBox,
  ScrollView,
  Select,
  SelectNative,
  Table,
  TableBody,
  TableDataCell,
  TableHead,
  TableHeadCell,
  TableRow,
  TextInput,
} from 'react95';
import { api } from '../api';
import type { ServiceDTO, ToolDTO, ToolsResponse } from '../types';

interface ToolsPanelProps {
  service: ServiceDTO;
  onNotify: (msg: string, isError?: boolean) => void;
  onAfterMutate: () => Promise<void>;
}

const activeTierOptions = [
  { value: '', label: '(default / complete)' },
  { value: 'core', label: 'core' },
  { value: 'extended', label: 'extended' },
  { value: 'complete', label: 'complete' },
];

const toolTierOptions = [
  { value: '', label: '(default)' },
  { value: 'core', label: 'core' },
  { value: 'extended', label: 'extended' },
  { value: 'complete', label: 'complete' },
];

export const ToolsPanel: React.FC<ToolsPanelProps> = ({ service, onNotify, onAfterMutate }) => {
  const [data, setData] = useState<ToolsResponse | null>(null);
  const [activeTier, setActiveTier] = useState('');
  const [codeMode, setCodeMode] = useState(false);
  const [noCodeMode, setNoCodeMode] = useState(false);
  const [toolsPath, setToolsPath] = useState('');
  const [tools, setTools] = useState<ToolDTO[]>([]);
  const [filter, setFilter] = useState('');
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [showJSON, setShowJSON] = useState(false);
  const [dirty, setDirty] = useState(false);
  const [restartNeeded, setRestartNeeded] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const res = await api.getTools(service.name);
      setData(res);
      setActiveTier(res.active_tier || '');
      setCodeMode(res.code_mode);
      setNoCodeMode(res.no_code_mode);
      setToolsPath(service.tools_config_path || '');
      setTools(res.tools || []);
      setDirty(false);
    } catch (err) {
      onNotify(err instanceof Error ? err.message : String(err), true);
    } finally {
      setLoading(false);
    }
  }, [service.name, service.tools_config_path, onNotify]);

  useEffect(() => {
    void load();
  }, [load]);

  const filteredTools = useMemo(() => {
    const q = filter.trim().toLowerCase();
    if (!q) return tools.map((t, idx) => ({ t, idx }));
    return tools
      .map((t, idx) => ({ t, idx }))
      .filter(
        ({ t }) =>
          t.name.toLowerCase().includes(q) ||
          (t.description || '').toLowerCase().includes(q) ||
          (t.effective_tier || '').toLowerCase().includes(q),
      );
  }, [tools, filter]);

  const handleSync = async () => {
    setSaving(true);
    try {
      const res = await api.syncTools(service.name);
      setData(res);
      setTools(res.tools || []);
      setDirty(false);
      onNotify(`Synced ${res.tools.length} tool(s) from live server for ${service.name}.`);
      await onAfterMutate();
    } catch (err) {
      onNotify(err instanceof Error ? err.message : String(err), true);
    } finally {
      setSaving(false);
    }
  };

  const handleSave = async (andRestart: boolean) => {
    setSaving(true);
    try {
      const res = await api.updateTools(service.name, {
        active_tier: activeTier,
        code_mode: codeMode,
        no_code_mode: noCodeMode,
        tools_config_path: toolsPath,
        tools,
        apply: true,
      });
      setData(res);
      setTools(res.tools || []);
      setDirty(false);
      if (andRestart) {
        const restartRes = await api.restartService(service.name);
        setRestartNeeded(false);
        onNotify(
          restartRes.ok
            ? `Saved tools & restarted ${service.name}.`
            : restartRes.error || 'Restart failed',
          !restartRes.ok,
        );
      } else {
        if (service.running) setRestartNeeded(true);
        onNotify(
          `Saved ${service.name} tool policy (${res.overrides_count} override(s) written to ${res.tools_config_path}).`,
        );
      }
      await onAfterMutate();
    } catch (err) {
      onNotify(err instanceof Error ? err.message : String(err), true);
    } finally {
      setSaving(false);
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
      <GroupBox label="Global Tier & Code Mode Policy">
        <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap', alignItems: 'flex-end' }}>
          <div style={{ minWidth: 180 }}>
            <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
              Active Tier (AST_MCP_TIER)
            </label>
            <Select
              options={activeTierOptions}
              value={activeTier}
              onChange={(opt) => {
                setActiveTier(opt.value);
                setDirty(true);
              }}
              width={180}
            />
          </div>

          <div style={{ minWidth: 220, flex: 1 }}>
            <label style={{ display: 'block', fontSize: 12, marginBottom: 4, fontWeight: 'bold' }}>
              tools.json Path (resolved: {data?.tools_config_path || service.resolved_tools_path})
            </label>
            <TextInput
              value={toolsPath}
              onChange={(e) => {
                setToolsPath(e.target.value);
                setDirty(true);
              }}
              placeholder="~/.astcache/tools.json"
              fullWidth
            />
          </div>

          <div style={{ display: 'flex', gap: 12, alignItems: 'center', paddingBottom: 4 }}>
            <Checkbox
              checked={codeMode}
              onChange={(e) => {
                setCodeMode(e.target.checked);
                if (e.target.checked) setNoCodeMode(false);
                setDirty(true);
              }}
              label="Code Mode"
            />
            <Checkbox
              checked={noCodeMode}
              onChange={(e) => {
                setNoCodeMode(e.target.checked);
                if (e.target.checked) setCodeMode(false);
                setDirty(true);
              }}
              label="No Code Mode"
            />
          </div>
        </div>
      </GroupBox>

      {restartNeeded && (
        <Frame variant="well" style={{ padding: '6px 10px', background: '#ffffd0', color: '#000', fontSize: 12 }}>
          Note: <b>{service.name}</b> loads <code>tools.json</code> and <code>AST_MCP_*</code> env vars at startup only. Restart to apply changes live.
        </Frame>
      )}

      <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center', justifyContent: 'space-between' }}>
        <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
          <Button size="sm" onClick={handleSync} disabled={saving || !service.mcp_endpoint}>
            Sync from Live Server
          </Button>
          <Button size="sm" onClick={() => setShowJSON(!showJSON)}>
            {showJSON ? 'Hide tools.json' : 'Preview tools.json'}
          </Button>
          <Button size="sm" onClick={load} disabled={loading}>
            Reload
          </Button>
        </div>

        <div style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
          <TextInput
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            placeholder="Filter tools..."
            style={{ width: 180 }}
          />
          <Button size="sm" primary onClick={() => handleSave(false)} disabled={saving}>
            {dirty ? 'Save & Write tools.json *' : 'Save & Write tools.json'}
          </Button>
          {service.running && (
            <Button size="sm" onClick={() => handleSave(true)} disabled={saving}>
              Save & Restart
            </Button>
          )}
        </div>
      </div>

      {showJSON && data && (
        <GroupBox label={`tools.json Overrides (${data.overrides_count})`}>
          <ScrollView style={{ maxHeight: 160, background: '#fff', color: '#000', padding: 8 }}>
            <pre style={{ margin: 0, fontFamily: 'monospace', fontSize: 12 }}>{data.overrides_json || '{}'}</pre>
          </ScrollView>
        </GroupBox>
      )}

      <div style={{ maxHeight: 300, overflow: 'auto', border: '2px inset #808080', background: '#fff' }}>
        <Table>
          <TableHead>
            <TableRow>
              <TableHeadCell style={{ width: 64 }}>Enabled</TableHeadCell>
              <TableHeadCell>Tool Name</TableHeadCell>
              <TableHeadCell style={{ width: 120 }}>Tier Override</TableHeadCell>
              <TableHeadCell style={{ width: 95 }}>Effective</TableHeadCell>
              <TableHeadCell>Description Override</TableHeadCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredTools.length === 0 ? (
              <TableRow>
                <td colSpan={5} style={{ textAlign: 'center', padding: 16 }}>
                  {loading
                    ? 'Loading tools...'
                    : tools.length === 0
                      ? 'No tools configured in config.yaml. Click "Sync from Live Server" while the service is running.'
                      : 'No tools match your filter.'}
                </td>
              </TableRow>
            ) : (
              filteredTools.map(({ t, idx }) => (
                <TableRow key={t.name}>
                  <TableDataCell style={{ textAlign: 'center' }}>
                    <Checkbox
                      checked={t.enabled}
                      onChange={(e) => {
                        const next = [...tools];
                        next[idx] = { ...t, enabled: e.target.checked };
                        setTools(next);
                        setDirty(true);
                      }}
                    />
                  </TableDataCell>
                  <TableDataCell style={{ fontFamily: 'monospace', fontWeight: 'bold' }}>
                    {t.name}
                  </TableDataCell>
                  <TableDataCell>
                    <SelectNative
                      options={toolTierOptions}
                      value={t.tier || ''}
                      onChange={(opt) => {
                        const next = [...tools];
                        next[idx] = { ...t, tier: opt.value };
                        setTools(next);
                        setDirty(true);
                      }}
                      width="100%"
                    />
                  </TableDataCell>
                  <TableDataCell style={{ fontSize: 12 }}>
                    {t.tier || t.effective_tier || '—'}
                  </TableDataCell>
                  <TableDataCell>
                    <TextInput
                      value={t.description || ''}
                      onChange={(e) => {
                        const next = [...tools];
                        next[idx] = { ...t, description: e.target.value };
                        setTools(next);
                        setDirty(true);
                      }}
                      placeholder="(server default)"
                      fullWidth
                    />
                  </TableDataCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
};
