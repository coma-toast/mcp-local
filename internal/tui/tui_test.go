package tui

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

var (
	keyEnter = tea.KeyMsg{Type: tea.KeyEnter}
	keyEsc   = tea.KeyMsg{Type: tea.KeyEsc}
	keyCtrlC = tea.KeyMsg{Type: tea.KeyCtrlC}
	keyUp    = tea.KeyMsg{Type: tea.KeyUp}
	keyDown  = tea.KeyMsg{Type: tea.KeyDown}
	keyTab   = tea.KeyMsg{Type: tea.KeyTab}
	keyLeft  = tea.KeyMsg{Type: tea.KeyLeft}
	keyRight = tea.KeyMsg{Type: tea.KeyRight}
	keyBack  = tea.KeyMsg{Type: tea.KeyBackspace}
)

func runes(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

// send feeds msgs through Update (rendering View after each) and returns the
// final model and the last command.
func send[M tea.Model](t *testing.T, m M, msgs ...tea.Msg) (M, tea.Cmd) {
	t.Helper()
	var cmd tea.Cmd
	for _, msg := range msgs {
		var next tea.Model
		next, cmd = m.Update(msg)
		m = next.(M)
		_ = m.View()
	}
	return m, cmd
}

func typeText(s string) []tea.Msg {
	var msgs []tea.Msg
	for _, r := range s {
		msgs = append(msgs, runes(string(r)))
	}
	return msgs
}

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

// --- tools TUI ---

func TestToolsNoPanicWithoutServicesOrTools(t *testing.T) {
	keys := []tea.Msg{keyDown, keyUp, keyEnter, keyLeft, keyRight, runes("j"), runes("k"), runes("h"), runes("l"), keyEnter, keyBack, keyEnter, keyLeft}
	for name, cfg := range map[string]*config.ManagerConfig{
		"no services": {},
		"no tools":    {Services: []config.ServiceConfig{{Name: "svc"}}},
	} {
		t.Run(name, func(t *testing.T) {
			m, _ := send(t, initialToolModel(cfg), keys...)
			if m.cancelled || m.saved {
				t.Fatalf("navigation keys must not quit: %+v", m)
			}
		})
	}
}

func TestToolsQuitKeysDiscard(t *testing.T) {
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{{Name: "svc", Tools: []config.ToolConfig{{Name: "t1", Enabled: true}}}}}
	for name, quit := range map[string]tea.Msg{"esc": keyEsc, "q": runes("q"), "ctrl+c": keyCtrlC} {
		t.Run(name, func(t *testing.T) {
			m, cmd := send(t, initialToolModel(cfg), keyEnter, keyLeft, quit)
			if !isQuit(cmd) || !m.cancelled || m.saved {
				t.Fatalf("want cancelled quit, got cancelled=%v saved=%v", m.cancelled, m.saved)
			}
			if !cfg.Services[0].Tools[0].Enabled {
				t.Fatal("toggle leaked into the caller's config before save")
			}
		})
	}
}

func TestToolsSaveKey(t *testing.T) {
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{{Name: "svc", Tools: []config.ToolConfig{{Name: "t1", Enabled: true}}}}}
	m, cmd := send(t, initialToolModel(cfg), keyEnter, keyLeft, runes("s"))
	if !isQuit(cmd) || !m.saved || m.cancelled {
		t.Fatalf("want saved quit, got saved=%v cancelled=%v", m.saved, m.cancelled)
	}
	if m.services[0].Tools[0].Enabled {
		t.Fatal("toggle not applied to working copy")
	}
}

func TestToolsDescEditTypesQAndS(t *testing.T) {
	cfg := &config.ManagerConfig{Services: []config.ServiceConfig{{Name: "svc", Tools: []config.ToolConfig{{Name: "t1"}}}}}
	msgs := append([]tea.Msg{keyEnter, keyEnter}, typeText("qs")...)
	m, _ := send(t, initialToolModel(cfg), append(msgs, keyEnter)...)
	if m.cancelled || m.saved || m.services[0].Tools[0].Description != "qs" {
		t.Fatalf("q/s must be text in the description field: %+v", m.services[0].Tools[0])
	}
}

// --- wizard ---

func wizardWith(names ...string) model {
	var svcs []config.ServiceConfig
	for _, n := range names {
		svcs = append(svcs, config.ServiceConfig{Name: n})
	}
	return initialModel(svcs)
}

