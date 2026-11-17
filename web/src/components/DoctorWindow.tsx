import React, { useCallback, useEffect, useState } from 'react';
import {
  Button,
  Frame,
  GroupBox,
  ScrollView,
  Table,
  TableBody,
  TableDataCell,
  TableHead,
  TableHeadCell,
  TableRow,
} from 'react95';
import { api } from '../api';
import { StatusLed } from '../icons';
import type { DoctorResponse, ValidateResponse } from '../types';

interface DoctorWindowProps {
  onNotify: (msg: string, isError?: boolean) => void;
  onAfterMutate: () => Promise<void>;
}

export const DoctorWindow: React.FC<DoctorWindowProps> = ({ onNotify, onAfterMutate }) => {
  const [validation, setValidation] = useState<ValidateResponse | null>(null);
  const [doctor, setDoctor] = useState<DoctorResponse | null>(null);
  const [busy, setBusy] = useState(false);

  const runDiagnostics = useCallback(
    async (fix = false) => {
      setBusy(true);
      try {
        const [valRes, docRes] = await Promise.all([api.validate(), api.doctor('', fix)]);
        setValidation(valRes);
        setDoctor(docRes);
        if (fix) {
          onNotify(
            docRes.remaining === 0
              ? 'Doctor --fix resolved all fixable issues.'
              : `Doctor --fix finished with ${docRes.remaining} remaining issue(s).`,
            docRes.remaining > 0,
          );
          await onAfterMutate();
        }
      } catch (err) {
        onNotify(err instanceof Error ? err.message : String(err), true);
      } finally {
        setBusy(false);
      }
    },
    [onNotify, onAfterMutate],
  );

  useEffect(() => {
    void runDiagnostics(false);
  }, [runDiagnostics]);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
      <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', justifyContent: 'space-between', alignItems: 'center' }}>
        <div style={{ fontSize: 12 }}>
          Checks config paths, binaries, dependencies, stale PIDs, orphan processes, port conflicts, and agent registrations.
        </div>
        <div style={{ display: 'flex', gap: 6 }}>
          <Button size="sm" onClick={() => runDiagnostics(false)} disabled={busy}>
            Run Checks
          </Button>
          <Button size="sm" primary onClick={() => runDiagnostics(true)} disabled={busy}>
            {busy ? 'Running...' : 'Doctor --fix'}
          </Button>
        </div>
      </div>

      {validation && (
        <GroupBox label={`Config Validation (${validation.valid ? 'Valid' : 'Errors Found'})`}>
          {validation.issues.length === 0 ? (
            <div style={{ fontSize: 12, color: '#006400', fontWeight: 'bold' }}>
              No configuration errors or warnings.
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              {validation.issues.map((iss, idx) => (
                <Frame
                  key={idx}
                  variant="well"
                  style={{
                    padding: '4px 8px',
                    fontSize: 12,
                    background: iss.level === 'error' ? '#ffe0e0' : '#ffffd0',
                    color: '#000',
                  }}
                >
                  <b>[{iss.level.toUpperCase()}] {iss.service}:</b> {iss.message}
                </Frame>
              ))}
            </div>
          )}
        </GroupBox>
      )}

      {doctor && (
        <GroupBox label={`Doctor Findings (${doctor.remaining} remaining)`}>
          <div style={{ maxHeight: 220, overflow: 'auto', border: '2px inset #808080', background: '#fff' }}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableHeadCell style={{ width: 150 }}>Service</TableHeadCell>
                  <TableHeadCell style={{ width: 120 }}>State</TableHeadCell>
                  <TableHeadCell>Findings</TableHeadCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {doctor.services.map((svc) => (
                  <TableRow key={svc.name}>
                    <TableDataCell style={{ fontWeight: 'bold' }}>{svc.name}</TableDataCell>
                    <TableDataCell>
                      <StatusLed state={svc.state} />
                    </TableDataCell>
                    <TableDataCell>
                      {svc.findings.length === 0 ? (
                        <span style={{ color: '#006400' }}>OK — no problems</span>
                      ) : (
                        <ul style={{ margin: 0, paddingLeft: 16 }}>
                          {svc.findings.map((f, i) => (
                            <li key={i}>
                              {f.fixed ? '✅ Fixed: ' : f.fixable ? '🔧 ' : '⚠️ '}
                              {f.message}
                              {f.error ? ` (error: ${f.error})` : ''}
                            </li>
                          ))}
                        </ul>
                      )}
                    </TableDataCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {doctor.output && (
            <ScrollView style={{ maxHeight: 130, marginTop: 8, background: '#000', color: '#c0ffc0', padding: 8 }}>
              <pre style={{ margin: 0, fontFamily: 'monospace', fontSize: 11, whiteSpace: 'pre-wrap' }}>
                {doctor.output}
              </pre>
            </ScrollView>
          )}
        </GroupBox>
      )}
    </div>
  );
};
