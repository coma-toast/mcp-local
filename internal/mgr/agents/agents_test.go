package agents

import (
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

func TestDefaultTargets(t *testing.T) {
	targets := DefaultTargets()
	if !targets.OpenCode {
		t.Error("DefaultTargets should have OpenCode=true")
	}
	if !targets.Cursor {
		t.Error("DefaultTargets should have Cursor=true")
	}
	if !targets.Claude {
		t.Error("DefaultTargets should have Claude=true")
	}
}

func TestTargetsFromConfig_WithAgentsConfig(t *testing.T) {
	cfg := config.ManagerConfig{
		Agents: config.AgentsConfig{
			OpenCode: true,
			Cursor:   false,
			Claude:   true,
		},
	}
	targets := TargetsFromConfig(cfg)
	if !targets.OpenCode {
		t.Error("Should respect OpenCode from config")
	}
	if targets.Cursor {
		t.Error("Should respect Cursor=false from config")
	}
	if !targets.Claude {
		t.Error("Should respect Claude from config")
	}
}

func TestTargetsFromConfig_EmptyAgentsConfig(t *testing.T) {
	cfg := config.ManagerConfig{
		Agents: config.AgentsConfig{}, // all zero values
	}
	targets := TargetsFromConfig(cfg)
	if !targets.OpenCode || !targets.Cursor || !targets.Claude {
		t.Error("Should default to all true when no agents config specified")
	}
}

func TestTargetsFromConfig_PartialAgentsConfig(t *testing.T) {
	// When at least one agent is explicitly set, others should default to false
	cfg := config.ManagerConfig{
		Agents: config.AgentsConfig{
			OpenCode: true,
			// Cursor and Claude not set = false
		},
	}
	targets := TargetsFromConfig(cfg)
	if !targets.OpenCode {
		t.Error("Should respect OpenCode=true")
	}
	if targets.Cursor {
		t.Error("Should default Cursor to false when any agent is explicitly configured")
	}
	if targets.Claude {
		t.Error("Should default Claude to false when any agent is explicitly configured")
	}
}
