# mcp-local

`mcp-local` is a unified local control plane for managing Model Context Protocol (MCP) servers. Instead of managing multiple binaries, ports, and agent configurations manually, `mcp-local` provides a single point of control for your local AI infrastructure.

## Why mcp-local?

When building an advanced AI agent setup, you often end up with a fragmented ecosystem:
- Multiple MCP servers running on different ports.
- Complex environment variables (like `ONNXRUNTIME_LIB`) needed for specific binaries.
- Manual updates to `opencode.jsonc` or `mcp.json` every time a server changes.
- Difficulty syncing your "local AI stack" across multiple machines.

`mcp-local` solves this by treating your MCP infrastructure as code. You define your servers in a single portable config, and `mcp-local` handles the rest.

## Key Features

### Unified Lifecycle Management
- `mcp-local start <service>` or `mcp-local start --all` — launches and registers with OpenCode and Cursor.
- `mcp-local stop <service>` / `stop --all`, `restart <service>` / `restart --all`.
- `mcp-local stop --deregister` (default) — removes agent entries when stopping.
- `mcp-local rebuild <service>` — runs `build_command` and re-registers.
- `mcp-local status` — bubbletea live table; `mcp-local status --plain` for scripting.
- stdio services detected via PID file (not just port).

### Dual Agent Registration
On `start` / `restart` / `register`, services are registered with both:
- **OpenCode**: `~/.config/opencode/opencode.jsonc` — HTTP as `remote`, stdio as `local`.
- **Cursor**: `~/.cursor/mcp.json` — HTTP as `{ "url": "..." }`, stdio as `{ "command": "...", "args": [...], "env": {...} }`.

Control which agents receive registrations via top-level config:
```yaml
agents:
  opencode: true
  cursor: true
```

### Tool & Tier Management
- `mcp-local tools` — TUI to enable/disable tools, change tiers, edit descriptions.
- `mcp-local tools sync <service>` — fetches tool list from a running MCP server via `tools/list` and populates config.
- On start, `active_tier` and `code_mode` are injected as `AST_MCP_TIER` / `AST_MCP_CODE_MODE` env vars.
- Disabled tools are passed as `AST_MCP_DISABLED_TOOLS=tool1,tool2` env var.
- Restart services after changing tool settings to apply.

### Portable Configuration
All settings in `~/.mcp-local/config.yaml`. Sync across machines via Dropbox, iCloud, or Git.

## Supported Agents
- **OpenCode**: Full integration with `opencode.jsonc` / `opencode.json`.
- **Cursor**: Full integration with `~/.cursor/mcp.json`.
- **Any MCP-compliant client**: Manages standard HTTP or stdio MCP servers.

## Supported MCP Servers
While it can manage *any* binary, it is optimized for:
- **ast-context-cache**: Full lifecycle, tool tier env injection, and TUI management.
- **model-proxy**: stdio process management and registration.
- **Custom Servers**: Add any binary to `config.yaml`.

## Getting Started

### Installation
```bash
git clone https://github.com/coma-toast/mcp-local.git
cd mcp-local
make          # writes ./bin/mcp-local
make test     # run unit tests
# Optional: install system-wide
# sudo make install
```

Add `./bin` to your PATH.

### Configuration
Create `~/.mcp-local/config.yaml` or copy the starter:
```bash
mkdir -p ~/.mcp-local
cp examples/config.starter.yaml ~/.mcp-local/config.yaml
```

Example config:
```yaml
agents:
  opencode: true
  cursor: true

services:
  - name: ast-context-cache
    command: ${HOME}/git/ast-context-cache/ast-mcp
    port: 7821
    type: http
    mcp_url: http://localhost:7821/mcp
    health_url: http://localhost:7821/health
    dashboard_url: http://localhost:7830
    log: ${HOME}/.mcp-local/ast-context-cache.log
    active_tier: complete
    env:
      ONNXRUNTIME_LIB: /opt/homebrew/lib/libonnxruntime.dylib
    build_command: "cd ${HOME}/git/ast-context-cache && make build"
    no_build_on_start: false   # set true to skip build on start
    deps:
      - ${HOME}/git/ast-context-cache/model/model.onnx

  - name: model-proxy
    command: ${HOME}/git/model-proxy/bin/model-proxy
    type: stdio
    args:
      - --config
      - ${HOME}/.config/model-proxy/config.yaml
```

### Usage
```bash
mcp-local start --all          # Start and register with all agents
mcp-local status               # Live table (q to quit)
mcp-local status --plain       # One-shot text output
mcp-local tools sync ast-context-cache   # Fetch tools from running server
mcp-local tools                # Tool tier TUI
mcp-local stop --all           # Stop and deregister (default)
mcp-local stop --all --deregister=false  # Stop without deregistering
mcp-local register             # Register all services with agents
mcp-local deregister           # Remove entries for stopped services
```

## Config Reference

| Field | Description |
|-------|-------------|
| `name` | Service identifier |
| `command` | Path to binary (supports `${HOME}` and `~/`) |
| `args` | Command-line arguments |
| `port` | TCP listen port (omit for stdio) |
| `type` | `http` or `stdio` |
| `mcp_url` | MCP endpoint URL for remote registration |
| `health_url` | HTTP health check endpoint |
| `dashboard_url` | Web dashboard URL |
| `log` | Log file path |
| `env` | Environment variables (key-value map) |
| `build_command` | Shell command to build the service |
| `no_build_on_start` | Skip build on `start` (default: false, i.e. build runs) |
| `deps` | Required file paths (checked before start) |
| `active_tier` | Injected as `AST_MCP_TIER` env var |
| `code_mode` | Injected as `AST_MCP_CODE_MODE=true` when true |
| `tools` | Per-tool config (name, enabled, tier, description) |

## Development

```bash
make        # build
make test   # run tests
make check  # test + vet + fmt
```

CI runs `go test ./...`, `go vet`, and `go build` on push/PR.

## License
MIT
