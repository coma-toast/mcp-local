import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Button, GroupBox, TextInput } from 'react95';
import { api } from '../api';
import { PRESET_CATEGORIES, PRESET_SOURCES, SERVICE_PRESETS } from '../presets';
import type { ServicePreset } from '../types';

interface PresetDropdownProps {
  selectedPreset: ServicePreset | null;
  onSelect: (preset: ServicePreset) => void;
}

function sourceBadgeStyle(source: string): React.CSSProperties {
  switch (source) {
    case 'zapier.com/mcp':
      return { background: '#ff4f00', color: '#ffffff' };
    case 'github.com/mcp':
      return { background: '#000080', color: '#ffffff' };
    case 'mcpservers.org':
      return { background: '#006644', color: '#ffffff' };
    case 'template':
      return { background: '#555555', color: '#ffffff' };
    default:
      return { background: '#4b0082', color: '#ffffff' };
  }
}

function typeBadgeStyle(type: string): React.CSSProperties {
  switch (type) {
    case 'http':
      return { background: '#e0f0ff', color: '#003366', border: '1px solid #003366' };
    case 'filesystem':
      return { background: '#fff4d6', color: '#664400', border: '1px solid #664400' };
    default:
      return { background: '#e8ffe8', color: '#005500', border: '1px solid #005500' };
  }
}

function commandPreview(preset: ServicePreset): string {
  const s = preset.service;
  if (s.type === 'filesystem') {
    return `filesystem: ${s.path || '~/projects'} (port ${s.port || 4000})`;
  }
  if (s.type === 'http' && !s.command) {
    return s.mcp_url || 'https://...';
  }
  const parts = [s.command, ...(s.args || [])].filter(Boolean);
  return parts.join(' ') || s.mcp_url || '';
}

