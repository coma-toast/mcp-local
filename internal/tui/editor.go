package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

type editField struct {
	label string
	value string
}

type editModel struct {
	fields    []editField
	inputs    []textinput.Model
	cursor    int
	svc       config.ServiceConfig
	errMsg    string
	saved     bool
	cancelled bool
}

func initialEditModel(svc *config.ServiceConfig) editModel {
	port := ""
	if svc.Port != 0 {
		port = strconv.Itoa(svc.Port)
	}
	fields := []editField{
		{"Name", svc.Name},
		{"Command", svc.Command},
		{"Port", port},
		{"Type (http/stdio)", svc.MCPType},
		{"MCP URL", svc.MCPURL},
		{"Health URL", svc.HealthURL},
		{"Dashboard URL", svc.DashboardURL},
		{"Log Path", svc.Log},
		{"Build Command", config.EffectiveBuild(*svc)},
		{"Description", svc.Description},
		{"Active Tier", svc.ActiveTier},
	}
	inputs := make([]textinput.Model, len(fields))
	for i, f := range fields {
		inputs[i] = textinput.New()
		inputs[i].Placeholder = f.label
		inputs[i].SetValue(f.value)
	}
	inputs[0].Focus()
	return editModel{fields: fields, inputs: inputs, svc: *svc}
}

func (m editModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update treats every printable key as text; only arrows and tab move focus.
func (m editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "up", "shift+tab":
			m.move((m.cursor - 1 + len(m.inputs)) % len(m.inputs))
			return m, nil
		case "down", "tab":
			m.move((m.cursor + 1) % len(m.inputs))
			return m, nil
		case "enter":
			svc, err := m.result()
			if err != nil {
				m.errMsg = err.Error()
				return m, nil
			}
			m.svc, m.saved = svc, true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.inputs[m.cursor], cmd = m.inputs[m.cursor].Update(msg)
	return m, cmd
}

func (m *editModel) move(i int) {
	m.inputs[m.cursor].Blur()
	m.cursor = i
	m.inputs[m.cursor].Focus()
}

// result applies the inputs to a copy of the original service.
func (m editModel) result() (config.ServiceConfig, error) {
	s := m.svc
	port, err := parsePort(m.inputs[2].Value())
	if err != nil {
		return s, err
	}
	if s.Name = strings.TrimSpace(m.inputs[0].Value()); s.Name == "" {
		return s, fmt.Errorf("name is required")
	}
	s.Command, s.Port, s.MCPType = m.inputs[1].Value(), port, m.inputs[3].Value()
	s.MCPURL, s.HealthURL, s.DashboardURL = m.inputs[4].Value(), m.inputs[5].Value(), m.inputs[6].Value()
	s.Log = m.inputs[7].Value()
	setBuild(&s, m.inputs[8].Value())
	s.Description, s.ActiveTier = m.inputs[9].Value(), m.inputs[10].Value()
	return s, nil
}

// setBuild stores v in whichever build field s already uses (build_cmd wins
// over build_command, see config.EffectiveBuild), defaulting to build_command,
// so that config.EffectiveBuild(*s) == v afterwards.
func setBuild(s *config.ServiceConfig, v string) {
	if s.BuildCmd == "" {
		s.BuildCommand = v
		return
	}
	s.BuildCmd = v
	if v == "" {
		s.BuildCommand = ""
	}
}

func (m editModel) View() string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("EDIT SERVICE") + "\n\n")
	for i, f := range m.fields {
		style := itemStyle
		if m.cursor == i {
			style = selectedStyle
		}
		s.WriteString(fmt.Sprintf("  %s %s\n", style.Render(f.label+":"), m.inputs[i].View()))
	}
	if m.errMsg != "" {
		s.WriteString("\n" + errStyle.Render("⚠️  "+m.errMsg) + "\n")
	}
	s.WriteString("\n(↑/↓/Tab navigate, Enter to save, Esc to cancel)")
	return borderStyle.Render(s.String())
}

// RunServiceEditor edits svc in place; it returns ErrCancelled, leaving svc
// untouched, when the user quits without saving.
func RunServiceEditor(svc *config.ServiceConfig) error {
	result, err := tea.NewProgram(initialEditModel(svc)).Run()
	if err != nil {
		return err
	}
	final := result.(editModel)
	if !final.saved || final.cancelled {
		return ErrCancelled
	}
	*svc = final.svc
	return nil
}
