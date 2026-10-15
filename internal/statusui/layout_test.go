package statusui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestHumanUptime(t *testing.T) {
	cases := map[string]string{
		"06-02:49:34": "6d 2h", // macOS ps etime with days
		"02:05:00":    "2h 5m",
		"12:30":       "12m",
		"00:40":       "40s",
		"93784":       "1d 2h", // Linux ps etimes (seconds)
		"-":           "-",
		"?":           "?",
		"":            "",
		"garbage":     "garbage",
	}
	for in, want := range cases {
		if got := humanUptime(in); got != want {
			t.Errorf("humanUptime(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFitColumnsDropsUptimeThenPID(t *testing.T) {
	rows := []row{{name: "ast-context-cache"}, {name: "obsidian"}}
	wide := fitColumns(rows, 120)
	if !wide.pid || !wide.uptime || wide.name != len("ast-context-cache") {
		t.Fatalf("wide: %+v", wide)
	}
	mid := fitColumns(rows, 45)
	if mid.uptime || !mid.pid {
		t.Fatalf("45 cols should drop UPTIME only: %+v", mid)
	}
	narrow := fitColumns(rows, 34)
	if narrow.uptime || narrow.pid {
		t.Fatalf("34 cols should drop UPTIME and PID: %+v", narrow)
	}
	if unknown := fitColumns(rows, 0); !unknown.pid || !unknown.uptime {
		t.Fatalf("unknown width keeps all columns: %+v", unknown)
	}
}

func TestViewFitsTerminalWidth(t *testing.T) {
	m := model{width: 40, rows: []row{
		{name: "ast-context-cache", running: true, pid: 96922, port: 7821, uptime: "06-02:49:34"},
		{name: "atlassian", port: 0},
	}}
	for _, line := range strings.Split(m.View(), "\n") {
		if w := lipgloss.Width(line); w > 40 {
			t.Errorf("line wider than 40 cols (%d): %q", w, line)
		}
	}
	if !strings.Contains(m.View(), "-") {
		t.Errorf("port 0 should render as -")
	}
}
