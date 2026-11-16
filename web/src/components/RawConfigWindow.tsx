import React, { useCallback, useEffect, useState } from 'react';
import { Button, GroupBox, ScrollView, TextInput } from 'react95';
import { api } from '../api';
import type { RawConfigResponse } from '../types';

interface RawConfigWindowProps {
  onNotify: (msg: string, isError?: boolean) => void;
  onAfterMutate: () => Promise<void>;
}

export const RawConfigWindow: React.FC<RawConfigWindowProps> = ({ onNotify, onAfterMutate }) => {
  const [data, setData] = useState<RawConfigResponse | null>(null);
  const [content, setContent] = useState('');
  const [dirty, setDirty] = useState(false);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    setBusy(true);
    try {
      const res = await api.getRawConfig();
      setData(res);
      setContent(res.content || '');
      setDirty(false);
    } catch (err) {
      onNotify(err instanceof Error ? err.message : String(err), true);
    } finally {
      setBusy(false);
    }
  }, [onNotify]);

  useEffect(() => {
    void load();
  }, [load]);

  const handleSave = async () => {
    setBusy(true);
    try {
      const res = await api.putRawConfig(content);
      setData(res);
      setContent(res.content);
      setDirty(false);
      onNotify('Saved config.yaml (with automatic backup).');
      await onAfterMutate();
    } catch (err) {
      onNotify(err instanceof Error ? err.message : String(err), true);
    } finally {
      setBusy(false);
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 8 }}>
        <span style={{ fontSize: 12, fontFamily: 'monospace', wordBreak: 'break-all' }}>
          {data?.path || '~/.mcp-local/config.yaml'} {dirty ? ' (modified)' : ''}
        </span>
        <div style={{ display: 'flex', gap: 6 }}>
          <Button size="sm" onClick={load} disabled={busy}>
            Reload
          </Button>
          <Button size="sm" primary onClick={handleSave} disabled={busy}>
            Validate & Save YAML
          </Button>
        </div>
      </div>

      <TextInput
        multiline
        rows={16}
        value={content}
        onChange={(e) => {
          setContent(e.target.value);
          setDirty(true);
        }}
        placeholder="# Paste or edit ~/.mcp-local/config.yaml here"
        fullWidth
        style={{
          fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace',
          fontSize: 12,
        }}
      />

      {data && data.backups && data.backups.length > 0 && (
        <GroupBox label={`Recent Backups (${data.backups.length})`}>
          <ScrollView style={{ maxHeight: 90, background: '#fff', color: '#000', padding: 6 }}>
            {data.backups.map((b) => (
              <div key={b.path} style={{ fontSize: 11, fontFamily: 'monospace', marginBottom: 2 }}>
                {b.timestamp} ({b.size} bytes) — {b.path}
              </div>
            ))}
          </ScrollView>
        </GroupBox>
      )}
    </div>
  );
};
