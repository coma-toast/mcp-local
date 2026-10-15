package config

import "fmt"

type AgentEntry struct {
	Type        string
	URL         string
	Headers     map[string]string
	Command     []string
	Enabled     bool
	Environment map[string]string
}

func ServiceToEntry(s ServiceConfig) (AgentEntry, error) {
	if s.IsHTTP() {
		url := s.MCPEndpoint()
		if url == "" {
			url = fmt.Sprintf("http://localhost:%d/mcp", s.Port)
		}
		return AgentEntry{
			Type:    "remote",
			URL:     url,
			Headers: copyMap(s.Headers),
			Enabled: true,
		}, nil
	}
	if s.Command == "" {
		return AgentEntry{}, fmt.Errorf("no command (and no mcp_url) configured")
	}
	cmd := []string{ExpandPath(s.Command)}
	cmd = append(cmd, s.Args...)
	return AgentEntry{
		Type:        "local",
		Command:     cmd,
		Enabled:     true,
		Environment: copyMap(s.Env),
	}, nil
}

func copyMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
