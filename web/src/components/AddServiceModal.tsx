import React, { useState } from 'react';
import { Checkbox } from 'react95';
import type { ServicePreset, ServiceUpsertRequest } from '../types';
import { PresetDropdown } from './PresetDropdown';
import { ServiceEditorForm, emptyUpsertRequest } from './ServiceEditorForm';

interface AddServiceModalProps {
  onCreate: (req: ServiceUpsertRequest, startAfter: boolean) => Promise<void>;
  onCancel: () => void;
}

export const AddServiceModal: React.FC<AddServiceModalProps> = ({ onCreate, onCancel }) => {
  const [template, setTemplate] = useState<ServiceUpsertRequest>(() => emptyUpsertRequest());
  const [selectedPreset, setSelectedPreset] = useState<ServicePreset | null>(null);
  const [startAfter, setStartAfter] = useState(true);
  const [saving, setSaving] = useState(false);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
      <PresetDropdown
        selectedPreset={selectedPreset}
        onSelect={(preset) => {
          setSelectedPreset(preset);
          setTemplate(preset.service);
        }}
      />

      <Checkbox
        checked={startAfter}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setStartAfter(e.target.checked)}
        label="Start & register with enabled agents immediately after creating"
      />

      <ServiceEditorForm
        initial={template}
        isNew
        saving={saving}
        onCancel={onCancel}
        onSubmit={async (req) => {
          setSaving(true);
          try {
            await onCreate(req, startAfter);
          } finally {
            setSaving(false);
          }
        }}
      />
    </div>
  );
};