func TestWizardCursorBounded(t *testing.T) {
	var downs []tea.Msg
	for range 10 {
		downs = append(downs, keyDown)
	}
	m, _ := send(t, wizardWith("a"), downs...)
	if m.cursor != m.lastRow() {
		t.Fatalf("cursor = %d, want %d", m.cursor, m.lastRow())
	}
	m, _ = send(t, m, keyEnter)
	if m.state != confirmState {
		t.Fatalf("Enter on last row: state = %v, want confirm", m.state)
	}
}

func TestWizardEscCancelsEveryScreen(t *testing.T) {
	screens := map[string][]tea.Msg{
		"list":    nil,
		"form":    {keyDown, keyEnter},
		"delete":  {runes("d")},
		"confirm": {keyDown, keyDown, keyEnter},
	}
	existing := &config.ManagerConfig{Services: []config.ServiceConfig{{Name: "a"}}}
	for name, prefix := range screens {
		for _, quit := range []tea.Msg{keyEsc, keyCtrlC} {
			t.Run(name, func(t *testing.T) {
				m, cmd := send(t, wizardWith("a"), append(prefix, quit)...)
				if !isQuit(cmd) {
					t.Fatal("want quit")
				}
				if cfg, err := wizardResult(existing, m); !errors.Is(err, ErrCancelled) || cfg != nil {
					t.Fatalf("got cfg=%v err=%v, want ErrCancelled", cfg, err)
				}
			})
		}
	}
}

func TestWizardEnterOnRowEditsNotDeletes(t *testing.T) {
	m, _ := send(t, wizardWith("a", "b"), keyEnter)
	if len(m.services) != 2 || m.state != formState || m.editIdx != 0 || m.inputs[0].Value() != "a" {
		t.Fatalf("Enter on a row must open the edit form: state=%v services=%d", m.state, len(m.services))
	}
	m, _ = send(t, m, append(typeText("2"), keyEnter)...)
	if len(m.services) != 2 || m.services[0].Name != "a2" {
		t.Fatalf("edit not applied: %+v", m.services)
	}
}

func TestWizardDeleteNeedsConfirmation(t *testing.T) {
	m, _ := send(t, wizardWith("a", "b"), keyDown, runes("d"))
	if m.state != deleteState || len(m.services) != 2 {
		t.Fatal("d must ask for confirmation first")
	}
	m, _ = send(t, m, runes("n"))
	if m.state != listState || len(m.services) != 2 {
		t.Fatal("n must keep the service")
	}
	m, _ = send(t, m, runes("d"), runes("y"))
	if len(m.services) != 1 || m.services[0].Name != "a" || m.cursor > m.lastRow() {
		t.Fatalf("y must delete b: %+v cursor=%d", m.services, m.cursor)
	}
}

func addService(t *testing.T, m model, fields ...string) model {
	t.Helper()
	for m.cursor < len(m.services) {
		m, _ = send(t, m, keyDown)
	}
	m, _ = send(t, m, keyEnter)
	for i, f := range fields {
		if i > 0 {
			m, _ = send(t, m, keyTab)
		}
		m, _ = send(t, m, typeText(f)...)
	}
	m, _ = send(t, m, keyEnter)
	return m
}

func TestWizardEnvValueWithEquals(t *testing.T) {
	m := addService(t, initialModel(nil), "svc", "cmd", "", "stdio", "URL=http://x?a=b,K=v")
	if m.state != listState || len(m.services) != 1 {
		t.Fatalf("service not added: %s", m.errMsg)
	}
	if got := m.services[0].Env["URL"]; got != "http://x?a=b" {
		t.Fatalf("URL = %q", got)
	}
	if got := m.services[0].Env["K"]; got != "v" {
		t.Fatalf("K = %q", got)
	}
}

func TestWizardInvalidPort(t *testing.T) {
	for _, port := range []string{"abc", "0", "70000", "-1"} {
		m := addService(t, initialModel(nil), "svc", "cmd", port)
		if m.state != formState || len(m.services) != 0 || m.errMsg == "" {
			t.Fatalf("port %q: want validation error, got state=%v services=%+v", port, m.state, m.services)
		}
	}
	m := addService(t, initialModel(nil), "svc", "cmd", "8080")
	if len(m.services) != 1 || m.services[0].Port != 8080 {
		t.Fatalf("valid port not saved: %+v", m.services)
	}
}

