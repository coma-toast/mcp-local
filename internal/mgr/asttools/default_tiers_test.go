package asttools

import "testing"

func TestDefaultBuiltinTier(t *testing.T) {
	if DefaultBuiltinTier("get_context_capsule") != "core" {
		t.Fatal("expected core")
	}
	if DefaultBuiltinTier("index_files") != "extended" {
		t.Fatal("expected extended")
	}
	if DefaultBuiltinTier("execute_code") != "complete" {
		t.Fatal("expected complete")
	}
	if DefaultBuiltinTier("unknown_tool") != "" {
		t.Fatal("unknown should be empty")
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
