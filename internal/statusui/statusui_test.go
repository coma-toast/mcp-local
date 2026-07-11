package statusui

import (
	"testing"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

func TestPrintPlain(t *testing.T) {
	cfg := &config.ManagerConfig{
		Services: []config.ServiceConfig{
			{Name: "service1", Port: 8080, MCPType: "http"},
			{Name: "service2", Port: 0, MCPType: "stdio"},
		},
	}

	// Note: PrintPlain writes to stdout, we can't easily capture it
	// Just test that it doesn't panic
	PrintPlain(cfg, "")

	// Test with filter
	PrintPlain(cfg, "service1")
}

func TestPrintPlain_EmptyConfig(t *testing.T) {
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{}}
	PrintPlain(cfg, "")
}
