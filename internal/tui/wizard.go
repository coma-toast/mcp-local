package tui

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

// ErrCancelled is returned by the TUIs when the user quits without saving.
var ErrCancelled = errors.New("cancelled, no changes saved")

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B4B4B4"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	errStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F5F"))

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1).
			MarginLeft(2)
)

type sessionState int

const (
	listState sessionState = iota
	formState
	deleteState
	confirmState
)

var wizardFields = []string{"Name", "Command", "Port", "Type (http/stdio)", "Env (K=V,K=V)", "Build Command", "Dependencies (csv)"}

type model struct {
	state     sessionState
	services  []config.ServiceConfig
	cursor    int
	inputs    []textinput.Model
	focused   int
	editIdx   int // index of the service the form edits, -1 when adding
	errMsg    string
	done      bool
	cancelled bool
}

func initialModel(existing []config.ServiceConfig) model {
	services := make([]config.ServiceConfig, len(existing))
	copy(services, existing)
	inputs := make([]textinput.Model, len(wizardFields))
	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].Placeholder = wizardFields[i]
	}
	inputs[0].Focus()
	return model{state: listState, services: services, inputs: inputs, editIdx: -1}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// lastRow is the cursor index of the final list row ([Save & Exit]).
func (m model) lastRow() int { return len(m.services) + 1 }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		if m.state == formState {
			var cmd tea.Cmd
			m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
			return m, cmd
		}
		return m, nil
	}
	switch key.String() {
	case "ctrl+c", "esc":
		m.cancelled = true
		return m, tea.Quit
	}
	switch m.state {
	case listState:
		return m.updateList(key)
	case formState:
		return m.updateForm(key)
	case deleteState:
		switch key.String() {
		case "y", "Y":
			if m.cursor < len(m.services) {
				m.services = append(m.services[:m.cursor], m.services[m.cursor+1:]...)
			}
			m.cursor = min(m.cursor, m.lastRow())
			m.state = listState
		case "n", "N":
			m.state = listState
		}
	case confirmState:
		switch key.String() {
		case "enter":
			m.done = true
			return m, tea.Quit
		case "n", "backspace":
			m.state = listState
		}
	}
	return m, nil
}

func (m model) updateList(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < m.lastRow() {
			m.cursor++
		}
	case "d", "delete":
		if m.cursor < len(m.services) {
			m.state = deleteState
		}
	case "enter":
		switch {
		case m.cursor < len(m.services):
			m.openForm(m.cursor)
		case m.cursor == len(m.services):
			m.openForm(-1)
		default:
			m.state = confirmState
		}
	}
	return m, nil
}

func (m *model) openForm(idx int) {
	m.editIdx, m.errMsg, m.state = idx, "", formState
	vals := make([]string, len(m.inputs))
	if idx >= 0 {
		s := m.services[idx]
		vals = []string{s.Name, s.Command, "", s.MCPType, formatEnv(s.Env), config.EffectiveBuild(s), strings.Join(s.Deps, ",")}
		if s.Port != 0 {
			vals[2] = strconv.Itoa(s.Port)
		}
	}
	for i := range m.inputs {
		m.inputs[i].SetValue(vals[i])
		m.inputs[i].Blur()
	}
	m.focused = 0
	m.inputs[0].Focus()
}

func (m model) updateForm(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "tab", "down":
		m.focus((m.focused + 1) % len(m.inputs))
		return m, nil
	case "shift+tab", "up":
		m.focus((m.focused - 1 + len(m.inputs)) % len(m.inputs))
		return m, nil
	case "enter":
		base := config.ServiceConfig{}
		if m.editIdx >= 0 {
			base = m.services[m.editIdx]
		}
		s, err := m.formService(base)
		if err != nil {
			m.errMsg = err.Error()
			return m, nil
		}
		if m.editIdx >= 0 {
			m.services[m.editIdx] = s
		} else {
			m.services = append(m.services, s)
		}
		m.state, m.errMsg = listState, ""
		return m, nil
	}
	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(key)
	return m, cmd
}

func (m *model) focus(i int) {
	m.inputs[m.focused].Blur()
	m.focused = i
	m.inputs[m.focused].Focus()
}

