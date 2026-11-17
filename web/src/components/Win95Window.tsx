import React, { useRef, useState } from 'react';
import { Button, Window, WindowContent, WindowHeader } from 'react95';

export interface WinBounds {
  x: number;
  y: number;
  width: number;
  height?: number;
}

interface Win95WindowProps {
  id: string;
  title: string;
  icon?: React.ReactNode;
  bounds: WinBounds;
  zIndex: number;
  active: boolean;
  minimized?: boolean;
  maximized?: boolean;
  onFocus: () => void;
  onClose: () => void;
  onMinimize?: () => void;
  onToggleMaximize?: () => void;
  onMove: (x: number, y: number) => void;
  children: React.ReactNode;
  footer?: React.ReactNode;
}

export const Win95Window: React.FC<Win95WindowProps> = ({
  title,
  icon,
  bounds,
  zIndex,
  active,
  minimized,
  maximized,
  onFocus,
  onClose,
  onMinimize,
  onToggleMaximize,
  onMove,
  children,
  footer,
}) => {
  const dragRef = useRef<{ startX: number; startY: number; origX: number; origY: number } | null>(null);
  const [dragging, setDragging] = useState(false);

  if (minimized) {
    return null;
  }

  const handleMouseDownHeader = (e: React.MouseEvent) => {
    if ((e.target as HTMLElement).closest('button')) return;
    onFocus();
    if (maximized) return;
    dragRef.current = {
      startX: e.clientX,
      startY: e.clientY,
      origX: bounds.x,
      origY: bounds.y,
    };
    setDragging(true);

    const onMouseMove = (ev: MouseEvent) => {
      if (!dragRef.current) return;
      const nx = Math.max(0, dragRef.current.origX + (ev.clientX - dragRef.current.startX));
      const ny = Math.max(0, dragRef.current.origY + (ev.clientY - dragRef.current.startY));
      onMove(nx, ny);
    };
    const onMouseUp = () => {
      dragRef.current = null;
      setDragging(false);
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
    };
    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
  };

  const style: React.CSSProperties = maximized
    ? {
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 48,
        width: '100%',
        height: 'calc(100vh - 48px)',
        zIndex,
        display: 'flex',
        flexDirection: 'column',
      }
    : {
        position: 'fixed',
        top: bounds.y,
        left: bounds.x,
        width: `min(${bounds.width}px, calc(100vw - 16px))`,
        height: bounds.height ? `min(${bounds.height}px, calc(100vh - 64px))` : undefined,
        maxHeight: 'calc(100vh - 60px)',
        zIndex,
        display: 'flex',
        flexDirection: 'column',
      };

  return (
    <Window style={style} onMouseDown={onFocus}>
      <WindowHeader
        active={active}
        onMouseDown={handleMouseDownHeader}
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          cursor: maximized ? 'default' : dragging ? 'grabbing' : 'grab',
          userSelect: 'none',
          gap: 8,
          flexShrink: 0,
        }}
      >
        <span style={{ display: 'flex', alignItems: 'center', gap: 6, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
          {icon}
          <span style={{ overflow: 'hidden', textOverflow: 'ellipsis' }}>{title}</span>
        </span>
        <span style={{ display: 'flex', alignItems: 'center', gap: 2, flexShrink: 0 }}>
          {onMinimize && (
            <Button size="sm" square onClick={onMinimize} title="Minimize" style={{ minWidth: 22, padding: 0, fontWeight: 'bold' }}>
              _
            </Button>
          )}
          {onToggleMaximize && (
            <Button size="sm" square onClick={onToggleMaximize} title={maximized ? 'Restore' : 'Maximize'} style={{ minWidth: 22, padding: 0, fontWeight: 'bold' }}>
              {maximized ? '❐' : '□'}
            </Button>
          )}
          <Button size="sm" square onClick={onClose} title="Close" style={{ minWidth: 22, padding: 0, fontWeight: 'bold' }}>
            ×
          </Button>
        </span>
      </WindowHeader>
      <WindowContent
        style={{
          flex: 1,
          minHeight: 0,
          display: 'flex',
          flexDirection: 'column',
          overflow: 'auto',
          padding: 10,
        }}
      >
        {children}
      </WindowContent>
      {footer}
    </Window>
  );
};
