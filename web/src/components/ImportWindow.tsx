import React, { useCallback, useEffect, useState } from 'react';
import {
  Button,
  Checkbox,
  Table,
  TableBody,
  TableDataCell,
  TableHead,
  TableHeadCell,
  TableRow,
} from 'react95';
import { api } from '../api';
import type { ImportCandidate } from '../types';

interface ImportWindowProps {
  onNotify: (msg: string, isError?: boolean) => void;
  onAfterMutate: () => Promise<void>;
}

export const ImportWindow: React.FC<ImportWindowProps> = ({ onNotify, onAfterMutate }) => {
  const [candidates, setCandidates] = useState<ImportCandidate[]>([]);
  const [selected, setSelected] = useState<Record<string, boolean>>({});
  const [overwrite, setOverwrite] = useState(false);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    setBusy(true);
    try {
      const res = await api.importPreview();
      const list = res.candidates || [];
      setCandidates(list);
      const initial: Record<string, boolean> = {};
      for (const c of list) {
        if (!c.exists) {
          initial[`${c.source}:${c.name}`] = true;
        }
      }
      setSelected(initial);
    } catch (err) {
      onNotify(err instanceof Error ? err.message : String(err), true);
    } finally {
      setBusy(false);
    }
  }, [onNotify]);

  useEffect(() => {
    void load();
  }, [load]);

  const handleImport = async () => {
    const items = candidates
      .filter((c) => selected[`${c.source}:${c.name}`])
      .map((c) => ({ source: c.source, name: c.name }));
    if (items.length === 0) {
      onNotify('Select at least one service to import.', true);
      return;
    }
    setBusy(true);
    try {
      const res = await api.importApply(items, overwrite);
      onNotify(`Imported ${res.imported} service(s) into config.yaml.`);
      await onAfterMutate();
      await load();
    } catch (err) {
      onNotify(err instanceof Error ? err.message : String(err), true);
    } finally {
      setBusy(false);
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 8 }}>
        <span style={{ fontSize: 12 }}>
          Discover MCP server entries already defined in OpenCode, Cursor, and Claude Desktop configs.
        </span>
        <Button size="sm" onClick={load} disabled={busy}>
          Rescan Agents
        </Button>
      </div>

      <div style={{ maxHeight: 280, overflow: 'auto', border: '2px inset #808080', background: '#fff' }}>
        <Table>
          <TableHead>
            <TableRow>
              <TableHeadCell style={{ width: 50 }}>Pick</TableHeadCell>
              <TableHeadCell style={{ width: 90 }}>Source</TableHeadCell>
              <TableHeadCell>Name</TableHeadCell>
              <TableHeadCell style={{ width: 80 }}>Mode</TableHeadCell>
              <TableHeadCell>Command / Endpoint</TableHeadCell>
              <TableHeadCell style={{ width: 100 }}>Status</TableHeadCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {candidates.length === 0 ? (
              <TableRow>
                <td colSpan={6} style={{ textAlign: 'center', padding: 16 }}>
                  {busy ? 'Scanning agent configs...' : 'No MCP entries found in OpenCode, Cursor, or Claude Desktop.'}
                </td>
              </TableRow>
            ) : (
              candidates.map((c) => {
                const key = `${c.source}:${c.name}`;
                return (
                  <TableRow key={key}>
                    <TableDataCell style={{ textAlign: 'center' }}>
                      <Checkbox
                        checked={!!selected[key]}
                        onChange={(e) => setSelected({ ...selected, [key]: e.target.checked })}
                      />
                    </TableDataCell>
                    <TableDataCell>{c.source}</TableDataCell>
                    <TableDataCell style={{ fontWeight: 'bold' }}>{c.name}</TableDataCell>
                    <TableDataCell>{c.mode}</TableDataCell>
                    <TableDataCell style={{ fontFamily: 'monospace', fontSize: 12 }}>
                      {c.mcp_url || [c.command, ...(c.args || [])].filter(Boolean).join(' ')}
                    </TableDataCell>
                    <TableDataCell>
                      {c.exists ? 'Already in config' : 'New'}
                    </TableDataCell>
                  </TableRow>
                );
              })
            )}
          </TableBody>
        </Table>
      </div>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 8 }}>
        <Checkbox
          checked={overwrite}
          onChange={(e) => setOverwrite(e.target.checked)}
          label="Overwrite services that already exist in config.yaml"
        />
        <Button primary onClick={handleImport} disabled={busy || candidates.length === 0}>
          Import Selected
        </Button>
      </div>
    </div>
  );
};
