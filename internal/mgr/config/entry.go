package config

import (
	"fmt"
	"strings"
)

type AgentEntry struct {
	Type        string
	URL         string
	Command     []string
	Enabled     bool
	Environment map[string]string
}

func ServiceToEntry(s ServiceConfig) (AgentEntry, error) {
	if s.IsHTTP() {
		url := strings.TrimSpace(s.MCPURL)
		if url == "" {
			url = fmt.Sprintf("http://localhost:%d/mcp", s.Port)
		}
		return AgentEntry{
			Type:    "remote",
			URL:     url,
			Enabled: true,
		}, nil
	}

	cmd := []string{ExpandPath(s.Command)}
	cmd = append(cmd, s.Args...)
	entry := AgentEntry{
		Type:    "local",
		Command: cmd,
		Enabled: true,
	}
	if len(s.Env) > 0 {
		entry.Environment = make(map[string]string, len(s.Env))
		for k, v := range s.Env {
			entry.Environment[k] = v
		}
	}
	return entry, nil
}
