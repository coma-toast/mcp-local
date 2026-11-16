import React, { useCallback, useEffect, useMemo, useState } from 'react';
import {
  AppBar,
  Button,
  Checkbox,
  Frame,
  GroupBox,
  MenuList,
  MenuListItem,
  ScrollView,
  Separator,
  Tab,
  TabBody,
  Tabs,
  Table,
  TableBody,
  TableDataCell,
  TableHead,
  TableHeadCell,
  TableRow,
  TextInput,
  Toolbar,
} from 'react95';
import { createGlobalStyle, ThemeProvider } from 'styled-components';
import { styleReset } from 'react95';
import ms_sans_serif from 'react95/dist/fonts/ms_sans_serif.woff2';
import ms_sans_serif_bold from 'react95/dist/fonts/ms_sans_serif_bold.woff2';
import original from 'react95/dist/themes/original';
import tokyoDark from 'react95/dist/themes/tokyoDark';
import vaporTeal from 'react95/dist/themes/vaporTeal';
import rainyDay from 'react95/dist/themes/rainyDay';
import modernDark from 'react95/dist/themes/modernDark';
import matrix from 'react95/dist/themes/matrix';
import lilac from 'react95/dist/themes/lilac';
import rose from 'react95/dist/themes/rose';
import hotdogStand from 'react95/dist/themes/hotdogStand';

import { api } from './api';
import { AddServiceModal } from './components/AddServiceModal';
import { AgentsWindow } from './components/AgentsWindow';
import { DoctorWindow } from './components/DoctorWindow';
import { ImportWindow } from './components/ImportWindow';
import { LogsPanel } from './components/LogsPanel';
import { RawConfigWindow } from './components/RawConfigWindow';
import { dtoToUpsertRequest, ServiceEditorForm } from './components/ServiceEditorForm';
import { ToolsPanel } from './components/ToolsPanel';
import { Win95Window, type WinBounds } from './components/Win95Window';
import {
  ComputerIcon,
  DoctorIcon,
  FloppyImportIcon,
  NetworkAgentsIcon,
  NotepadYamlIcon,
  StatusLed,
  ToolsGearIcon,
  WinLogoIcon,
} from './icons';
import type { HealthResponse, InfoResponse, ServiceDTO, StateResponse } from './types';

const THEMES = {
  original: { label: 'Windows 95 Classic', theme: original, bg: '#008080' },
  rainyDay: { label: 'Rainy Day', theme: rainyDay, bg: '#3e5368' },
  vaporTeal: { label: 'Vapor Teal', theme: vaporTeal, bg: '#016868' },
  lilac: { label: 'Lilac', theme: lilac, bg: '#5a4e7c' },
  rose: { label: 'Rose', theme: rose, bg: '#734a54' },
  tokyoDark: { label: 'Tokyo Dark', theme: tokyoDark, bg: '#16161e' },
  modernDark: { label: 'Modern Dark', theme: modernDark, bg: '#1e1e1e' },
  matrix: { label: 'Matrix', theme: matrix, bg: '#081808' },
  hotdogStand: { label: 'Hotdog Stand', theme: hotdogStand, bg: '#b00000' },
} as const;

type ThemeKey = keyof typeof THEMES;
type WindowId = 'manager' | 'add' | 'agents' | 'doctor' | 'import' | 'raw' | 'delete';

interface WinState {
  open: boolean;
  minimized: boolean;
  maximized: boolean;
  z: number;
  bounds: WinBounds;
}

const GlobalStyles = createGlobalStyle<{ $bg: string }>`
  ${styleReset}
  @font-face {
    font-family: 'ms_sans_serif';
    src: url('${ms_sans_serif}') format('woff2');
    font-weight: 400;
    font-style: normal;
  }
  @font-face {
    font-family: 'ms_sans_serif';
    src: url('${ms_sans_serif_bold}') format('woff2');
    font-weight: bold;
    font-style: normal;
  }
  html, body, #root {
    height: 100%;
    margin: 0;
    padding: 0;
    overflow: hidden;
    font-family: 'ms_sans_serif', sans-serif;
    background: ${(p) => p.$bg};
  }
  * {
    box-sizing: border-box;
  }
  code, pre {
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  }
`;

