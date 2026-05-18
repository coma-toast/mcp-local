package tui

import (
	"fmt"
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
	fields   []editField
	inputs   []textinput.Model
	cursor   int
	quitting bool
}

func initialEditModel(svc *config.ServiceConfig) editModel {
	fields := []editField{
		{"Name", svc.Name},
		{"Command", svc.Command},
		{"Port", fmt.Sprintf("%d", svc.Port)},
		{"Type (http/stdio)", svc.MCPType},
		{"MCP URL", svc.MCPURL},
		{"Health URL", svc.HealthURL},
		{"Dashboard URL", svc.DashboardURL},
		{"Log Path", svc.Log},
		{"Build Command", svc.BuildCommand},
		{"Description", svc.Description},
		{"Active Tier", svc.ActiveTier},
	}
	inputs := make([]textinput.Model, len(fields))
	for i, f := range fields {
		inputs[i] = textinput.New()
		inputs[i].Placeholder = f.label
		inputs[i].SetValue(f.value)
		if i == 0 {
			inputs[i].Focus()
		}
	}
	return editModel{
		fields: fields,
		inputs: inputs,
		cursor: 0,
	}
}

func (m editModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.inputs[m.cursor].Blur()
				m.cursor--
				m.inputs[m.cursor].Focus()
			}
			return m, nil
		case "down", "j":
			if m.cursor < len(m.inputs)-1 {
				m.inputs[m.cursor].Blur()
				m.cursor++
				m.inputs[m.cursor].Focus()
			}
			return m, nil
		case "enter":
			m.quitting = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.inputs[m.cursor], cmd = m.inputs[m.cursor].Update(msg)
	return m, cmd
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
	s.WriteString("\n(↑/↓ navigate, Enter to save, Esc to cancel)")
	return borderStyle.Render(s.String())
}

func RunServiceEditor(svc *config.ServiceConfig) error {
	m := initialEditModel(svc)
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return err
	}
	final := result.(editModel)
	if final.quitting && final.cursor == 0 && final.inputs[0].Value() == "" {
		return fmt.Errorf("editor cancelled")
	}
	svc.Name = final.inputs[0].Value()
	svc.Command = final.inputs[1].Value()
	if port := final.inputs[2].Value(); port != "" {
		fmt.Sscanf(port, "%d", &svc.Port)
	}
	svc.MCPType = final.inputs[3].Value()
	svc.MCPURL = final.inputs[4].Value()
	svc.HealthURL = final.inputs[5].Value()
	svc.DashboardURL = final.inputs[6].Value()
	svc.Log = final.inputs[7].Value()
	svc.BuildCommand = final.inputs[8].Value()
	svc.Description = final.inputs[9].Value()
	svc.ActiveTier = final.inputs[10].Value()
	return nil
}
