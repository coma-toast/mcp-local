import React, { useCallback, useEffect, useState } from 'react';
import { Button, Checkbox, GroupBox, ScrollView, Tab, TabBody, Tabs } from 'react95';
import { api } from '../api';
import type { AgentConfigView, AgentsDTO } from '../types';

interface AgentsWindowProps {
  agents: AgentsDTO;
  onNotify: (msg: string, isError?: boolean) => void;
  onAfterMutate: () => Promise<void>;
}

export const AgentsWindow: React.FC<AgentsWindowProps> = ({
  agents,
  onNotify,
  onAfterMutate,
}) => {
  const [form, setForm] = useState<AgentsDTO>(agents);
  const [configs, setConfigs] = useState<AgentConfigView[]>([]);
  const [tab, setTab] = useState(0);
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    setForm(agents);
  }, [agents]);

  const loadConfigs = useCallback(async () => {
    try {
      const res = await api.agentConfigs();
      setConfigs(res.configs || []);
    } catch (err) {
      onNotify(err instanceof Error ? err.message : String(err), true);
    }
  }, [onNotify]);

  useEffect(() => {
    void loadConfigs();
  }, [loadConfigs]);

  const handleSaveAgents = async () => {
    setBusy(true);
    try {
      await api.updateAgents(form);
      onNotify('Saved agent registration targets to config.yaml.');
      await onAfterMutate();
      await loadConfigs();
    } catch (err) {
      onNotify(err instanceof Error ? err.message : String(err), true);
    } finally {
      setBusy(false);
    }
  };

  const handleRegisterAll = async () => {
    setBusy(true);
    try {
      const res = await api.registerAll();
      onNotify(res.output || (res.ok ? 'Registered all services.' : res.error || 'Failed'), !res.ok);
      await onAfterMutate();
      await loadConfigs();
    } catch (err) {
      onNotify(err instanceof Error ? err.message : String(err), true);
    } finally {
      setBusy(false);
    }
  };

  const selectedConfig = configs[tab];

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
      <GroupBox label="Managed Agent Targets (config.yaml -> agents:)">
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          <Checkbox
            checked={form.managed}
            onChange={(e) => setForm({ ...form, managed: e.target.checked })}
            label="Explicitly manage agents: key (when unchecked, legacy mode enables all three)"
          />
          <div style={{ display: 'flex', gap: 18, flexWrap: 'wrap', paddingLeft: 8 }}>
            <Checkbox
              checked={form.managed ? form.opencode : true}
              disabled={!form.managed}
              onChange={(e) => setForm({ ...form, opencode: e.target.checked })}
              label="OpenCode"
            />
            <Checkbox
              checked={form.managed ? form.cursor : true}
              disabled={!form.managed}
              onChange={(e) => setForm({ ...form, cursor: e.target.checked })}
              label="Cursor"
            />
            <Checkbox
              checked={form.managed ? form.claude : true}
              disabled={!form.managed}
              onChange={(e) => setForm({ ...form, claude: e.target.checked })}
              label="Claude Desktop (uses mcp-local bridge for HTTP)"
            />
          </div>
          <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', flexWrap: 'wrap' }}>
            <Button size="sm" onClick={loadConfigs} disabled={busy}>
              Refresh Files
            </Button>
            <Button size="sm" onClick={handleRegisterAll} disabled={busy}>
              Register All Enabled Services
            </Button>
            <Button size="sm" primary onClick={handleSaveAgents} disabled={busy}>
              Save Agent Targets
            </Button>
          </div>
        </div>
      </GroupBox>

      {configs.length > 0 && (
        <div>
          <Tabs value={tab} onChange={(v) => setTab(v)}>
            {configs.map((c, idx) => (
              <Tab key={c.host} value={idx}>
                {c.host} {c.exists ? '' : '(missing)'}
              </Tab>
            ))}
          </Tabs>
          <TabBody style={{ padding: 10 }}>
            {selectedConfig && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                <div style={{ fontSize: 12, fontFamily: 'monospace', wordBreak: 'break-all' }}>
                  <b>Path:</b> {selectedConfig.path} | <b>Enabled:</b> {selectedConfig.enabled ? 'yes' : 'no'}
                </div>
                {selectedConfig.error ? (
                  <div style={{ color: '#c00000', fontSize: 12 }}>{selectedConfig.error}</div>
                ) : (
                  <ScrollView style={{ height: 240, background: '#fff', color: '#000', padding: 8 }}>
                    <pre style={{ margin: 0, fontFamily: 'monospace', fontSize: 12 }}>
                      {selectedConfig.exists
                        ? selectedConfig.content
                        : '(File does not exist yet on disk)'}
                    </pre>
                  </ScrollView>
                )}
              </div>
            )}
          </TabBody>
        </div>
      )}
    </div>
  );
};
