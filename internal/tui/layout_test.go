package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

func TestVisibleWindow(t *testing.T) {
	cases := []struct{ n, cursor, offset, height, start, end int }{
		{10, 0, 0, 0, 0, 10},  // unknown height: everything
		{5, 4, 0, 10, 0, 5},   // fits
		{30, 0, 0, 10, 0, 10}, // top
		{30, 12, 0, 10, 3, 13},
		{30, 12, 5, 10, 5, 15}, // cursor inside window: keep offset
		{30, 2, 5, 10, 2, 12},  // cursor above window: scroll up
		{30, 29, 0, 10, 20, 30},
	}
	for _, c := range cases {
		s, e := visibleWindow(c.n, c.cursor, c.offset, c.height)
		if s != c.start || e != c.end {
			t.Errorf("visibleWindow(%d,%d,%d,%d) = %d,%d want %d,%d", c.n, c.cursor, c.offset, c.height, s, e, c.start, c.end)
		}
	}
}

func TestPadAndTruncateUseDisplayWidth(t *testing.T) {
	if got := lipgloss.Width(pad("✓ on", 7)); got != 7 {
		t.Errorf("pad width = %d", got)
	}
	if got := truncate("abcdefghij", 5); lipgloss.Width(got) != 5 || !strings.HasSuffix(got, "…") {
		t.Errorf("truncate = %q", got)
	}
}

func TestTierLabel(t *testing.T) {
	if got := tierLabel(config.ToolConfig{Name: "x", Tier: "complete"}); got != "complete" {
		t.Errorf("configured tier: %q", got)
	}
	if got := tierLabel(config.ToolConfig{Name: "get_context_capsule"}); got != "core (default)" {
		t.Errorf("builtin tier: %q", got)
	}
	if got := tierLabel(config.ToolConfig{Name: "some_new_tool"}); got != "server default" {
		t.Errorf("unknown tool: %q", got)
	}
}

func manyTools(n int) []config.ToolConfig {
	tools := make([]config.ToolConfig, n)
	for i := range tools {
		tools[i] = config.ToolConfig{Name: fmt.Sprintf("tool_number_%02d", i), Enabled: i%2 == 0,
			Description: strings.Repeat("a long description that keeps going ", 6)}
	}
	return tools
}

func TestToolTableFitsBoxAndScrolls(t *testing.T) {
	for _, width := range []int{80, 140} {
		m := initialToolModel(&config.ManagerConfig{Services: []config.ServiceConfig{{Name: "svc", Tools: manyTools(40)}}})
		next, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 30})
		m = next.(toolModel)
		m.state = toolListState
		table := m.toolTable(m.services[0].Tools)
		lines := strings.Split(strings.TrimRight(table, "\n"), "\n")
		for _, line := range lines {
			if w := lipgloss.Width(line); w > contentWidth(width) {
				t.Errorf("width %d: line is %d cols, box holds %d: %q", width, w, contentWidth(width), line)
			}
		}
		if rows := len(lines) - 2; rows > m.toolRowsHeight()+1 {
			t.Errorf("width %d: %d rows rendered for a %d-row window", width, rows, m.toolRowsHeight())
		}
		if !strings.Contains(table, "↓") {
			t.Errorf("width %d: expected a '↓ N more' marker for 40 tools in 30 rows", width)
		}
		// moving the cursor to the end scrolls the window and keeps it in view
		for i := 0; i < 39; i++ {
			next, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
			m = next.(toolModel)
		}
		if !strings.Contains(m.toolTable(m.services[0].Tools), "tool_number_39") {
			t.Errorf("width %d: last tool not visible after scrolling to it", width)
		}
	}
}

func TestDescInputWidthTracksTerminal(t *testing.T) {
	m := initialToolModel(&config.ManagerConfig{})
	next, _ := m.Update(tea.WindowSizeMsg{Width: 90, Height: 30})
	if got := next.(toolModel).input.Width; got != contentWidth(90)-3 {
		t.Errorf("input width = %d, want %d", got, contentWidth(90)-3)
	}
}