// formService applies the form's values on top of base, so fields the form
// does not show (args, tools, URLs, ...) survive an edit.
func (m model) formService(base config.ServiceConfig) (config.ServiceConfig, error) {
	s := base
	s.Name = strings.TrimSpace(m.inputs[0].Value())
	if s.Name == "" {
		return s, fmt.Errorf("name is required")
	}
	for i, other := range m.services {
		if i != m.editIdx && other.Name == s.Name {
			return s, fmt.Errorf("a service named %q already exists", s.Name)
		}
	}
	port, err := parsePort(m.inputs[2].Value())
	if err != nil {
		return s, err
	}
	s.Command, s.Port, s.MCPType = m.inputs[1].Value(), port, m.inputs[3].Value()
	// Re-parsing an untouched env would split values that contain commas.
	if envStr := m.inputs[4].Value(); envStr != formatEnv(base.Env) {
		s.Env = parseEnv(envStr)
	}
	setBuild(&s, m.inputs[5].Value())
	s.Deps = nil
	for _, d := range strings.Split(m.inputs[6].Value(), ",") {
		if d = strings.TrimSpace(d); d != "" {
			s.Deps = append(s.Deps, d)
		}
	}
	return s, nil
}

// parsePort accepts an empty string as "no port" (stdio services).
func parsePort(v string) (int, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, nil
	}
	port, err := strconv.Atoi(v)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid port %q: must be a number from 1 to 65535", v)
	}
	return port, nil
}

// parseEnv parses K=V,K=V; only the first "=" splits, so values may contain "=".
func parseEnv(v string) map[string]string {
	var env map[string]string
	for _, p := range strings.Split(v, ",") {
		kv := strings.SplitN(p, "=", 2)
		k := strings.TrimSpace(kv[0])
		if len(kv) != 2 || k == "" {
			continue
		}
		if env == nil {
			env = map[string]string{}
		}
		env[k] = kv[1]
	}
	return env
}

func formatEnv(env map[string]string) string {
	pairs := make([]string, 0, len(env))
	for k, v := range env {
		pairs = append(pairs, k+"="+v)
	}
	sort.Strings(pairs)
	return strings.Join(pairs, ",")
}

func (m model) View() string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("MCP-LOCAL CONFIG WIZARD") + "\n\n")
	row := func(i int, label string) {
		cursor, style := "  ", itemStyle
		if m.cursor == i {
			cursor, style = "> ", selectedStyle
		}
		s.WriteString(cursor + style.Render(label) + "\n")
	}
	switch m.state {
	case listState:
		s.WriteString("Current Services:\n")
		for i, svc := range m.services {
			row(i, svc.Name)
		}
		s.WriteString("\n")
		row(len(m.services), "[Add Service]")
		row(m.lastRow(), "[Save & Exit]")
		s.WriteString("\n(↑/↓ navigate, Enter select/edit, d delete, Esc quit without saving)")
	case formState:
		title := "Add New Service:"
		if m.editIdx >= 0 {
			title = "Edit Service:"
		}
		s.WriteString(title + "\n\n")
		for i, input := range m.inputs {
			label := itemStyle.Render(wizardFields[i] + ": ")
			if m.focused == i {
				label = selectedStyle.Render(wizardFields[i] + ": ")
			}
			s.WriteString(label + " " + input.View() + "\n")
		}
		if m.errMsg != "" {
			s.WriteString("\n" + errStyle.Render("⚠️  "+m.errMsg) + "\n")
		}
		s.WriteString("\n(Tab/↑/↓ switch fields, Enter save service, Esc quit without saving)")
	case deleteState:
		if m.cursor < len(m.services) {
			s.WriteString(fmt.Sprintf("Delete service %q? (y/n)", m.services[m.cursor].Name))
		}
	case confirmState:
		s.WriteString("Confirm Configuration:\n\n")
		for _, svc := range m.services {
			s.WriteString(fmt.Sprintf("- %s (%s:%d)\n", svc.Name, svc.MCPType, svc.Port))
		}
		s.WriteString("\n" + selectedStyle.Render("[Enter save and exit, n back, Esc quit without saving]"))
	}
	return borderStyle.Render(s.String())
}

// wizardResult turns the final model into a config, keeping every non-service
// setting (e.g. agents) from existing.
func wizardResult(existing *config.ManagerConfig, final model) (*config.ManagerConfig, error) {
	if !final.done || final.cancelled {
		return nil, ErrCancelled
	}
	cfg := &config.ManagerConfig{}
	if existing != nil {
		*cfg = *existing
	}
	cfg.Services = final.services
	return cfg, nil
}

func RunConfigWizard() (*config.ManagerConfig, error) {
	existing, _ := config.LoadConfig()
	var services []config.ServiceConfig
	if existing != nil {
		services = existing.Services
	}
	m, err := tea.NewProgram(initialModel(services)).Run()
	if err != nil {
		return nil, err
	}
	return wizardResult(existing, m.(model))
}