func TestWizardNewServiceUsesBuildCommand(t *testing.T) {
	m := addService(t, initialModel(nil), "svc", "cmd", "", "", "", "make")
	if s := m.services[0]; s.BuildCommand != "make" || s.BuildCmd != "" {
		t.Fatalf("build fields = %q/%q", s.BuildCmd, s.BuildCommand)
	}
}

func TestWizardResultPreservesAgents(t *testing.T) {
	existing := &config.ManagerConfig{Agents: &config.AgentsConfig{Cursor: true, Claude: true}, Services: []config.ServiceConfig{{Name: "a"}}}
	m, _ := send(t, initialModel(existing.Services), keyDown, keyDown, keyEnter, keyEnter)
	cfg, err := wizardResult(existing, m)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Agents == nil || *cfg.Agents != *existing.Agents || len(cfg.Services) != 1 {
		t.Fatalf("got %+v", cfg)
	}
}

// --- editor ---

func TestEditorJKAreText(t *testing.T) {
	m := initialEditModel(&config.ServiceConfig{Name: "svc", Log: "log"})
	for range 8 {
		m, _ = send(t, m, keyDown)
	}
	m, _ = send(t, m, typeText("jkmake")...)
	if m.cursor != 8 || m.inputs[8].Value() != "jkmake" || m.inputs[7].Value() != "log" {
		t.Fatalf("cursor=%d build=%q log=%q", m.cursor, m.inputs[8].Value(), m.inputs[7].Value())
	}
	m, _ = send(t, m, keyTab, tea.KeyMsg{Type: tea.KeyShiftTab}, keyUp)
	if m.cursor != 7 {
		t.Fatalf("tab/shift+tab/up: cursor = %d, want 7", m.cursor)
	}
}

func TestEditorEscCancels(t *testing.T) {
	for _, quit := range []tea.Msg{keyEsc, keyCtrlC} {
		m, cmd := send(t, initialEditModel(&config.ServiceConfig{Name: "svc"}), append(typeText("x"), quit)...)
		if !isQuit(cmd) || !m.cancelled || m.saved {
			t.Fatalf("want cancelled, got cancelled=%v saved=%v", m.cancelled, m.saved)
		}
	}
}

func TestEditorInvalidPort(t *testing.T) {
	m := initialEditModel(&config.ServiceConfig{Name: "svc", Port: 80})
	m, _ = send(t, m, keyDown, keyDown, keyBack, keyBack, runes("x"), keyEnter)
	if m.saved || m.errMsg == "" {
		t.Fatal("bad port must not save")
	}
}

func TestEditorBuildField(t *testing.T) {
	cases := []struct {
		name               string
		svc                config.ServiceConfig
		shown              string
		wantCmd, wantCmmnd string
	}{
		{"build_cmd in use", config.ServiceConfig{Name: "s", BuildCmd: "old", BuildCommand: "stale"}, "old", "new", "stale"},
		{"build_command in use", config.ServiceConfig{Name: "s", BuildCommand: "old"}, "old", "", "new"},
		{"new service", config.ServiceConfig{Name: "s"}, "", "", "new"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := initialEditModel(&tc.svc)
			if got := m.inputs[8].Value(); got != tc.shown {
				t.Fatalf("shown %q, want %q", got, tc.shown)
			}
			m.inputs[8].SetValue("new")
			m, cmd := send(t, m, keyEnter)
			if !isQuit(cmd) || !m.saved {
				t.Fatalf("not saved: %s", m.errMsg)
			}
			if m.svc.BuildCmd != tc.wantCmd || m.svc.BuildCommand != tc.wantCmmnd || config.EffectiveBuild(m.svc) != "new" {
				t.Fatalf("got %q/%q", m.svc.BuildCmd, m.svc.BuildCommand)
			}
		})
	}
}

func TestSetBuildClearKeepsEffectiveEmpty(t *testing.T) {
	s := config.ServiceConfig{BuildCmd: "a", BuildCommand: "b"}
	setBuild(&s, "")
	if got := config.EffectiveBuild(s); got != "" {
		t.Fatalf("clearing the field left effective build %q", got)
	}
}

func TestEditorKeepsUnshownFields(t *testing.T) {
	svc := config.ServiceConfig{Name: "svc", Args: []string{"--x"}, Env: map[string]string{"A": "1"}, Port: 80}
	m, _ := send(t, initialEditModel(&svc), keyEnter)
	if !m.saved || len(m.svc.Args) != 1 || m.svc.Env["A"] != "1" || m.svc.Port != 80 {
		t.Fatalf("got %+v", m.svc)
	}
}
