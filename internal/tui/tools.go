package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/coma-toast/mcp-local/internal/mgr/asttools"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

type toolState int

const (
	serviceListState toolState = iota
	toolListState
	descEditState
)

const toolsHelp = "s: save & quit, q/Esc: quit without saving"

type toolModel struct {
	state     toolState
	services  []config.ServiceConfig // working copy; the caller's config is untouched until save
	svcIdx    int
	toolIdx   int
	cursor    int
	input     textinput.Model
	saved     bool
	cancelled bool
}

func initialToolModel(cfg *config.ManagerConfig) toolModel {
	ti := textinput.New()
	ti.Placeholder = "Enter tool description..."
	ti.Focus()
	services := make([]config.ServiceConfig, len(cfg.Services))
	for i, svc := range cfg.Services {
		svc.Tools = append([]config.ToolConfig(nil), svc.Tools...)
		services[i] = svc
	}
	return toolModel{state: serviceListState, services: services, input: ti}
}

func (m toolModel) Init() tea.Cmd {
	return textinput.Blink
}

// tools returns the selected service's tools, or nil when svcIdx is out of range.
func (m toolModel) tools() []config.ToolConfig {
	if m.svcIdx < 0 || m.svcIdx >= len(m.services) {
		return nil
	}
	return m.services[m.svcIdx].Tools
}

func (m toolModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		if m.state == descEditState {
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
		return m, nil
	}
	if key.String() == "ctrl+c" {
		m.cancelled = true
		return m, tea.Quit
	}
	switch m.state {
	case descEditState:
		return m.updateDesc(key)
	case serviceListState, toolListState:
		switch key.String() {
		case "q", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "s":
			m.saved = true
			return m, tea.Quit
		}
	}
	if m.state == serviceListState {
		m.updateServices(key)
	} else {
		return m.updateTools(key)
	}
	return m, nil
}

func (m *toolModel) updateServices(key tea.KeyMsg) {
	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.services)-1 {
			m.cursor++
		}
	case "enter":
		if m.cursor < len(m.services) {
			m.svcIdx, m.state, m.cursor = m.cursor, toolListState, 0
		}
	}
}

func (m toolModel) updateTools(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	tools := m.tools()
	if key.String() == "backspace" {
		m.state, m.cursor = serviceListState, m.svcIdx
		return m, nil
	}
	if m.cursor < 0 || m.cursor >= len(tools) {
		return m, nil
	}
	tool := &tools[m.cursor]
	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(tools)-1 {
			m.cursor++
		}
	case "left", "h":
		tool.Enabled = !tool.Enabled
	case "right", "l":
		tiers := []string{"core", "extended", "complete"}
		currentTier := asttools.EffectiveTierForDisplay(tool.Name, tool.Tier)
		nextIdx := 0
		for i, t := range tiers {
			if t == currentTier {
				nextIdx = (i + 1) % len(tiers)
				break
			}
		}
		tool.Tier = tiers[nextIdx]
	case "enter":
		m.toolIdx = m.cursor
		m.input.SetValue(tool.Description)
		m.state = descEditState
		return m, textinput.Blink
	}
	return m, nil
}

func (m toolModel) updateDesc(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		m.state = toolListState
		return m, nil
	case "enter":
		if tools := m.tools(); m.toolIdx >= 0 && m.toolIdx < len(tools) {
			tools[m.toolIdx].Description = m.input.Value()
		}
		m.state = toolListState
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(key)
	return m, cmd
}

func (m toolModel) View() string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("MCP-LOCAL TOOL MANAGER") + "\n\n")
	switch m.state {
	case serviceListState:
		if len(m.services) == 0 {
			s.WriteString("No services configured. Add one with: mcp-local add\n")
			s.WriteString("\n(q/Esc to quit)")
			break
		}
		s.WriteString("Select a service to manage its tools:\n\n")
		for i, svc := range m.services {
			cursor, style := "  ", itemStyle
			if m.cursor == i {
				cursor, style = "> ", selectedStyle
			}
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(svc.Name)))
		}
		s.WriteString("\n(↑/↓ navigate, Enter select, " + toolsHelp + ")")
	case toolListState:
		if m.svcIdx < 0 || m.svcIdx >= len(m.services) {
			s.WriteString("No service selected.\n\n(Backspace to return)")
			break
		}
		svc := m.services[m.svcIdx]
		s.WriteString(fmt.Sprintf("Managing tools for: %s\n\n", selectedStyle.Render(svc.Name)))
		if len(svc.Tools) == 0 {
			s.WriteString("No tools defined for this service.\n")
			s.WriteString("Run: mcp-local tools sync " + svc.Name + "\n")
			s.WriteString("\n(Backspace to return, " + toolsHelp + ")")
			break
		}
		s.WriteString(fmt.Sprintf("%-20s %-10s %-15s %s\n", "TOOL", "ENABLED", "TIER", "DESCRIPTION"))
		s.WriteString(strings.Repeat("-", 60) + "\n")
		for i, tool := range svc.Tools {
			cursor, style := "  ", itemStyle
			if m.cursor == i {
				cursor, style = "> ", selectedStyle
			}
			enabled := "❌"
			if tool.Enabled {
				enabled = "✅"
			}
			tierLabel := tool.Tier
			if tierLabel == "" {
				tierLabel = asttools.EffectiveTierForDisplay(tool.Name, "") + " (builtin)"
			}
			s.WriteString(fmt.Sprintf("%s%s %-10s %-15s %s\n",
				cursor, style.Render(tool.Name), enabled, tierLabel, tool.Description))
		}
		s.WriteString("\n(Left/Right: Toggle/Tier, Enter: Edit Desc, Backspace: Back, " + toolsHelp + ")")
		s.WriteString("\nAfter save: mcp-local tools apply " + svc.Name + " && mcp-local restart " + svc.Name)
	case descEditState:
		s.WriteString("Edit Tool Description:\n\n")
		s.WriteString(m.input.View() + "\n")
		s.WriteString("\n(Enter to keep, Esc to discard this edit)")
	}
	return borderStyle.Render(s.String())
}

// RunToolManager saves cfg only when the user presses s; any other exit
// returns ErrCancelled and leaves cfg unchanged.
func RunToolManager(cfg *config.ManagerConfig) error {
	m, err := tea.NewProgram(initialToolModel(cfg)).Run()
	if err != nil {
		return err
	}
	final := m.(toolModel)
	if !final.saved || final.cancelled {
		return ErrCancelled
	}
	cfg.Services = final.services
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}
	fmt.Println("\n⚠️  Changes saved to config.yaml. Restart services to apply env/tool settings.")
	return nil
}
