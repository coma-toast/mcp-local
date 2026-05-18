package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

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
	confirmState
)

type model struct {
	state    sessionState
	services []config.ServiceConfig
	cursor   int
	inputs   []textinput.Model
	focused  int
}

func initialModel(existing []config.ServiceConfig) model {
	services := make([]config.ServiceConfig, len(existing))
	copy(services, existing)

	inputs := make([]textinput.Model, 7)
	fields := []string{"Name", "Command", "Port", "Type (http/stdio)", "Env (K=V,K=V)", "Build Command", "Dependencies (csv)"}

	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].Placeholder = fields[i]
		inputs[i].Focus()
		if i != 0 {
			inputs[i].Blur()
		}
	}

	return model{
		state:    listState,
		services: services,
		cursor:   0,
		inputs:   inputs,
		focused:  0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

		if m.state == listState {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.services)+2 {
					m.cursor++
				}
			case "enter":
				if m.cursor == len(m.services) {
					m.state = formState
					for i := range m.inputs {
						m.inputs[i].Blur()
					}
					m.inputs[0].Focus()
					return m, nil
				} else if m.cursor == len(m.services)+1 {
					m.state = confirmState
					return m, nil
				} else {
					idx := m.cursor
					m.services = append(m.services[:idx], m.services[idx+1:]...)
					if m.cursor >= len(m.services)+2 {
						m.cursor = len(m.services) + 1
					}
				}
			}
		} else if m.state == formState {
			switch msg.String() {
			case "tab":
				m.inputs[m.focused].Blur()
				m.focused = (m.focused + 1) % len(m.inputs)
				m.inputs[m.focused].Focus()
			case "shift+tab":
				m.inputs[m.focused].Blur()
				m.focused = (m.focused - 1 + len(m.inputs)) % len(m.inputs)
				m.inputs[m.focused].Focus()
			case "enter":
				s := config.ServiceConfig{
					Name:     m.inputs[0].Value(),
					Command:  m.inputs[1].Value(),
					MCPType:  m.inputs[3].Value(),
					BuildCmd: m.inputs[5].Value(),
				}
				port, _ := strconv.Atoi(m.inputs[2].Value())
				s.Port = port

				envStr := m.inputs[4].Value()
				if envStr != "" {
					s.Env = make(map[string]string)
					pairs := strings.Split(envStr, ",")
					for _, p := range pairs {
						kv := strings.Split(p, "=")
						if len(kv) == 2 {
							s.Env[kv[0]] = kv[1]
						}
					}
				}

				depsStr := m.inputs[6].Value()
				if depsStr != "" {
					s.Deps = strings.Split(depsStr, ",")
				}

				m.services = append(m.services, s)
				m.state = listState
				for i := range m.inputs {
					m.inputs[i].SetValue("")
				}
				return m, nil
			}
			var inputCmd tea.Cmd
			m.inputs[m.focused], inputCmd = m.inputs[m.focused].Update(msg)
			return m, inputCmd
		} else if m.state == confirmState {
			switch msg.String() {
			case "enter":
				return m, tea.Quit
			case "esc":
				m.state = listState
				return m, nil
			}
		}
	}

	if m.state == formState {
		var inputCmd tea.Cmd
		m.inputs[m.focused], inputCmd = m.inputs[m.focused].Update(msg)
		return m, inputCmd
	}
	return m, cmd
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("MCP-LOCAL CONFIG WIZARD") + "\n\n")

	switch m.state {
	case listState:
		s.WriteString("Current Services:\n")
		for i, svc := range m.services {
			cursor := "  "
			style := itemStyle
			if m.cursor == i {
				cursor = "> "
				style = selectedStyle
			}
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(svc.Name)))
		}

		s.WriteString("\n")

		addCursor := "  "
		addStyle := itemStyle
		if m.cursor == len(m.services) {
			addCursor = "> "
			addStyle = selectedStyle
		}
		s.WriteString(fmt.Sprintf("%s%s\n", addCursor, addStyle.Render("[Add Service]")))

		saveCursor := "  "
		saveStyle := itemStyle
		if m.cursor == len(m.services)+1 {
			saveCursor = "> "
			saveStyle = selectedStyle
		}
		s.WriteString(fmt.Sprintf("%s%s\n", saveCursor, saveStyle.Render("[Save & Exit]")))

		s.WriteString("\n(Use arrow keys to navigate, Enter to select, Ctrl+C to quit)")

	case formState:
		s.WriteString("Add New Service:\n\n")
		fields := []string{"Name", "Command", "Port", "Type (http/stdio)", "Env (K=V,K=V)", "Build Command", "Dependencies (csv)"}

		for i, input := range m.inputs {
			label := itemStyle.Render(fields[i] + ": ")
			if m.focused == i {
				label = selectedStyle.Render(fields[i] + ": ")
			}
			s.WriteString(label + " " + input.View() + "\n")
		}
		s.WriteString("\n(Tab to switch fields, Enter to add service)")

	case confirmState:
		s.WriteString("Confirm Configuration:\n\n")
		for _, svc := range m.services {
			s.WriteString(fmt.Sprintf("- %s (%s:%d)\n", svc.Name, svc.MCPType, svc.Port))
		}
		s.WriteString("\n")
		s.WriteString(selectedStyle.Render("[Press Enter to Save and Exit]"))
	}

	return borderStyle.Render(s.String())
}

func RunConfigWizard() (*config.ManagerConfig, error) {
	existing, _ := config.LoadConfig()
	var services []config.ServiceConfig
	if existing != nil {
		services = existing.Services
	}
	p := tea.NewProgram(initialModel(services))
	m, err := p.Run()
	if err != nil {
		return nil, err
	}

	finalModel := m.(model)
	if finalModel.state != confirmState {
		return nil, fmt.Errorf("wizard cancelled")
	}

	return &config.ManagerConfig{Services: finalModel.services}, nil
}
