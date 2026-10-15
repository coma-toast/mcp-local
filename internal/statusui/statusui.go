package statusui

import (
	"fmt"
	"io"
	"os"
	"strconv"
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
)

type tickMsg time.Time

type row struct {
	name    string
	running bool
	pid     int
	port    int
	uptime  string
	mode    process.Mode
	stale   bool
}

// label is the STATUS column: services mcp-local never runs say why instead of "stopped".
func (r row) label() string {
	switch {
	case r.running:
		return "running"
	case r.mode == process.Remote:
		return "remote"
	case r.mode == process.Stdio:
		return "on-demand"
	case r.stale:
		return "stale"
	}
	return "stopped"
}

type model struct {
	byName   map[string]config.ServiceConfig
	order    []string
	filter   string
	rows     []row
	quitting bool
	width    int // terminal width from tea.WindowSizeMsg; 0 until reported
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
			rows = append(rows, row{name: st.Name, running: st.Running, pid: st.PID, port: st.Port, uptime: st.Uptime, mode: st.Mode, stale: st.Stale})
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
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
	cols := fitColumns(m.rows, m.width)
	header := pad("NAME", cols.name) + " " + pad("STATUS", colStatus)
	if cols.pid {
		header += " " + pad("PID", colPID)
	}
	header += " " + pad("PORT", colPort)
	if cols.uptime {
		header += " UPTIME"
	}
	out := headerStyle.Render(header) + "\n"
	for _, r := range m.rows {
		style := stoppedStyle
		switch r.label() {
		case "running":
			style = runningStyle
		case "remote", "on-demand":
			style = dimStyle
		}
		line := pad(cut(r.name, cols.name), cols.name) + " " + style.Render(pad(r.label(), colStatus))
		if cols.pid {
			if r.pid > 0 {
				line += " " + pad(strconv.Itoa(r.pid), colPID)
			} else {
				line += " " + dimStyle.Render(pad("-", colPID))
			}
		}
		if r.port > 0 {
			line += " " + pad(strconv.Itoa(r.port), colPort)
		} else {
			line += " " + dimStyle.Render(pad("-", colPort))
		}
		if cols.uptime && r.running {
			line += " " + humanUptime(r.uptime)
		}
		out += line + "\n"
	}
	out += dimStyle.Render("\nPress q to quit")
	return out
}

const colStatus, colPID, colPort, colUptime = 9, 7, 5, 8

type columns struct {
	name        int
	pid, uptime bool
}

// fitColumns sizes NAME to the longest service name (max 32) and drops UPTIME,
// then PID, when the terminal is too narrow; NAME is cut last. width 0 = unknown.
func fitColumns(rows []row, width int) columns {
	c := columns{name: len("NAME"), pid: true, uptime: true}
	for _, r := range rows {
		if len(r.name) > c.name {
			c.name = len(r.name)
		}
	}
	if c.name > 32 {
		c.name = 32
	}
	if width <= 0 {
		return c
	}
	total := func() int {
		t := c.name + 1 + colStatus + 1 + colPort
		if c.pid {
			t += 1 + colPID
		}
		if c.uptime {
			t += 1 + colUptime
		}
		return t
	}
	if total() > width {
		c.uptime = false
	}
	if total() > width {
		c.pid = false
	}
	if over := total() - width; over > 0 && c.name-over >= 8 {
		c.name -= over
	}
	return c
}

func pad(s string, width int) string {
	if n := len([]rune(s)); n < width {
		return s + strings.Repeat(" ", width-n)
	}
	return s
}

// cut shortens s to width runes, ending in "…" when it had to cut.
func cut(s string, width int) string {
	r := []rune(s)
	if len(r) <= width || width < 2 {
		return s
	}
	return string(r[:width-1]) + "…"
}

// humanUptime turns ps elapsed time ("[[dd-]hh:]mm:ss" on macOS, plain seconds
// on Linux) into "6d 2h", "3h 5m", "12m" or "40s"; anything else is returned as is.
func humanUptime(s string) string {
	s = strings.TrimSpace(s)
	secs, ok := elapsedSeconds(s)
	if !ok {
		return s
	}
	d, h, m := secs/86400, secs%86400/3600, secs%3600/60
	switch {
	case d > 0:
		return fmt.Sprintf("%dd %dh", d, h)
	case h > 0:
		return fmt.Sprintf("%dh %dm", h, m)
	case m > 0:
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%ds", secs)
}

func elapsedSeconds(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n, n >= 0
	}
	days := 0
	if i := strings.Index(s, "-"); i > 0 {
		d, err := strconv.Atoi(s[:i])
		if err != nil {
			return 0, false
		}
		days, s = d, s[i+1:]
	}
	parts := strings.Split(s, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return 0, false
	}
	total := 0
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return 0, false
		}
		total = total*60 + n
	}
	return days*86400 + total, true
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
	colName, colStatus, colPID, colPort := 24, 10, 8, 7
	fmt.Fprintf(w, "%-*s %-*s %-*s %-*s\n", colName, "NAME", colStatus, "STATUS", colPID, "PID", colPort, "PORT")
	fmt.Fprintln(w, "------------------------------------------------------------")
	for _, r := range rows {
		pid := "-"
		if r.pid > 0 {
			pid = fmt.Sprintf("%d", r.pid)
		}
		fmt.Fprintf(w, "%-*s %-*s %-*s %-*d\n", colName, r.name, colStatus, r.label(), colPID, pid, colPort, r.port)
	}
}
