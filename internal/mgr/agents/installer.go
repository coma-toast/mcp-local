package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/coma-toast/mcp-local/internal/mgr/config"
)

// Delegation to `ast-mcp install|uninstall` for ast-context-cache services.
//
// Contract of the ast CLI (v4+):
//
//	<command> install|uninstall --target <opencode|cursor|claude_desktop> --component mcp --yes --json --mcp-url <url>
//	exit 0 ok, 2 confirmation required, 3 conflict/parse error, 4 unsupported
//	--json: {"changes":[{"path","kind","diff","skipped","reason"}],"status":[{"target","component","status","path"}],"warnings":[...]}

const MinASTMCPMajor = 4

var (
	installerTimeout = 60 * time.Second
	versionTimeout   = 10 * time.Second
	semverRe         = regexp.MustCompile(`v?(\d+)\.(\d+)\.(\d+)`)
)

// astTargets maps mcp-local host names to ast-mcp install targets.
var astTargets = map[string]string{
	HostOpenCode: "opencode",
	HostCursor:   "cursor",
	HostClaude:   "claude_desktop",
}

type installChange struct {
	Path    string `json:"path"`
	Kind    string `json:"kind"`
	Diff    string `json:"diff"`
	Skipped bool   `json:"skipped"`
	Reason  string `json:"reason"`
}

type installStatus struct {
	Target    string `json:"target"`
	Component string `json:"component"`
	Status    string `json:"status"`
	Path      string `json:"path"`
}

type installResult struct {
	Changes  []installChange `json:"changes"`
	Status   []installStatus `json:"status"`
	Warnings []string        `json:"warnings"`
}

// installers caches `--version` checks per command for one register/deregister pass.
type installers struct {
	versions map[string]error
}

func newInstallers() *installers {
	return &installers{versions: map[string]error{}}
}

// delegate runs the service's installer for host when the service has one and it is usable.
// handled=false means the caller must use mcp-local's native writer.
func (in *installers) delegate(s config.ServiceConfig, host, verb string) (handled, found bool, err error) {
	if s.EffectiveInstaller() != config.InstallerASTMCP {
		return false, false, nil
	}
	target, ok := astTargets[host]
	if !ok {
		return false, false, nil
	}
	url := s.MCPEndpoint()
	if url == "" {
		return false, false, nil
	}
	cmd := config.ExpandPath(strings.TrimSpace(s.Command))
	if cmd == "" {
		cmd = config.InstallerASTMCP
	}
	if err := in.checkVersion(cmd); err != nil {
		fmt.Fprintf(Warn, "  ⚠️  %s: %v; using mcp-local's %s writer\n", s.Name, err, host)
		return false, false, nil
	}
	res, code, stderr, err := runInstaller(cmd, verb, target, url)
	for _, w := range res.Warnings {
		fmt.Fprintf(Warn, "  ⚠️  %s (%s): %s\n", s.Name, host, w)
	}
	switch {
	case err != nil:
		return true, false, fmt.Errorf("%s %s: %w", cmd, verb, err)
	case code == 0:
		return true, res.changed(), nil
	case code == 4:
		fmt.Fprintf(Warn, "  ⚠️  %s: %s %s does not support target %s; using mcp-local's writer\n", s.Name, cmd, verb, target)
		return false, false, nil
	}
	return true, false, fmt.Errorf("%s %s --target %s exited %d (%s)%s", cmd, verb, target, code, exitMeaning(code), res.detail(stderr))
}

func (r installResult) changed() bool {
	for _, c := range r.Changes {
		if !c.Skipped {
			return true
		}
	}
	return false
}

func (r installResult) detail(stderr string) string {
	var parts []string
	for _, c := range r.Changes {
		if c.Reason != "" {
			parts = append(parts, fmt.Sprintf("%s: %s", c.Path, c.Reason))
		}
	}
	parts = append(parts, r.Warnings...)
	if s := strings.TrimSpace(stderr); s != "" {
		parts = append(parts, s)
	}
	if len(parts) == 0 {
		return ""
	}
	return ": " + strings.Join(parts, "; ")
}

func exitMeaning(code int) string {
	switch code {
	case 2:
		return "confirmation required"
	case 3:
		return "conflict or parse error"
	case 4:
		return "unsupported"
	}
	return "failed"
}

func (in *installers) checkVersion(cmd string) error {
	if err, ok := in.versions[cmd]; ok {
		return err
	}
	err := checkVersion(cmd)
	in.versions[cmd] = err
	return err
}

func checkVersion(cmd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), versionTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, cmd, "--version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s --version failed: %v", cmd, err)
	}
	major, ver, ok := parseVersion(string(out))
	if !ok {
		return fmt.Errorf("%s --version: no version in %q", cmd, strings.TrimSpace(string(out)))
	}
	if major < MinASTMCPMajor {
		return fmt.Errorf("%s is %s (< %d.0.0, no install command)", cmd, ver, MinASTMCPMajor)
	}
	return nil
}

func parseVersion(s string) (int, string, bool) {
	m := semverRe.FindStringSubmatch(s)
	if m == nil {
		return 0, "", false
	}
	major, _ := strconv.Atoi(m[1])
	return major, strings.TrimPrefix(m[0], "v"), true
}

func runInstaller(cmd, verb, target, url string) (installResult, int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), installerTimeout)
	defer cancel()
	c := exec.CommandContext(ctx, cmd, verb, "--target", target, "--component", "mcp", "--yes", "--json", "--mcp-url", url)
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	var res installResult
	_ = json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &res)
	if exitErr, ok := err.(*exec.ExitError); ok {
		return res, exitErr.ExitCode(), stderr.String(), nil
	}
	return res, 0, stderr.String(), err
}
