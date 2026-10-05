package asttools

import "testing"

func TestDefaultBuiltinTier(t *testing.T) {
	tests := map[string]string{
		"get_context_capsule":    "core",
		"diff_impact":            "core",
		"check_symbol_exists":    "core",
		"check_deletion_safety":  "core",
		"fetch_context":          "core",
		"recall_memory":          "core",
		"handoff":                "core",
		"open_handoff":           "core",
		"scratchpad":             "core",
		"index_files":            "extended",
		"store_context":          "extended",
		"store_memory":           "extended",
		"fetch_doc":              "extended",
		"report_kv_repair_event": "extended",
		"execute_code":           "complete",
		"unknown_tool":           "",
	}
	for name, want := range tests {
		if got := DefaultBuiltinTier(name); got != want {
			t.Errorf("DefaultBuiltinTier(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestEffectiveTierForDisplay(t *testing.T) {
	if EffectiveTierForDisplay("index_files", "") != "extended" {
		t.Fatal("empty config should show builtin")
	}
	if EffectiveTierForDisplay("index_files", "core") != "core" {
		t.Fatal("configured tier wins")
	}
}
