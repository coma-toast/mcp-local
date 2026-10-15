package agents

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/coma-toast/mcp-local/internal/mgr/claudedesktop"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/cursor"
	"github.com/coma-toast/mcp-local/internal/mgr/jsonagent"
	"github.com/coma-toast/mcp-local/internal/mgr/opencode"
)

const (
	HostOpenCode = "opencode"
	HostCursor   = "cursor"
	HostClaude   = "claude"
)

// Warn receives non-fatal warnings (installer fallbacks, installer warnings).
var Warn io.Writer = os.Stderr

type Targets struct {
	OpenCode bool
	Cursor   bool
	Claude   bool
}

func DefaultTargets() Targets {
	return Targets{OpenCode: true, Cursor: true, Claude: true}
}

// TargetsFromConfig is the `agents:` block as written (all false = no hosts); without an
// `agents:` key every host is enabled (legacy behavior).
func TargetsFromConfig(cfg config.ManagerConfig) Targets {
	if cfg.Agents == nil {
		return DefaultTargets()
	}
	return Targets{OpenCode: cfg.Agents.OpenCode, Cursor: cfg.Agents.Cursor, Claude: cfg.Agents.Claude}
}

// Hosts lists the enabled host names in a stable order.
func (t Targets) Hosts() []string {
	var hosts []string
	if t.OpenCode {
		hosts = append(hosts, HostOpenCode)
	}
	if t.Cursor {
		hosts = append(hosts, HostCursor)
	}
	if t.Claude {
		hosts = append(hosts, HostClaude)
	}
	return hosts
}

// RegisterAll registers services with every enabled host. A failure on one host or service
// does not stop the others; all failures are returned joined.
func RegisterAll(services []config.ServiceConfig, targets Targets) error {
	inst := newInstallers()
	var errs []error
	for _, host := range targets.Hosts() {
		var native []config.ServiceConfig
		for _, s := range services {
			handled, _, err := inst.delegate(s, host, "install")
			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %s: %w", host, s.Name, err))
				continue
			}
			if !handled {
				native = append(native, s)
			}
		}
		if len(native) == 0 {
			continue
		}
		if err := registerNative(host, native); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", host, err))
		}
	}
	return errors.Join(errs...)
}

// DeregisterAll removes svc from every enabled host (the service is named explicitly, so
// entries mcp-local didn't write are removed too), continuing past per-host errors.
func DeregisterAll(svc config.ServiceConfig, targets Targets) error {
	inst := newInstallers()
	var errs []error
	for _, host := range targets.Hosts() {
		if _, err := inst.deregister(svc, host, true); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", host, err))
		}
	}
	return errors.Join(errs...)
}

// Removal is what Deregister did on one host.
type Removal int

const (
	NotFound Removal = iota // no entry
	Removed                 // entry removed
	NotOwned                // entry left alone: mcp-local didn't write it
)

// Deregister removes svc from one host. Unless explicit, an entry that mcp-local's own writer
// didn't record writing is left in place (NotOwned), since it may be the user's.
func Deregister(svc config.ServiceConfig, host string, explicit bool) (Removal, error) {
	return newInstallers().deregister(svc, host, explicit)
}

func (in *installers) deregister(svc config.ServiceConfig, host string, explicit bool) (Removal, error) {
	handled, found, err := in.delegate(svc, host, "uninstall")
	if err != nil || handled {
		return removal(found), err
	}
	a, err := hostAgent(host)
	if err != nil {
		return NotFound, err
	}
	if !explicit && !a.Owns(svc.Name) {
		if a.Has(svc.Name) {
			return NotOwned, nil
		}
		return NotFound, nil
	}
	found, err = a.Deregister(svc.Name)
	return removal(found), err
}

func removal(found bool) Removal {
	if found {
		return Removed
	}
	return NotFound
}

func hostAgent(host string) (jsonagent.Agent, error) {
	switch host {
	case HostOpenCode:
		return opencode.Agent(), nil
	case HostCursor:
		return cursor.Agent(), nil
	case HostClaude:
		return claudedesktop.Agent(), nil
	}
	return jsonagent.Agent{}, fmt.Errorf("unknown agent host %q", host)
}

func registerNative(host string, services []config.ServiceConfig) error {
	a, err := hostAgent(host)
	if err != nil {
		return err
	}
	return a.RegisterServices(services)
}
