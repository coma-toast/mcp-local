package asttools

// DefaultBuiltinTier returns ast-context-cache's compiled-in tier for a tool name.
// Empty string means unknown tool — omit tier override until the user sets one in the TUI.
// Keep in sync with internal/mcp/tools.go Tier fields in github.com/coma-toast/ast-context-cache.
func DefaultBuiltinTier(name string) string {
	switch name {
	case "get_context_capsule", "index_status", "search_semantic", "get_project_map",
		"get_file_context", "get_impact_graph", "diff_impact", "check_symbol_exists",
		"check_deletion_safety", "search_docs", "list_doc_sources", "retrieve",
		"fetch_context", "list_context", "search_context", "recall_memory",
		"handoff", "open_handoff", "scratchpad":
		return "core"
	case "execute_code":
		return "complete"
	case "index_files", "cache_summary", "analyze_dead_code", "analyze_complexity",
		"export_bundle", "import_bundle", "fetch_doc", "add_doc_source", "remove_doc_source",
		"update_doc_source", "store_context", "flush_context", "store_memory",
		"forget_memory", "report_kv_repair_event":
		return "extended"
	default:
		return ""
	}
}

// EffectiveTierForDisplay is the tier shown in the TUI and used when cycling tiers on an unset override.
func EffectiveTierForDisplay(name, configured string) string {
	if configured != "" {
		return configured
	}
	return DefaultBuiltinTier(name)
}
