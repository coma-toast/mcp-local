import React from 'react';

interface IconProps {
  size?: number;
}

export const ComputerIcon: React.FC<IconProps> = ({ size = 32 }) => (
  <svg width={size} height={size} viewBox="0 0 16 16" shapeRendering="crispEdges">
    <rect x="2" y="1" width="12" height="9" fill="#c0c0c0" stroke="#000" strokeWidth="1" />
    <rect x="3" y="2" width="10" height="7" fill="#008080" />
    <rect x="4" y="3" width="3" height="1" fill="#00ffff" />
    <rect x="4" y="5" width="6" height="1" fill="#ffff00" />
    <rect x="6" y="10" width="4" height="2" fill="#808080" />
    <rect x="1" y="12" width="14" height="3" fill="#c0c0c0" stroke="#000" strokeWidth="1" />
    <rect x="10" y="13" width="3" height="1" fill="#008000" />
  </svg>
);

export const NetworkAgentsIcon: React.FC<IconProps> = ({ size = 32 }) => (
  <svg width={size} height={size} viewBox="0 0 16 16" shapeRendering="crispEdges">
    <rect x="1" y="2" width="6" height="5" fill="#c0c0c0" stroke="#000" strokeWidth="1" />
    <rect x="2" y="3" width="4" height="3" fill="#000080" />
    <rect x="9" y="2" width="6" height="5" fill="#c0c0c0" stroke="#000" strokeWidth="1" />
    <rect x="10" y="3" width="4" height="3" fill="#008080" />
    <rect x="3" y="7" width="2" height="3" fill="#808080" />
    <rect x="11" y="7" width="2" height="3" fill="#808080" />
    <rect x="2" y="10" width="12" height="2" fill="#008000" />
    <rect x="7" y="12" width="2" height="2" fill="#008000" />
    <rect x="4" y="14" width="8" height="1" fill="#000" />
  </svg>
);

export const DoctorIcon: React.FC<IconProps> = ({ size = 32 }) => (
  <svg width={size} height={size} viewBox="0 0 16 16" shapeRendering="crispEdges">
    <rect x="2" y="4" width="12" height="10" fill="#ffffff" stroke="#000" strokeWidth="1" />
    <rect x="5" y="2" width="6" height="2" fill="#808080" stroke="#000" strokeWidth="1" />
    <rect x="7" y="6" width="2" height="6" fill="#ff0000" />
    <rect x="5" y="8" width="6" height="2" fill="#ff0000" />
  </svg>
);

export const FloppyImportIcon: React.FC<IconProps> = ({ size = 32 }) => (
  <svg width={size} height={size} viewBox="0 0 16 16" shapeRendering="crispEdges">
    <rect x="1" y="1" width="14" height="14" fill="#000080" stroke="#000" strokeWidth="1" />
    <rect x="3" y="1" width="9" height="5" fill="#c0c0c0" />
    <rect x="8" y="2" width="3" height="3" fill="#000080" />
    <rect x="3" y="8" width="10" height="6" fill="#ffffff" />
    <rect x="4" y="10" width="8" height="1" fill="#808080" />
    <rect x="4" y="12" width="6" height="1" fill="#808080" />
  </svg>
);

export const NotepadYamlIcon: React.FC<IconProps> = ({ size = 32 }) => (
  <svg width={size} height={size} viewBox="0 0 16 16" shapeRendering="crispEdges">
    <rect x="2" y="1" width="12" height="14" fill="#ffffff" stroke="#000" strokeWidth="1" />
    <rect x="2" y="1" width="12" height="3" fill="#008080" />
    <rect x="4" y="6" width="8" height="1" fill="#000" />
    <rect x="4" y="8" width="6" height="1" fill="#000080" />
    <rect x="4" y="10" width="8" height="1" fill="#000" />
    <rect x="4" y="12" width="5" height="1" fill="#808080" />
  </svg>
);

export const ToolsGearIcon: React.FC<IconProps> = ({ size = 16 }) => (
  <svg width={size} height={size} viewBox="0 0 16 16" shapeRendering="crispEdges">
    <rect x="6" y="1" width="4" height="14" fill="#808080" />
    <rect x="1" y="6" width="14" height="4" fill="#808080" />
    <rect x="3" y="3" width="10" height="10" fill="#c0c0c0" stroke="#000" strokeWidth="1" />
    <rect x="6" y="6" width="4" height="4" fill="#008080" stroke="#000" strokeWidth="1" />
  </svg>
);

export const WinLogoIcon: React.FC<IconProps> = ({ size = 16 }) => (
  <svg width={size} height={size} viewBox="0 0 16 16" shapeRendering="crispEdges">
    <rect x="2" y="2" width="5" height="5" fill="#ff0000" />
    <rect x="9" y="2" width="5" height="5" fill="#00a800" />
    <rect x="2" y="9" width="5" height="5" fill="#0000ff" />
    <rect x="9" y="9" width="5" height="5" fill="#ffff00" />
  </svg>
);

export const StatusLed: React.FC<{ state: string }> = ({ state }) => {
  let color = '#808080';
  let label = state;
  switch (state) {
    case 'running':
      color = '#00b000';
      break;
    case 'degraded':
    case 'orphaned':
      color = '#d8a000';
      break;
    case 'port conflict':
    case 'stale':
      color = '#e00000';
      break;
    case 'stopped':
      color = '#808080';
      break;
    default:
      color = '#808080';
  }
  return (
    <span style={{ display: 'inline-flex', alignItems: 'center', gap: 6, whiteSpace: 'nowrap' }}>
      <span
        aria-hidden="true"
        style={{
          width: 9,
          height: 9,
          display: 'inline-block',
          background: color,
          border: '1px solid #000',
          boxShadow: 'inset 1px 1px 0 rgba(255,255,255,0.55)',
          flexShrink: 0,
        }}
      />
      <span>{label}</span>
    </span>
  );
};
