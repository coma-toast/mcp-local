package statusui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
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
