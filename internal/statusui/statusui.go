package statusui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/runtime"
)

var (
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	runningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	stoppedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type tickMsg time.Time

type row struct {
	name    string
	running bool
	pid     int
	port    int
	uptime  string
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
			rows = append(rows, row{name: st.Name, running: st.Running, pid: st.PID, port: st.Port, uptime: st.Uptime})
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
	colName, colStatus, colPID, colPort, colUptime := 24, 10, 8, 7, 12
	header := fmt.Sprintf("%-*s %-*s %-*s %-*s %-*s",
		colName, "NAME", colStatus, "STATUS", colPID, "PID", colPort, "PORT", colUptime, "UPTIME")
	out := headerStyle.Render(header) + "\n"
	for _, r := range m.rows {
		var statusStr string
		if r.running {
			statusStr = runningStyle.Render(fmt.Sprintf("%-*s", colStatus, "running"))
		} else {
			statusStr = stoppedStyle.Render(fmt.Sprintf("%-*s", colStatus, "stopped"))
		}
		pidStr := dimStyle.Render(fmt.Sprintf("%-*s", colPID, "-"))
		if r.pid > 0 {
			pidStr = fmt.Sprintf("%-*d", colPID, r.pid)
		}
		out += fmt.Sprintf("%-*s %s %s %-*d %-*s\n",
			colName, r.name, statusStr, pidStr, colPort, r.port, colUptime, r.uptime)
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

// PrintPlain prints a one-shot table (script-friendly).
func PrintPlain(cfg *config.ManagerConfig, filter string) {
	m := newModel(cfg, filter)
	var rows []row
	if cmd := m.refresh(); cmd != nil {
		if msg := cmd(); msg != nil {
			rows, _ = msg.([]row)
		}
	}
	colName, colStatus, colPID, colPort := 24, 10, 8, 7
	fmt.Printf("%-*s %-*s %-*s %-*s\n", colName, "NAME", colStatus, "STATUS", colPID, "PID", colPort, "PORT")
	fmt.Println("------------------------------------------------------------")
	for _, r := range rows {
		st := "stopped"
		if r.running {
			st = "running"
		}
		pid := "-"
		if r.pid > 0 {
			pid = fmt.Sprintf("%d", r.pid)
		}
		fmt.Printf("%-*s %-*s %-*s %-*d\n", colName, r.name, colStatus, st, colPID, pid, colPort, r.port)
	}
}