function formatClock(d: Date): string {
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

export default function App() {
  const [themeKey, setThemeKey] = useState<ThemeKey>(() => {
    const saved = localStorage.getItem('mcp-local-95-theme') as ThemeKey | null;
    return saved && saved in THEMES ? saved : 'original';
  });
  const [info, setInfo] = useState<InfoResponse | null>(null);
  const [state, setState] = useState<StateResponse>({
    config_exists: false,
    agents: { managed: false, opencode: true, cursor: true, claude: true },
    services: [],
  });
  const [selectedName, setSelectedName] = useState<string>('');
  const [detailTab, setDetailTab] = useState<number>(0);
  const [filter, setFilter] = useState('');
  const [autoPoll, setAutoPoll] = useState(true);
  const [busy, setBusy] = useState(false);
  const [actionLog, setActionLog] = useState<string>('');
  const [healthResult, setHealthResult] = useState<HealthResponse | null>(null);
  const [toast, setToast] = useState<{ msg: string; isError: boolean } | null>(null);
  const [startOpen, setStartOpen] = useState(false);
  const [themesOpen, setThemesOpen] = useState(false);
  const [clock, setClock] = useState(() => formatClock(new Date()));
  const [deleteTarget, setDeleteTarget] = useState<string>('');
  const [deleteDeregister, setDeleteDeregister] = useState(true);

  const [topZ, setTopZ] = useState(20);
  const [activeWin, setActiveWin] = useState<WindowId>('manager');
  const [windows, setWindows] = useState<Record<WindowId, WinState>>({
    manager: {
      open: true,
      minimized: false,
      maximized: false,
      z: 10,
      bounds: { x: 108, y: 14, width: 1040, height: 700 },
    },
    add: {
      open: false,
      minimized: false,
      maximized: false,
      z: 11,
      bounds: { x: 150, y: 30, width: 820 },
    },
    agents: {
      open: false,
      minimized: false,
      maximized: false,
      z: 12,
      bounds: { x: 190, y: 55, width: 780 },
    },
    doctor: {
      open: false,
      minimized: false,
      maximized: false,
      z: 13,
      bounds: { x: 210, y: 70, width: 780 },
    },
    import: {
      open: false,
      minimized: false,
      maximized: false,
      z: 14,
      bounds: { x: 230, y: 85, width: 760 },
    },
    raw: {
      open: false,
      minimized: false,
      maximized: false,
      z: 15,
      bounds: { x: 200, y: 45, width: 820 },
    },
    delete: {
      open: false,
      minimized: false,
      maximized: false,
      z: 19,
      bounds: { x: 340, y: 180, width: 420 },
    },
  });

  const notify = useCallback((msg: string, isError = false) => {
    setToast({ msg, isError });
  }, []);

  useEffect(() => {
    if (!toast) return;
    const t = setTimeout(() => setToast(null), 6000);
    return () => clearTimeout(t);
  }, [toast]);

  useEffect(() => {
    const t = setInterval(() => setClock(formatClock(new Date())), 15000);
    return () => clearInterval(t);
  }, []);

  const refresh = useCallback(async () => {
    try {
      const [inf, st] = await Promise.all([api.info(), api.state()]);
      setInfo(inf);
      setState(st);
      setSelectedName((prev) => {
        if (prev && st.services.some((s) => s.name === prev)) return prev;
        return st.services[0]?.name || '';
      });
    } catch (err) {
      notify(err instanceof Error ? err.message : String(err), true);
    }
  }, [notify]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  useEffect(() => {
    if (!autoPoll) return;
    const t = setInterval(() => {
      void api.state().then(setState).catch(() => {});
    }, 5000);
    return () => clearInterval(t);
  }, [autoPoll]);

  const focusWindow = useCallback(
    (id: WindowId) => {
      const nextZ = topZ + 1;
      setTopZ(nextZ);
      setActiveWin(id);
      setWindows((prev) => ({
        ...prev,
        [id]: { ...prev[id], open: true, minimized: false, z: nextZ },
      }));
    },
    [topZ],
  );

  const closeWindow = (id: WindowId) => {
    setWindows((prev) => ({
      ...prev,
      [id]: { ...prev[id], open: false, minimized: false },
    }));
  };

  const minimizeWindow = (id: WindowId) => {
    setWindows((prev) => ({
      ...prev,
      [id]: { ...prev[id], minimized: true },
    }));
  };

  const toggleMaximize = (id: WindowId) => {
    setWindows((prev) => ({
      ...prev,
      [id]: { ...prev[id], maximized: !prev[id].maximized },
    }));
    focusWindow(id);
  };

  const moveWindow = (id: WindowId, x: number, y: number) => {
    setWindows((prev) => ({
      ...prev,
      [id]: { ...prev[id], bounds: { ...prev[id].bounds, x, y } },
    }));
  };

  const selectedService: ServiceDTO | undefined = useMemo(
    () => state.services.find((s) => s.name === selectedName),
    [state.services, selectedName],
  );

  const editorInitial = useMemo(
    () => (selectedService ? dtoToUpsertRequest(selectedService) : null),
    [selectedService],
  );

  useEffect(() => {
    setHealthResult(null);
  }, [selectedName]);

  const filteredServices = useMemo(() => {
    const q = filter.trim().toLowerCase();
    if (!q) return state.services;
    return state.services.filter(
      (s) =>
        s.name.toLowerCase().includes(q) ||
        s.mode.toLowerCase().includes(q) ||
        s.state.toLowerCase().includes(q) ||
        (s.description || '').toLowerCase().includes(q) ||
        (s.command || '').toLowerCase().includes(q) ||
        (s.mcp_url || '').toLowerCase().includes(q),
    );
  }, [state.services, filter]);

  const runAction = async (label: string, fn: () => Promise<{ ok: boolean; output: string; error?: string }>) => {
    setBusy(true);
    try {
      const res = await fn();
      setActionLog(res.output || res.error || label);
      notify(res.ok ? `${label}: OK` : res.error || `${label} failed`, !res.ok);
      await refresh();
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err);
      setActionLog(msg);
      notify(msg, true);
    } finally {
      setBusy(false);
    }
  };

  const handleHealthCheck = async (name: string) => {
    setBusy(true);
    try {
      const res = await api.healthService(name);
      setHealthResult(res);
      notify(
        res.ok
          ? `Health OK (${res.status_code}, ${res.latency_ms}ms)`
          : res.error || `Health HTTP ${res.status_code}`,
        !res.ok,
      );
    } catch (err) {
      notify(err instanceof Error ? err.message : String(err), true);
    } finally {
      setBusy(false);
    }
  };

  const runningCount = state.services.filter((s) => s.running).length;
  const activeTheme = THEMES[themeKey];

  const windowTitles: Record<WindowId, string> = {
    manager: 'MCP Control Panel',
    add: 'Add New MCP Service',
    agents: 'Agent Sync & IDE Configs',
    doctor: 'Doctor & Config Validator',
    import: 'Import MCP Servers',
    raw: 'Notepad — config.yaml',
    delete: 'Confirm Delete',
  };

  const desktopShortcuts: { id: WindowId; label: string; icon: React.ReactNode }[] = [
    { id: 'manager', label: 'MCP Manager', icon: <ComputerIcon size={32} /> },
    { id: 'agents', label: 'Agent Sync', icon: <NetworkAgentsIcon size={32} /> },
    { id: 'doctor', label: 'Doctor', icon: <DoctorIcon size={32} /> },
    { id: 'import', label: 'Import IDEs', icon: <FloppyImportIcon size={32} /> },
    { id: 'raw', label: 'config.yaml', icon: <NotepadYamlIcon size={32} /> },
  ];

  const agentBadge = (svc: ServiceDTO, host: string, short: string) => {
    const enabled =
      !state.agents.managed ||
      (host === 'opencode' && state.agents.opencode) ||
      (host === 'cursor' && state.agents.cursor) ||
      (host === 'claude' && state.agents.claude);
    if (!enabled) {
      return (
        <span key={host} title={`${host}: disabled in agents config`} style={{ opacity: 0.35, fontSize: 11 }}>
          {short}
        </span>
      );
    }
    const ok = !!svc.registered?.[host];
    return (
      <span
        key={host}
        title={`${host}: ${ok ? 'registered' : 'not registered'}`}
        style={{
          fontSize: 10,
          fontWeight: 'bold',
          padding: '1px 4px',
          border: '1px solid #000',
          background: ok ? '#c8f7c5' : '#ffd6d6',
          color: '#000',
        }}
      >
        {short}
      </span>
    );
  };

  return (
    <ThemeProvider theme={activeTheme.theme}>
      <GlobalStyles $bg={activeTheme.bg} />

      <div
        style={{ position: 'fixed', inset: 0, bottom: 48, overflow: 'hidden' }}
        onMouseDown={() => {
          if (startOpen) {
            setStartOpen(false);
            setThemesOpen(false);
          }
        }}
      >
        {/* Desktop Icons */}
        <div
          style={{
            position: 'absolute',
            top: 14,
            left: 12,
            display: 'flex',
            flexDirection: 'column',
            gap: 16,
            zIndex: 1,
          }}
        >
          {desktopShortcuts.map((sc) => (
            <button
              key={sc.id}
              type="button"
              onClick={() => focusWindow(sc.id)}
              style={{
                width: 78,
                background: 'transparent',
                border: 'none',
                color: '#ffffff',
                textShadow: '1px 1px 0 #000000',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                gap: 4,
                cursor: 'pointer',
                padding: 4,
                fontSize: 12,
                fontFamily: 'inherit',
              }}
            >
              {sc.icon}
              <span>{sc.label}</span>
            </button>
          ))}
        </div>

        {/* Main Window: MCP Control Panel */}
        {windows.manager.open && (
          <Win95Window
            id="manager"
            title={`MCP Control Panel — ${info?.config_path || '~/.mcp-local/config.yaml'}`}
            icon={<ComputerIcon size={16} />}
            bounds={windows.manager.bounds}
            zIndex={windows.manager.z}
            active={activeWin === 'manager'}
            minimized={windows.manager.minimized}
            maximized={windows.manager.maximized}
            onFocus={() => focusWindow('manager')}
            onClose={() => closeWindow('manager')}
            onMinimize={() => minimizeWindow('manager')}
            onToggleMaximize={() => toggleMaximize('manager')}
            onMove={(x, y) => moveWindow('manager', x, y)}
            footer={
              <Frame
                variant="well"
                style={{
                  padding: '3px 8px',
                  fontSize: 12,
                  display: 'flex',
                  justifyContent: 'space-between',
                  gap: 12,
                  flexShrink: 0,
                }}
              >
                <span>
                  {toast ? (
                    <span style={{ color: toast.isError ? '#c00000' : '#006400', fontWeight: 'bold' }}>
                      {toast.msg}
                    </span>
                  ) : (
                    <span>
                      {state.services.length} service(s) | {runningCount} running
                      {info?.config_is_symlink ? ` | symlink -> ${info.config_symlink_target}` : ''}
                    </span>
                  )}
                </span>
                <span>v{info?.version || 'dev'}</span>
              </Frame>
            }
          >
            <div style={{ display: 'flex', flexDirection: 'column', gap: 10, height: '100%' }}>
              {/* Toolbar */}
              <Toolbar style={{ flexWrap: 'wrap', gap: 6, justifyContent: 'space-between' }}>
                <div style={{ display: 'flex', gap: 4, flexWrap: 'wrap', alignItems: 'center' }}>
                  <Button size="sm" primary onClick={() => focusWindow('add')}>
                    + New Service
                  </Button>
                  <Button size="sm" onClick={() => runAction('Start All', api.startAll)} disabled={busy}>
                    Start All
                  </Button>
                  <Button size="sm" onClick={() => runAction('Stop All', () => api.stopAll(true))} disabled={busy}>
                    Stop All
                  </Button>
                  <Button size="sm" onClick={() => runAction('Restart All', api.restartAll)} disabled={busy}>
                    Restart All
                  </Button>
                  <Button size="sm" onClick={() => runAction('Register All', api.registerAll)} disabled={busy}>
                    Register All
                  </Button>
                  <Button size="sm" onClick={() => focusWindow('agents')}>
                    Agents
                  </Button>
                  <Button size="sm" onClick={() => focusWindow('doctor')}>
                    Doctor
                  </Button>
                  <Button size="sm" onClick={() => focusWindow('import')}>
                    Import
                  </Button>
                  <Button size="sm" onClick={() => focusWindow('raw')}>
                    YAML
                  </Button>
                </div>

                <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                  <TextInput
                    value={filter}
                    onChange={(e) => setFilter(e.target.value)}
                    placeholder="Search services..."
                    style={{ width: 150 }}
                  />
                  <Checkbox
                    checked={autoPoll}
                    onChange={(e) => setAutoPoll(e.target.checked)}
                    label="Poll 5s"
                  />
                  <Button size="sm" onClick={refresh} disabled={busy}>
                    Refresh
                  </Button>
                </div>
              </Toolbar>

              {state.config_error && (
                <Frame variant="well" style={{ padding: 8, background: '#ffe0e0', color: '#900000', fontSize: 12 }}>
                  <b>Config Error:</b> {state.config_error}{' '}
                  <Button size="sm" onClick={() => focusWindow('raw')} style={{ marginLeft: 8 }}>
                    Open config.yaml
                  </Button>
                </Frame>
              )}

              {/* Services Table */}
              <div
                style={{
                  minHeight: 150,
                  maxHeight: 225,
                  overflow: 'auto',
                  border: '2px inset #808080',
                  background: '#fff',
                  flexShrink: 0,
                }}
              >
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableHeadCell>Service</TableHeadCell>
                      <TableHeadCell style={{ width: 90 }}>Mode</TableHeadCell>
                      <TableHeadCell style={{ width: 125 }}>Status</TableHeadCell>
                      <TableHeadCell style={{ width: 70 }}>Port</TableHeadCell>
                      <TableHeadCell style={{ width: 75 }}>PID</TableHeadCell>
                      <TableHeadCell style={{ width: 118 }}>Agents</TableHeadCell>
                      <TableHeadCell style={{ width: 65 }}>Tools</TableHeadCell>
                      <TableHeadCell style={{ width: 160 }}>Quick Actions</TableHeadCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {filteredServices.length === 0 ? (
                      <TableRow>
                        <td colSpan={8} style={{ textAlign: 'center', padding: 20 }}>
                          {state.services.length === 0 ? (
                            <div style={{ display: 'flex', flexDirection: 'column', gap: 8, alignItems: 'center' }}>
                              <span>
                                No MCP services configured yet
                                {!state.config_exists ? ' (~/.mcp-local/config.yaml will be created on save).' : '.'}
                              </span>
                              <div style={{ display: 'flex', gap: 8 }}>
                                <Button size="sm" primary onClick={() => focusWindow('add')}>
                                  + Add First Service
                                </Button>
                                <Button size="sm" onClick={() => focusWindow('import')}>
                                  Import from IDEs
                                </Button>
                              </div>
                            </div>
                          ) : (
                            'No services match your search.'
                          )}
                        </td>
                      </TableRow>
                    ) : (
                      filteredServices.map((svc) => {
                        const isSel = svc.name === selectedName;
                        return (
                          <TableRow
                            key={svc.name}
                            onClick={() => setSelectedName(svc.name)}
                            style={{
                              cursor: 'pointer',
                              background: isSel ? '#000080' : undefined,
                              color: isSel ? '#ffffff' : undefined,
                            }}
                          >
                            <TableDataCell style={{ fontWeight: 'bold' }}>{svc.name}</TableDataCell>
                            <TableDataCell>{svc.mode}</TableDataCell>
                            <TableDataCell>
                              <StatusLed state={svc.state} />
                            </TableDataCell>
                            <TableDataCell>{svc.port || '—'}</TableDataCell>
                            <TableDataCell>{svc.pid || '—'}</TableDataCell>
                            <TableDataCell>
                              <div style={{ display: 'flex', gap: 3 }}>
                                {agentBadge(svc, 'opencode', 'OC')}
                                {agentBadge(svc, 'cursor', 'CR')}
                                {agentBadge(svc, 'claude', 'CL')}
                              </div>
                            </TableDataCell>
                            <TableDataCell>{svc.tools?.length || 0}</TableDataCell>
                            <TableDataCell onClick={(e) => e.stopPropagation()}>
                              <div style={{ display: 'flex', gap: 4 }}>
                                {svc.running ? (
                                  <Button
                                    size="sm"
                                    disabled={busy}
                                    onClick={() =>
                                      runAction(`Stop ${svc.name}`, () => api.stopService(svc.name, true))
                                    }
                                  >
                                    Stop
                                  </Button>
                                ) : (
                                  <Button
                                    size="sm"
                                    disabled={busy}
                                    onClick={() =>
                                      runAction(`Start ${svc.name}`, () => api.startService(svc.name, true))
                                    }
                                  >
                                    Start
                                  </Button>
                                )}
                                <Button
                                  size="sm"
                                  disabled={busy}
                                  onClick={() =>
                                    runAction(`Restart ${svc.name}`, () => api.restartService(svc.name))
                                  }
                                >
                                  Restart
                                </Button>
                              </div>
                            </TableDataCell>
                          </TableRow>
                        );
                      })
                    )}
                  </TableBody>
                </Table>
              </div>

              {/* Detail Inspector */}
              {selectedService ? (
                <div style={{ flex: 1, minHeight: 0, display: 'flex', flexDirection: 'column' }}>
                  <Tabs value={detailTab} onChange={(v) => setDetailTab(v)}>
                    <Tab value={0}>Overview & Controls</Tab>
                    <Tab value={1}>Edit Configuration</Tab>
                    <Tab value={2}>
                      Tools & Tiers ({selectedService.tools?.length || 0})
                    </Tab>
                    <Tab value={3}>Logs</Tab>
                  </Tabs>
                  <TabBody style={{ flex: 1, overflow: 'auto', padding: 10 }}>
                    {detailTab === 0 && (
                      <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
                        <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap', alignItems: 'center', justifyContent: 'space-between' }}>
                          <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
                            <Button
                              size="sm"
                              primary
                              disabled={busy}
                              onClick={() =>
                                runAction(`Start ${selectedService.name}`, () =>
                                  api.startService(selectedService.name, true),
                                )
                              }
                            >
                              Start
                            </Button>
                            <Button
                              size="sm"
                              disabled={busy}
                              onClick={() =>
                                runAction(`Stop ${selectedService.name}`, () =>
                                  api.stopService(selectedService.name, true),
                                )
                              }
                            >
                              Stop
                            </Button>
                            <Button
                              size="sm"
                              disabled={busy}
                              onClick={() =>
                                runAction(`Restart ${selectedService.name}`, () =>
                                  api.restartService(selectedService.name),
                                )
                              }
                            >
                              Restart
                            </Button>
                            {selectedService.effective_build && (
                              <Button
                                size="sm"
                                disabled={busy}
                                onClick={() =>
                                  runAction(`Rebuild ${selectedService.name}`, () =>
                                    api.rebuildService(selectedService.name),
                                  )
                                }
                              >
                                Rebuild
                              </Button>
                            )}
                            <Button
                              size="sm"
                              disabled={busy}
                              onClick={() =>
                                runAction(`Register ${selectedService.name}`, () =>
                                  api.registerService(selectedService.name),
                                )
                              }
                            >
                              Register Agents
                            </Button>
                            <Button
                              size="sm"
                              disabled={busy}
                              onClick={() =>
                                runAction(`Deregister ${selectedService.name}`, () =>
                                  api.deregisterService(selectedService.name),
                                )
                              }
                            >
                              Deregister
                            </Button>
                            {selectedService.health_url && (
                              <Button
                                size="sm"
                                disabled={busy}
                                onClick={() => handleHealthCheck(selectedService.name)}
                              >
                                Test Health
                              </Button>
                            )}
                            {selectedService.dashboard_url && (
                              <Button
                                size="sm"
                                onClick={() =>
                                  window.open(selectedService.dashboard_url, '_blank', 'noopener,noreferrer')
                                }
                              >
                                Open Dashboard ↗
                              </Button>
                            )}
                          </div>

                          <Button
                            size="sm"
                            disabled={busy}
                            onClick={() => {
                              setDeleteTarget(selectedService.name);
                              setDeleteDeregister(true);
                              focusWindow('delete');
                            }}
                          >
                            Delete Service
                          </Button>
                        </div>

                        <GroupBox label={`Service Details: ${selectedService.name}`}>
                          <div
                            style={{
                              display: 'grid',
                              gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))',
                              gap: '6px 16px',
                              fontSize: 12,
                            }}
                          >
                            <div>
                              <b>Mode:</b> {selectedService.mode} | <b>Status:</b>{' '}
                              <StatusLed state={selectedService.state} />
                            </div>
                            <div>
                              <b>Detail:</b> {selectedService.detail}
                            </div>
                            <div>
                              <b>PID / Uptime:</b>{' '}
                              {selectedService.pid
                                ? `${selectedService.pid}${selectedService.uptime ? ` (${selectedService.uptime})` : ''}`
                                : '—'}
                            </div>
                            <div>
                              <b>Ports:</b>{' '}
                              {selectedService.port || '—'}
                              {selectedService.extra_ports?.length
                                ? ` (extra: ${selectedService.extra_ports.join(', ')})`
                                : ''}
                            </div>
                            <div style={{ gridColumn: '1 / -1', fontFamily: 'monospace', wordBreak: 'break-all' }}>
                              <b>Command:</b>{' '}
                              {selectedService.path
                                ? `embedded filesystem (${selectedService.path})`
                                : selectedService.command
                                  ? [selectedService.command, ...(selectedService.args || [])].join(' ')
                                  : '(remote-only)'}
                            </div>
                            <div style={{ fontFamily: 'monospace', wordBreak: 'break-all' }}>
                              <b>MCP Endpoint:</b> {selectedService.mcp_endpoint || '—'}
                            </div>
                            <div style={{ fontFamily: 'monospace', wordBreak: 'break-all' }}>
                              <b>Health URL:</b> {selectedService.health_url || '—'}
                            </div>
                            <div style={{ fontFamily: 'monospace', wordBreak: 'break-all' }}>
                              <b>Log Path:</b> {selectedService.resolved_log_path}
                            </div>
                            <div style={{ fontFamily: 'monospace', wordBreak: 'break-all' }}>
                              <b>Installer:</b> {selectedService.effective_installer || 'native (mcp-local)'}
                            </div>
                            <div>
                              <b>Active Tier:</b> {selectedService.active_tier || '(default)'} |{' '}
                              <b>Tool Overrides:</b> {selectedService.overrides_count}
                            </div>
                            <div>
                              <b>Registered:</b>{' '}
                              {Object.entries(selectedService.registered || {})
                                .map(([k, v]) => `${k}: ${v ? 'yes' : 'no'}`)
                                .join(' | ')}
                            </div>
                          </div>
                        </GroupBox>

                        {healthResult && (
                          <GroupBox
                            label={`Health Probe (${healthResult.ok ? 'Healthy' : 'Failed'} — ${healthResult.latency_ms}ms)`}
                          >
                            <div style={{ fontSize: 12, marginBottom: 4 }}>
                              <b>URL:</b> <code>{healthResult.url || '—'}</code>{' '}
                              {healthResult.status_code ? `| HTTP ${healthResult.status_code}` : ''}
                              {healthResult.error ? ` | Error: ${healthResult.error}` : ''}
                            </div>
                            {healthResult.body && (
                              <ScrollView style={{ maxHeight: 110, background: '#fff', color: '#000', padding: 6 }}>
                                <pre style={{ margin: 0, fontSize: 11 }}>{healthResult.body}</pre>
                              </ScrollView>
                            )}
                          </GroupBox>
                        )}

                        {actionLog && (
                          <GroupBox label="Last Action Output">
                            <ScrollView style={{ maxHeight: 110, background: '#000', color: '#c0ffc0', padding: 6 }}>
                              <pre style={{ margin: 0, fontSize: 11, whiteSpace: 'pre-wrap' }}>{actionLog}</pre>
                            </ScrollView>
                          </GroupBox>
                        )}
                      </div>
                    )}

                    {detailTab === 1 && editorInitial && (
                      <ServiceEditorForm
                        initial={editorInitial}
                        saving={busy}
                        onSubmit={async (req) => {
                          setBusy(true);
                          try {
                            const updated = await api.updateService(selectedService.name, req);
                            setSelectedName(updated.name);
                            notify(`Saved changes to ${updated.name}.`);
                            await refresh();
                          } catch (err) {
                            notify(err instanceof Error ? err.message : String(err), true);
                          } finally {
                            setBusy(false);
                          }
                        }}
                      />
                    )}

                    {detailTab === 2 && (
                      <ToolsPanel
                        service={selectedService}
                        onNotify={notify}
                        onAfterMutate={refresh}
                      />
                    )}

                    {detailTab === 3 && (
                      <LogsPanel service={selectedService} onNotify={notify} />
                    )}
                  </TabBody>
                </div>
              ) : (
                <Frame variant="well" style={{ padding: 16, textAlign: 'center', fontSize: 12 }}>
                  Select a service above to inspect controls, edit configuration, manage tool tiers, or view logs.
                </Frame>
              )}
            </div>
          </Win95Window>
        )}

        {/* Add Service Wizard Window */}
        {windows.add.open && (
          <Win95Window
            id="add"
            title="Add New MCP Service Wizard"
            icon={<ToolsGearIcon size={16} />}
            bounds={windows.add.bounds}
            zIndex={windows.add.z}
            active={activeWin === 'add'}
            minimized={windows.add.minimized}
            maximized={windows.add.maximized}
            onFocus={() => focusWindow('add')}
            onClose={() => closeWindow('add')}
            onMinimize={() => minimizeWindow('add')}
            onToggleMaximize={() => toggleMaximize('add')}
            onMove={(x, y) => moveWindow('add', x, y)}
          >
            <AddServiceModal
              onCancel={() => closeWindow('add')}
              onCreate={async (req, startAfter) => {
                try {
                  const created = await api.createService(req);
                  setSelectedName(created.name);
                  closeWindow('add');
                  if (startAfter) {
                    const startRes = await api.startService(created.name, true);
                    setActionLog(startRes.output);
                    notify(
                      startRes.ok
                        ? `Created, started, and registered ${created.name}.`
                        : `Created ${created.name}, start reported: ${startRes.error}`,
                      !startRes.ok,
                    );
                  } else {
                    notify(`Created service ${created.name}.`);
                  }
                  await refresh();
                } catch (err) {
                  notify(err instanceof Error ? err.message : String(err), true);
                }
              }}
            />
          </Win95Window>
        )}

        {/* Agent Sync Window */}
        {windows.agents.open && (
          <Win95Window
            id="agents"
            title="Agent Sync & IDE Configs"
            icon={<NetworkAgentsIcon size={16} />}
            bounds={windows.agents.bounds}
            zIndex={windows.agents.z}
            active={activeWin === 'agents'}
            minimized={windows.agents.minimized}
            maximized={windows.agents.maximized}
            onFocus={() => focusWindow('agents')}
            onClose={() => closeWindow('agents')}
            onMinimize={() => minimizeWindow('agents')}
            onToggleMaximize={() => toggleMaximize('agents')}
            onMove={(x, y) => moveWindow('agents', x, y)}
          >
            <AgentsWindow agents={state.agents} onNotify={notify} onAfterMutate={refresh} />
          </Win95Window>
        )}

        {/* Doctor & Validator Window */}
        {windows.doctor.open && (
          <Win95Window
            id="doctor"
            title="MCP Doctor & Config Validator"
            icon={<DoctorIcon size={16} />}
            bounds={windows.doctor.bounds}
            zIndex={windows.doctor.z}
            active={activeWin === 'doctor'}
            minimized={windows.doctor.minimized}
            maximized={windows.doctor.maximized}
            onFocus={() => focusWindow('doctor')}
            onClose={() => closeWindow('doctor')}
            onMinimize={() => minimizeWindow('doctor')}
            onToggleMaximize={() => toggleMaximize('doctor')}
            onMove={(x, y) => moveWindow('doctor', x, y)}
          >
            <DoctorWindow onNotify={notify} onAfterMutate={refresh} />
          </Win95Window>
        )}

        {/* Import Wizard Window */}
        {windows.import.open && (
          <Win95Window
            id="import"
            title="Import MCP Servers from Agent IDEs"
            icon={<FloppyImportIcon size={16} />}
            bounds={windows.import.bounds}
            zIndex={windows.import.z}
            active={activeWin === 'import'}
            minimized={windows.import.minimized}
            maximized={windows.import.maximized}
            onFocus={() => focusWindow('import')}
            onClose={() => closeWindow('import')}
            onMinimize={() => minimizeWindow('import')}
            onToggleMaximize={() => toggleMaximize('import')}
            onMove={(x, y) => moveWindow('import', x, y)}
          >
            <ImportWindow onNotify={notify} onAfterMutate={refresh} />
          </Win95Window>
        )}

        {/* Raw config.yaml Notepad Window */}
        {windows.raw.open && (
          <Win95Window
            id="raw"
            title="Notepad — config.yaml"
            icon={<NotepadYamlIcon size={16} />}
            bounds={windows.raw.bounds}
            zIndex={windows.raw.z}
            active={activeWin === 'raw'}
            minimized={windows.raw.minimized}
            maximized={windows.raw.maximized}
            onFocus={() => focusWindow('raw')}
            onClose={() => closeWindow('raw')}
            onMinimize={() => minimizeWindow('raw')}
            onToggleMaximize={() => toggleMaximize('raw')}
            onMove={(x, y) => moveWindow('raw', x, y)}
          >
            <RawConfigWindow onNotify={notify} onAfterMutate={refresh} />
          </Win95Window>
        )}

        {/* Confirm Delete Dialog */}
        {windows.delete.open && deleteTarget && (
          <Win95Window
            id="delete"
            title={`Confirm Delete — ${deleteTarget}`}
            bounds={windows.delete.bounds}
            zIndex={windows.delete.z}
            active={activeWin === 'delete'}
            onFocus={() => focusWindow('delete')}
            onClose={() => closeWindow('delete')}
            onMove={(x, y) => moveWindow('delete', x, y)}
          >
            <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
              <div style={{ fontSize: 13 }}>
                Remove <b>{deleteTarget}</b> from <code>config.yaml</code>?
              </div>
              <Checkbox
                checked={deleteDeregister}
                onChange={(e) => setDeleteDeregister(e.target.checked)}
                label="Also deregister from OpenCode, Cursor, and Claude Desktop"
              />
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
                <Button onClick={() => closeWindow('delete')}>Cancel</Button>
                <Button
                  primary
                  disabled={busy}
                  onClick={async () => {
                    const name = deleteTarget;
                    closeWindow('delete');
                    await runAction(`Delete ${name}`, () => api.deleteService(name, deleteDeregister));
                  }}
                >
                  Delete
                </Button>
              </div>
            </div>
          </Win95Window>
        )}
      </div>

      {/* Start Menu */}
      {startOpen && (
        <div
          style={{
            position: 'fixed',
            bottom: 46,
            left: 4,
            zIndex: 9999,
            display: 'flex',
            alignItems: 'flex-end',
          }}
        >
          <MenuList style={{ minWidth: 215 }}>
            <MenuListItem onClick={() => { focusWindow('manager'); setStartOpen(false); }}>
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <ComputerIcon size={16} /> MCP Control Panel
              </span>
            </MenuListItem>
            <MenuListItem onClick={() => { focusWindow('add'); setStartOpen(false); }}>
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <ToolsGearIcon size={16} /> New MCP Service...
              </span>
            </MenuListItem>
            <Separator />
            <MenuListItem onClick={() => { focusWindow('agents'); setStartOpen(false); }}>
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <NetworkAgentsIcon size={16} /> Agent Sync & IDEs
              </span>
            </MenuListItem>
            <MenuListItem onClick={() => { focusWindow('doctor'); setStartOpen(false); }}>
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <DoctorIcon size={16} /> Doctor & Validator
              </span>
            </MenuListItem>
            <MenuListItem onClick={() => { focusWindow('import'); setStartOpen(false); }}>
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <FloppyImportIcon size={16} /> Import from Agents...
              </span>
            </MenuListItem>
            <MenuListItem onClick={() => { focusWindow('raw'); setStartOpen(false); }}>
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <NotepadYamlIcon size={16} /> Edit config.yaml
              </span>
            </MenuListItem>
            <Separator />
            <MenuListItem onClick={() => setThemesOpen(!themesOpen)}>
              <span>Theme: {activeTheme.label} ▸</span>
            </MenuListItem>
            <MenuListItem
              onClick={() => {
                setStartOpen(false);
                void refresh();
              }}
            >
              <span>Refresh All</span>
            </MenuListItem>
          </MenuList>

          {themesOpen && (
            <MenuList style={{ marginLeft: 2, minWidth: 185 }}>
              {(Object.keys(THEMES) as ThemeKey[]).map((key) => (
                <MenuListItem
                  key={key}
                  onClick={() => {
                    setThemeKey(key);
                    localStorage.setItem('mcp-local-95-theme', key);
                    setThemesOpen(false);
                    setStartOpen(false);
                  }}
                >
                  <span>
                    {key === themeKey ? '✓ ' : '   '}
                    {THEMES[key].label}
                  </span>
                </MenuListItem>
              ))}
            </MenuList>
          )}
        </div>
      )}

      {/* Win95 Bottom Taskbar */}
      <AppBar style={{ top: 'auto', bottom: 0, height: 48, zIndex: 9998 }}>
        <Toolbar style={{ justifyContent: 'space-between', height: '100%', padding: '0 6px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 6, flex: 1, minWidth: 0, overflow: 'hidden' }}>
            <Button
              onClick={(e) => {
                e.stopPropagation();
                setStartOpen(!startOpen);
                setThemesOpen(false);
              }}
              active={startOpen}
              style={{ fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: 6, flexShrink: 0 }}
            >
              <WinLogoIcon size={16} />
              Start
            </Button>

            <Separator orientation="vertical" size={28} />

            <div style={{ display: 'flex', gap: 4, overflow: 'hidden' }}>
              {(Object.keys(windows) as WindowId[])
                .filter((id) => windows[id].open)
                .map((id) => {
                  const isFocused = activeWin === id && !windows[id].minimized;
                  return (
                    <Button
                      key={id}
                      active={isFocused}
                      onClick={() => {
                        if (isFocused) {
                          minimizeWindow(id);
                        } else {
                          focusWindow(id);
                        }
                      }}
                      style={{
                        maxWidth: 175,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                        fontWeight: isFocused ? 'bold' : 'normal',
                      }}
                    >
                      {windowTitles[id]}
                    </Button>
                  );
                })}
            </div>
          </div>

          {/* System Tray */}
          <Frame
            variant="well"
            style={{
              padding: '4px 10px',
              display: 'flex',
              alignItems: 'center',
              gap: 10,
              fontSize: 12,
              flexShrink: 0,
            }}
          >
            <span title="Running services / total configured">
              {runningCount}/{state.services.length} UP
            </span>
            <span>{clock}</span>
          </Frame>
        </Toolbar>
      </AppBar>
    </ThemeProvider>
  );
}
