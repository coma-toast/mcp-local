package statusui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/process"
	"github.com/coma-toast/mcp-local/internal/mgr/runtime"
)

var (
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	runningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	stoppedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	orphanStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
)

type tickMsg time.Time

type row struct {
	name    string
	running bool
	pid     int
	port    int
	uptime  string
	mode    process.Mode
	state   config.StateKind
	detail  string
}

// label is the STATUS column: services mcp-local never runs say why instead of "stopped".
func (r row) label() string {
	if r.state == config.StateStopped && !r.running {
		switch r.mode {
		case process.Remote:
			return "remote"
		case process.Stdio:
			return "on-demand"
		}
	}
	if r.state == "" {
		return string(config.StateStopped)
	}
	return string(r.state)
}

// detailText is the DETAIL column; it stays empty when the state says it all.
func (r row) detailText() string {
	switch r.detail {
	case "healthy", "running", "stopped":
		return ""
	}
	return r.detail
}

type model struct {
	byName   map[string]config.ServiceConfig
	order    []string
	filter   string
	rows     []row
	quitting bool
}

func newModel(cfg *config.ManagerConfig, filter string) model {
	byName := make(map[string]config.ServiceConfig)
	order := cfg.SortedNames()
	for _, s := range cfg.Services {
		byName[s.Name] = s
	}
	return model{byName: byName, order: order, filter: filter}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.refresh(), tick())
}

func tick() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) refresh() tea.Cmd {
	return func() tea.Msg {
		var rows []row
		for _, name := range m.order {
			if m.filter != "" && name != m.filter {
				continue
			}
			svc := m.byName[name]
			st := runtime.CheckStatus(svc)
			rows = append(rows, row{name: st.Name, running: st.Running, pid: st.PID, port: st.Port, uptime: st.Uptime, mode: st.Mode, state: st.State, detail: st.Detail})
		}
		return rows
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	case tickMsg:
		return m, tea.Batch(m.refresh(), tick())
	case []row:
		m.rows = msg
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	colName, colStatus, colPID, colPort, colUptime := 24, 13, 8, 7, 12
	header := fmt.Sprintf("%-*s %-*s %-*s %-*s %-*s %s",
		colName, "NAME", colStatus, "STATUS", colPID, "PID", colPort, "PORT", colUptime, "UPTIME", "DETAIL")
	out := headerStyle.Render(header) + "\n"
	for _, r := range m.rows {
		style := stoppedStyle
		switch r.label() {
		case "running":
			style = runningStyle
		case "degraded", "stale":
			style = warnStyle
		case "orphaned":
			style = orphanStyle
		case "remote", "on-demand":
			style = dimStyle
		}
		statusStr := style.Render(fmt.Sprintf("%-*s", colStatus, r.label()))
		pidStr := dimStyle.Render(fmt.Sprintf("%-*s", colPID, "-"))
		if r.pid > 0 {
			pidStr = fmt.Sprintf("%-*d", colPID, r.pid)
		}
		out += fmt.Sprintf("%-*s %s %s %-*d %-*s %s\n",
			colName, r.name, statusStr, pidStr, colPort, r.port, colUptime, r.uptime, dimStyle.Render(r.detailText()))
	}
	out += dimStyle.Render("\nPress q to quit")
	return out
}

// Run starts the bubbletea status UI.
func Run(cfg *config.ManagerConfig, filter string) error {
	m := newModel(cfg, filter)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

// PrintPlain prints a one-shot table (script-friendly) to stdout.
func PrintPlain(cfg *config.ManagerConfig, filter string) {
	FprintPlain(os.Stdout, cfg, filter)
}

func FprintPlain(w io.Writer, cfg *config.ManagerConfig, filter string) {
	m := newModel(cfg, filter)
	var rows []row
	if cmd := m.refresh(); cmd != nil {
		if msg := cmd(); msg != nil {
			rows, _ = msg.([]row)
		}
	}
	colName, colStatus, colPID, colPort := 24, 13, 8, 7
	fmt.Fprintf(w, "%-*s %-*s %-*s %-*s %s\n", colName, "NAME", colStatus, "STATUS", colPID, "PID", colPort, "PORT", "DETAIL")
	fmt.Fprintln(w, "--------------------------------------------------------------------------------")
	for _, r := range rows {
		pid := "-"
		if r.pid > 0 {
			pid = fmt.Sprintf("%d", r.pid)
		}
		line := fmt.Sprintf("%-*s %-*s %-*s %-*d %s", colName, r.name, colStatus, r.label(), colPID, pid, colPort, r.port, r.detailText())
		fmt.Fprintln(w, strings.TrimRight(line, " "))
	}
}