export const PresetDropdown: React.FC<PresetDropdownProps> = ({ selectedPreset, onSelect }) => {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState('');
  const [sourceFilter, setSourceFilter] = useState<string>('all');
  const [categoryFilter, setCategoryFilter] = useState<string>('All Categories');
  const [activeIndex, setActiveIndex] = useState(0);

  const [livePresets, setLivePresets] = useState<ServicePreset[]>([]);
  const [liveLoading, setLiveLoading] = useState(false);

  const containerRef = useRef<HTMLDivElement | null>(null);
  const inputRef = useRef<HTMLInputElement | null>(null);
  const listRef = useRef<HTMLDivElement | null>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const onDocMouseDown = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener('mousedown', onDocMouseDown);
    return () => document.removeEventListener('mousedown', onDocMouseDown);
  }, []);

  // Debounced live MCP Registry search when user types >= 2 chars
  useEffect(() => {
    const trimmed = query.trim();
    if (trimmed.length < 2) {
      setLivePresets([]);
      setLiveLoading(false);
      return;
    }
    let cancelled = false;
    setLiveLoading(true);
    const timer = window.setTimeout(async () => {
      try {
        const res = await api.searchRegistry(trimmed);
        if (!cancelled) {
          setLivePresets(res.presets || []);
        }
      } catch {
        if (!cancelled) {
          setLivePresets([]);
        }
      } finally {
        if (!cancelled) {
          setLiveLoading(false);
        }
      }
    }, 260);
    return () => {
      cancelled = true;
      window.clearTimeout(timer);
    };
  }, [query]);

  // Filter curated + live presets
  const filteredPresets = useMemo(() => {
    const q = query.trim().toLowerCase();
    const tokens = q ? q.split(/\s+/).filter(Boolean) : [];

    const matchesPreset = (p: ServicePreset): boolean => {
      if (sourceFilter === 'popular' && !p.popular) return false;
      if (
        sourceFilter !== 'all' &&
        sourceFilter !== 'popular' &&
        p.source !== sourceFilter
      ) {
        return false;
      }
      if (categoryFilter !== 'All Categories' && p.category !== categoryFilter) {
        return false;
      }
      if (tokens.length === 0) return true;
      const haystack = [
        p.title,
        p.service.name,
        p.category,
        p.source,
        p.description,
        p.service.command,
        ...(p.service.args || []),
        p.service.mcp_url,
        ...(p.tags || []),
      ]
        .filter(Boolean)
        .join(' ')
        .toLowerCase();
      return tokens.every((tok) => haystack.includes(tok));
    };

    const curated = SERVICE_PRESETS.filter(matchesPreset);
    const seenNames = new Set(curated.map((c) => c.service.name.toLowerCase()));

    const extraLive: ServicePreset[] = [];
    if (
      tokens.length > 0 &&
      (sourceFilter === 'all' || sourceFilter === 'github.com/mcp') &&
      categoryFilter === 'All Categories'
    ) {
      for (const lp of livePresets) {
        const key = lp.service.name.toLowerCase();
        if (!seenNames.has(key)) {
          seenNames.add(key);
          extraLive.push({
            ...lp,
            source: 'MCP Registry Live',
          });
        }
      }
    }

    return [...curated, ...extraLive];
  }, [query, sourceFilter, categoryFilter, livePresets]);

  useEffect(() => {
    setActiveIndex(0);
  }, [query, sourceFilter, categoryFilter]);

  const handleSelect = (preset: ServicePreset) => {
    // Deep-clone service so user edits don't mutate preset definition
    const cloned: ServicePreset = {
      ...preset,
      service: {
        ...preset.service,
        args: [...(preset.service.args || [])],
        env: { ...(preset.service.env || {}) },
        headers: { ...(preset.service.headers || {}) },
        deps: [...(preset.service.deps || [])],
        extra_ports: [...(preset.service.extra_ports || [])],
      },
    };
    onSelect(cloned);
    setOpen(false);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!open && (e.key === 'ArrowDown' || e.key === 'Enter')) {
      setOpen(true);
      return;
    }
    if (!open) return;

    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setActiveIndex((prev) => Math.min(prev + 1, Math.max(0, filteredPresets.length - 1)));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setActiveIndex((prev) => Math.max(0, prev - 1));
    } else if (e.key === 'Enter') {
      if (filteredPresets[activeIndex]) {
        e.preventDefault();
        handleSelect(filteredPresets[activeIndex]);
      }
    } else if (e.key === 'Escape') {
      e.preventDefault();
      setOpen(false);
    }
  };

  return (
    <GroupBox label="Popular Service Presets (github.com/mcp • mcpservers.org • zapier.com/mcp)">
      <div ref={containerRef} style={{ position: 'relative', display: 'flex', flexDirection: 'column', gap: 8 }}>
        {/* Top bar: Dropdown toggle button + Find-as-you-type search box */}
        <div style={{ display: 'flex', gap: 6, alignItems: 'center', flexWrap: 'wrap' }}>
          <Button
            type="button"
            onClick={() => {
              const next = !open;
              setOpen(next);
              if (next) {
                window.setTimeout(() => inputRef.current?.focus(), 20);
              }
            }}
            active={open}
            style={{ fontWeight: 'bold', minWidth: 235, justifyContent: 'space-between' }}
          >
            <span>
              📋 {selectedPreset ? `Preset: ${selectedPreset.service.name}` : 'Browse Popular Presets'}
            </span>
            <span style={{ marginLeft: 8 }}>{open ? '▴' : '▾'}</span>
          </Button>

          <div style={{ flex: 1, minWidth: 240, display: 'flex', gap: 4, alignItems: 'center' }}>
            <TextInput
              ref={inputRef}
              value={query}
              placeholder="Find as you type: zapier, playwright, github, postgres, context7, grafana, notion..."
              onFocus={() => setOpen(true)}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                setQuery(e.target.value);
                if (!open) setOpen(true);
              }}
              onKeyDown={handleKeyDown}
              fullWidth
            />
            {query && (
              <Button
                size="sm"
                type="button"
                onClick={() => {
                  setQuery('');
                  inputRef.current?.focus();
                }}
                title="Clear search"
              >
                ✕
              </Button>
            )}
          </div>
        </div>

        {/* Quick-access popular chips when dropdown is closed */}
        <div style={{ display: 'flex', gap: 4, flexWrap: 'wrap', alignItems: 'center' }}>
          <span style={{ fontSize: 11, color: '#333', marginRight: 2 }}>Quick pick:</span>
          {[
            'zapier-mcp-all',
            'github-stdio',
            'playwright-mcp',
            'context7-stdio',
            'fetch-mcp',
            'postgres-mcp',
            'brave-search',
            'sequential-thinking',
            'filesystem-embedded',
          ].map((id) => {
            const p = SERVICE_PRESETS.find((x) => x.id === id);
            if (!p) return null;
            const isSelected = selectedPreset?.id === p.id;
            return (
              <Button
                key={p.id}
                size="sm"
                type="button"
                active={isSelected}
                onClick={() => handleSelect(p)}
                style={{ fontSize: 11, height: 24, padding: '0 6px' }}
                title={p.description}
              >
                {p.service.name}
              </Button>
            );
          })}
          <Button
            size="sm"
            type="button"
            onClick={() => {
              setOpen(true);
              window.setTimeout(() => inputRef.current?.focus(), 20);
            }}
            style={{ fontSize: 11, height: 24, padding: '0 8px', fontWeight: 'bold' }}
          >
            All ({SERVICE_PRESETS.length}+) ▾
          </Button>
        </div>

        {/* Floating Find-As-You-Type Dropdown Menu */}
        {open && (
          <div
            style={{
              position: 'absolute',
              top: 40,
              left: 0,
              right: 0,
              zIndex: 100,
              background: '#c0c0c0',
              borderTop: '2px solid #ffffff',
              borderLeft: '2px solid #ffffff',
              borderRight: '2px solid #000000',
              borderBottom: '2px solid #000000',
              boxShadow: '4px 4px 0px rgba(0,0,0,0.45)',
              padding: 8,
              display: 'flex',
              flexDirection: 'column',
              gap: 6,
            }}
          >
            {/* Source Filter Tabs */}
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 6, flexWrap: 'wrap' }}>
              <div style={{ display: 'flex', gap: 4, flexWrap: 'wrap' }}>
                {PRESET_SOURCES.map((src) => (
                  <Button
                    key={src.id}
                    size="sm"
                    type="button"
                    active={sourceFilter === src.id}
                    onClick={() => setSourceFilter(src.id)}
                    style={{ fontSize: 11, height: 24, padding: '0 7px' }}
                  >
                    {src.label}
                  </Button>
                ))}
              </div>
              <div style={{ fontSize: 11, color: '#222' }}>
                {liveLoading
                  ? 'Searching live MCP Registry...'
                  : `${filteredPresets.length} matching service${filteredPresets.length === 1 ? '' : 's'}`}
              </div>
            </div>

            {/* Category Filter Pills */}
            <div style={{ display: 'flex', gap: 4, flexWrap: 'wrap', paddingBottom: 2, borderBottom: '1px solid #808080' }}>
              {PRESET_CATEGORIES.map((cat) => (
                <button
                  key={cat}
                  type="button"
                  onClick={() => setCategoryFilter(cat)}
                  style={{
                    fontSize: 11,
                    padding: '2px 7px',
                    cursor: 'pointer',
                    border: categoryFilter === cat ? '1px solid #000080' : '1px solid #808080',
                    background: categoryFilter === cat ? '#000080' : '#dfdfdf',
                    color: categoryFilter === cat ? '#ffffff' : '#000000',
                    fontWeight: categoryFilter === cat ? 'bold' : 'normal',
                  }}
                >
                  {cat}
                </button>
              ))}
            </div>

            {/* Sunken Scrollable Results List */}
            <div
              ref={listRef}
              style={{
                maxHeight: 310,
                overflowY: 'auto',
                background: '#ffffff',
                borderTop: '2px solid #808080',
                borderLeft: '2px solid #808080',
                borderRight: '2px solid #dfdfdf',
                borderBottom: '2px solid #dfdfdf',
              }}
            >
              {filteredPresets.length === 0 ? (
                <div style={{ padding: 16, textAlign: 'center', fontSize: 12, color: '#444' }}>
                  {liveLoading
                    ? 'Searching MCP Registry...'
                    : `No presets matched "${query}". Try clearing filters or searching another keyword.`}
                </div>
              ) : (
                filteredPresets.map((preset, idx) => {
                  const isHighlighted = idx === activeIndex;
                  const cmdText = commandPreview(preset);
                  return (
                    <div
                      key={preset.id}
                      onMouseEnter={() => setActiveIndex(idx)}
                      onClick={() => handleSelect(preset)}
                      style={{
                        padding: '7px 10px',
                        cursor: 'pointer',
                        borderBottom: '1px solid #e4e4e4',
                        background: isHighlighted ? '#000080' : idx % 2 === 0 ? '#ffffff' : '#f7f7f7',
                        color: isHighlighted ? '#ffffff' : '#000000',
                        display: 'flex',
                        flexDirection: 'column',
                        gap: 3,
                      }}
                    >
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 8 }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: 6, flexWrap: 'wrap' }}>
                          <span style={{ fontWeight: 'bold', fontSize: 13 }}>
                            {preset.popular ? '★ ' : ''}
                            {preset.title}
                          </span>
                          <span
                            style={{
                              fontSize: 10,
                              padding: '1px 5px',
                              fontFamily: 'monospace',
                              ...(isHighlighted
                                ? { background: '#ffffff', color: '#000080' }
                                : typeBadgeStyle(preset.service.type || 'stdio')),
                            }}
                          >
                            {(preset.service.type || 'stdio').toUpperCase()}
                          </span>
                          <span
                            style={{
                              fontSize: 10,
                              padding: '1px 5px',
                              ...sourceBadgeStyle(preset.source),
                            }}
                          >
                            {preset.source}
                          </span>
                        </div>
                        <span
                          style={{
                            fontSize: 11,
                            opacity: 0.85,
                            whiteSpace: 'nowrap',
                          }}
                        >
                          {preset.category}
                        </span>
                      </div>

                      <div
                        style={{
                          fontSize: 11,
                          color: isHighlighted ? '#e8e8ff' : '#333333',
                          lineHeight: 1.3,
                        }}
                      >
                        {preset.description}
                      </div>

                      {cmdText && (
                        <div
                          style={{
                            fontSize: 11,
                            fontFamily: 'monospace',
                            color: isHighlighted ? '#ffff99' : '#004080',
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                          }}
                        >
                          ▸ {cmdText}
                        </div>
                      )}
                    </div>
                  );
                })
              )}
            </div>

            {/* Footer bar with keyboard hints & Close button */}
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', fontSize: 11 }}>
              <span style={{ color: '#333' }}>
                Tip: Use <b>↑/↓</b> + <b>Enter</b> to select, or type any package name to search the live MCP Registry.
              </span>
              <Button size="sm" type="button" onClick={() => setOpen(false)}>
                Close ▴
              </Button>
            </div>
          </div>
        )}

        {/* Selected Preset Setup Banner */}
        {selectedPreset && (
          <div
            style={{
              background: '#ffffe1',
              border: '1px solid #808000',
              padding: '6px 10px',
              fontSize: 12,
              display: 'flex',
              flexDirection: 'column',
              gap: 3,
            }}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 8 }}>
              <div>
                <b>Loaded Preset:</b> {selectedPreset.title}{' '}
                <span
                  style={{
                    fontSize: 10,
                    padding: '1px 5px',
                    marginLeft: 4,
                    ...sourceBadgeStyle(selectedPreset.source),
                  }}
                >
                  {selectedPreset.source}
                </span>
              </div>
              {selectedPreset.source_url && (
                <a
                  href={selectedPreset.source_url}
                  target="_blank"
                  rel="noreferrer"
                  style={{ color: '#000080', fontWeight: 'bold', fontSize: 11 }}
                >
                  Open Docs / Repo ↗
                </a>
              )}
            </div>
            {selectedPreset.setup_hint && (
              <div style={{ color: '#443000', fontSize: 11 }}>
                💡 <b>Setup Note:</b> {selectedPreset.setup_hint}
              </div>
            )}
          </div>
        )}
      </div>
    </GroupBox>
  );
};
