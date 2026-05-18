package agents

import (
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/cursor"
	"github.com/coma-toast/mcp-local/internal/mgr/opencode"
)

type Targets struct {
	OpenCode bool
	Cursor   bool
}

func DefaultTargets() Targets {
	return Targets{OpenCode: true, Cursor: true}
}

func TargetsFromConfig(cfg config.ManagerConfig) Targets {
	if cfg.Agents.OpenCode || cfg.Agents.Cursor {
		return Targets{
			OpenCode: cfg.Agents.OpenCode,
			Cursor:   cfg.Agents.Cursor,
		}
	}
	return DefaultTargets()
}

func RegisterAll(services []config.ServiceConfig, targets Targets) error {
	if targets.OpenCode {
		if err := opencode.RegisterServices(services); err != nil {
			return err
		}
	}
	if targets.Cursor {
		if err := cursor.RegisterServices(services); err != nil {
			return err
		}
	}
	return nil
}

func DeregisterAll(name string, targets Targets) error {
	if targets.OpenCode {
		if _, err := opencode.Deregister(name); err != nil {
			return err
		}
	}
	if targets.Cursor {
		if err := cursor.Deregister(name); err != nil {
			return err
		}
	}
	return nil
}
