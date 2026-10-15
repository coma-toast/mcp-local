package statusui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
)

func TestPrintPlain(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cfg := &config.ManagerConfig{
		Services: []config.ServiceConfig{
			{Name: "service1", Port: 0, MCPType: "http", Command: "srv", MCPURL: "http://localhost:1/mcp"},
			{Name: "service2", Port: 0, MCPType: "stdio", Command: "npx"},
			{Name: "service3", MCPURL: "https://example.com/mcp"},
		},
	}
	var buf bytes.Buffer
	FprintPlain(&buf, cfg, "")
	out := buf.String()
	for _, want := range []string{"service1                 stopped", "service2                 on-demand", "service3                 remote"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	buf.Reset()
	FprintPlain(&buf, cfg, "service1")
	if strings.Contains(buf.String(), "service2") {
		t.Errorf("filter ignored:\n%s", buf.String())
	}
}

func TestPrintPlain_EmptyConfig(t *testing.T) {
	PrintPlain(&config.ManagerConfig{}, "")
}

func TestRowLabelAndDetail(t *testing.T) {
	cases := []struct {
		r            row
		label, extra string
	}{
		{row{state: config.StateRunning, detail: "healthy", running: true}, "running", ""},
		{row{state: config.StateOrphaned, detail: "port 7830 held by an untracked ast-mcp (pid 1)"}, "orphaned", "port 7830 held by an untracked ast-mcp (pid 1)"},
		{row{state: config.StateStopped, detail: "stopped", mode: process.Remote}, "remote", ""},
		{row{state: config.StateStopped, detail: "stopped", mode: process.Stdio}, "on-demand", ""},
		{row{state: config.StateStale, detail: "PID file names pid 9"}, "stale", "PID file names pid 9"},
	}
	for _, c := range cases {
		if c.r.label() != c.label || c.r.detailText() != c.extra {
			t.Errorf("%+v: label=%q detail=%q, want %q/%q", c.r, c.r.label(), c.r.detailText(), c.label, c.extra)
		}
	}
}
