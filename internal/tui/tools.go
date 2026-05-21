package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

type toolState int

const (
	serviceListState toolState = iota
	toolListState
	descEditState
)

type toolModel struct {
	state   toolState
	cfg     *config.ManagerConfig
	svcIdx  int
	toolIdx int
	cursor  int
	input   textinput.Model
	focused bool
}

func initialToolModel(cfg *config.ManagerConfig) toolModel {
	ti := textinput.New()
	ti.Placeholder = "Enter tool description..."
	ti.Focus()

	return toolModel{
		state:   serviceListState,
		cfg:     cfg,
		svcIdx:  0,
		toolIdx: 0,
		cursor:  0,
		input:   ti,
		focused: false,
	}
}

func (m toolModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m toolModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.state == descEditState {
				m.state = toolListState
				return m, nil
			}
			return m, tea.Quit
		}

		if m.state == serviceListState {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.cfg.Services)-1 {
					m.cursor++
				}
			case "enter":
				m.svcIdx = m.cursor
				m.state = toolListState
				m.cursor = 0
				return m, nil
			}
		} else if m.state == toolListState {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.cfg.Services[m.svcIdx].Tools)-1 {
					m.cursor++
				}
			case "left", "h":
				m.cfg.Services[m.svcIdx].Tools[m.cursor].Enabled = !m.cfg.Services[m.svcIdx].Tools[m.cursor].Enabled
			case "right", "l":
				tiers := []string{"core", "extended", "complete"}
				currentTier := m.cfg.Services[m.svcIdx].Tools[m.cursor].Tier
				nextIdx := 0
				for i, t := range tiers {
					if t == currentTier {
						nextIdx = (i + 1) % len(tiers)
						break
					}
				}
				m.cfg.Services[m.svcIdx].Tools[m.cursor].Tier = tiers[nextIdx]
			case "enter":
				m.toolIdx = m.cursor
				m.input.SetValue(m.cfg.Services[m.svcIdx].Tools[m.cursor].Description)
				m.state = descEditState
				return m, textinput.Blink
			case "backspace", "esc":
				m.state = serviceListState
				m.cursor = 0
				return m, nil
			}
		} else if m.state == descEditState {
			switch msg.String() {
			case "enter":
				m.cfg.Services[m.svcIdx].Tools[m.toolIdx].Description = m.input.Value()
				m.state = toolListState
				return m, nil
			}
			var inputCmd tea.Cmd
			m.input, inputCmd = m.input.Update(msg)
			return m, inputCmd
		}
	}

	if m.state == descEditState {
		var inputCmd tea.Cmd
		m.input, inputCmd = m.input.Update(msg)
		return m, inputCmd
	}

	return m, cmd
}

func (m toolModel) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("MCP-LOCAL TOOL MANAGER") + "\n\n")

	switch m.state {
	case serviceListState:
		s.WriteString("Select a service to manage its tools:\n\n")
		for i, svc := range m.cfg.Services {
			cursor := "  "
			style := itemStyle
			if m.cursor == i {
				cursor = "> "
				style = selectedStyle
			}
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(svc.Name)))
		}
		s.WriteString("\n(Use arrow keys, Enter to select, Ctrl+C to quit)")

	case toolListState:
		svc := m.cfg.Services[m.svcIdx]
		s.WriteString(fmt.Sprintf("Managing tools for: %s\n\n", selectedStyle.Render(svc.Name)))

		if len(svc.Tools) == 0 {
			s.WriteString("No tools defined for this service.\n")
			s.WriteString("Run: mcp-local tools sync " + svc.Name + "\n")
			s.WriteString("\n(Press Backspace to return)")
		} else {
			s.WriteString(fmt.Sprintf("%-20s %-10s %-15s %s\n", "TOOL", "ENABLED", "TIER", "DESCRIPTION"))
			s.WriteString(strings.Repeat("-", 60) + "\n")

			for i, tool := range svc.Tools {
				cursor := "  "
				style := itemStyle
				if m.cursor == i {
					cursor = "> "
					style = selectedStyle
				}

				enabled := "❌"
				if tool.Enabled {
					enabled = "✅"
				}

				s.WriteString(fmt.Sprintf("%s%s %-10s %-15s %s\n",
					cursor, style.Render(tool.Name), enabled, tool.Tier, tool.Description))
			}
			s.WriteString("\n(Left/Right: Toggle/Tier, Enter: Edit Desc, Backspace: Back)")
			s.WriteString("\nAfter save: mcp-local tools apply " + svc.Name + " && mcp-local restart " + svc.Name)
		}

	case descEditState:
		s.WriteString("Edit Tool Description:\n\n")
		s.WriteString(m.input.View() + "\n")
		s.WriteString("\n(Enter to save, Esc to cancel)")
	}

	return borderStyle.Render(s.String())
}

func RunToolManager(cfg *config.ManagerConfig) error {
	p := tea.NewProgram(initialToolModel(cfg))
	m, err := p.Run()
	if err != nil {
		return err
	}

	finalModel := m.(toolModel)
	if err := config.SaveConfig(finalModel.cfg); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Println("\n⚠️  Changes saved to config.yaml. Restart services to apply env/tool settings.")
	return nil
}
