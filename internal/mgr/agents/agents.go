package agents

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/coma-toast/mcp-local/internal/mgr/claudedesktop"
	"github.com/coma-toast/mcp-local/internal/mgr/config"
	"github.com/coma-toast/mcp-local/internal/mgr/cursor"
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

func TargetsFromConfig(cfg config.ManagerConfig) Targets {
	if cfg.Agents.OpenCode || cfg.Agents.Cursor || cfg.Agents.Claude {
		return Targets{
			OpenCode: cfg.Agents.OpenCode,
			Cursor:   cfg.Agents.Cursor,
			Claude:   cfg.Agents.Claude,
		}
	}
	return DefaultTargets()
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

// DeregisterAll removes svc from every enabled host, continuing past per-host errors.
func DeregisterAll(svc config.ServiceConfig, targets Targets) error {
	inst := newInstallers()
	var errs []error
	for _, host := range targets.Hosts() {
		if _, err := inst.deregister(svc, host); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", host, err))
		}
	}
	return errors.Join(errs...)
}

// Deregister removes svc from one host; found reports whether an entry was removed.
func Deregister(svc config.ServiceConfig, host string) (bool, error) {
	return newInstallers().deregister(svc, host)
}

func (in *installers) deregister(svc config.ServiceConfig, host string) (bool, error) {
	handled, found, err := in.delegate(svc, host, "uninstall")
	if err != nil || handled {
		return found, err
	}
	return deregisterNative(host, svc.Name)
}

func registerNative(host string, services []config.ServiceConfig) error {
	switch host {
	case HostOpenCode:
		return opencode.RegisterServices(services)
	case HostCursor:
		return cursor.RegisterServices(services)
	case HostClaude:
		return claudedesktop.RegisterServices(services)
	}
	return fmt.Errorf("unknown agent host %q", host)
}

func deregisterNative(host, name string) (bool, error) {
	switch host {
	case HostOpenCode:
		return opencode.Deregister(name)
	case HostCursor:
		return cursor.Deregister(name)
	case HostClaude:
		return claudedesktop.Deregister(name)
	}
	return false, fmt.Errorf("unknown agent host %q", host)
}
