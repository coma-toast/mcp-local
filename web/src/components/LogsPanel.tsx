import React, { useCallback, useEffect, useRef, useState } from 'react';
import { Button, Checkbox, ScrollView } from 'react95';
import { api } from '../api';
import type { ServiceDTO } from '../types';

interface LogsPanelProps {
  service: ServiceDTO;
  onNotify: (msg: string, isError?: boolean) => void;
}

export const LogsPanel: React.FC<LogsPanelProps> = ({ service, onNotify }) => {
  const [lines, setLines] = useState<string[]>([]);
  const [logPath, setLogPath] = useState(service.resolved_log_path);
  const [exists, setExists] = useState(true);
  const [follow, setFollow] = useState(false);
  const [autoScroll, setAutoScroll] = useState(true);
  const [loading, setLoading] = useState(false);
  const preRef = useRef<HTMLPreElement | null>(null);

  const loadSnapshot = useCallback(async () => {
    setLoading(true);
    try {
      const res = await api.getLogs(service.name, 250);
      setLogPath(res.path);
      setExists(res.exists);
      setLines(res.lines || []);
    } catch (err) {
      onNotify(err instanceof Error ? err.message : String(err), true);
    } finally {
      setLoading(false);
    }
  }, [service.name, onNotify]);

  useEffect(() => {
    if (follow) return;
    void loadSnapshot();
  }, [follow, loadSnapshot]);

  useEffect(() => {
    if (!follow) return;
    setLines([]);
    const es = new EventSource(`/api/services/${encodeURIComponent(service.name)}/logs/stream`);
    es.addEventListener('log', (ev) => {
      setExists(true);
      setLines((prev) => {
        const next = [...prev, (ev as MessageEvent).data as string];
        return next.length > 1000 ? next.slice(next.length - 1000) : next;
      });
    });
    es.addEventListener('info', (ev) => {
      setLines((prev) => [...prev, `[info] ${(ev as MessageEvent).data as string}`]);
    });
    es.onerror = () => {
      // Browser auto-reconnects SSE; keep quiet unless closed.
    };
    return () => es.close();
  }, [follow, service.name]);

  useEffect(() => {
    if (autoScroll && preRef.current) {
      const parent = preRef.current.parentElement;
      if (parent) {
        parent.scrollTop = parent.scrollHeight;
      }
    }
  }, [lines, autoScroll]);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
      <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center', justifyContent: 'space-between' }}>
        <span style={{ fontSize: 12, fontFamily: 'monospace', wordBreak: 'break-all' }}>
          {logPath} {!exists && '(file not created yet)'}
        </span>
        <div style={{ display: 'flex', gap: 8, alignItems: 'center', flexWrap: 'wrap' }}>
          <Checkbox
            checked={follow}
            onChange={(e) => setFollow(e.target.checked)}
            label="Follow Live (SSE)"
          />
          <Checkbox
            checked={autoScroll}
            onChange={(e) => setAutoScroll(e.target.checked)}
            label="Auto-scroll"
          />
          <Button size="sm" onClick={loadSnapshot} disabled={loading || follow}>
            Refresh
          </Button>
          <Button size="sm" onClick={() => setLines([])}>
            Clear
          </Button>
          <Button
            size="sm"
            onClick={() => {
              void navigator.clipboard?.writeText(lines.join('\n'));
              onNotify('Copied log lines to clipboard.');
            }}
          >
            Copy
          </Button>
        </div>
      </div>

      <ScrollView
        style={{
          height: 290,
          background: '#000000',
          color: '#c0ffc0',
          padding: 8,
        }}
      >
        <pre
          ref={preRef}
          style={{
            margin: 0,
            fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace',
            fontSize: 12,
            lineHeight: 1.4,
            whiteSpace: 'pre-wrap',
            wordBreak: 'break-word',
          }}
        >
          {lines.length > 0
            ? lines.join('\n')
            : loading
              ? 'Reading log file...'
              : follow
                ? 'Waiting for log output...'
                : '(empty log)'}
        </pre>
      </ScrollView>
    </div>
  );
};
