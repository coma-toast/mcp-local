# mcp-local

`mcp-local` is a unified local control plane for managing Model Context Protocol (MCP) servers. Instead of managing multiple binaries, ports, and agent configurations manually, `mcp-local` provides a single point of control for your local AI infrastructure.

## Why mcp-local?

When building an advanced AI agent setup, you often end up with a fragmented ecosystem:
- Multiple MCP servers running on different ports.
- Complex environment variables (like `ONNXRUNTIME_LIB`) needed for specific binaries.
- Manual updates to `opencode.jsonc` or `mcp.json` every time a server changes.
- Difficulty syncing your "local AI stack" across multiple machines.

`mcp-local` solves this by treating your MCP infrastructure as code. You define your servers in a single portable config, and `mcp-local` handles the rest.

**AI agents:** See [AGENTS.md](AGENTS.md) for command workflows, paths, ast-context-cache tool tiers, and what not to edit manually.

## Key Features

### Unified Lifecycle Management
- `mcp-local start <service>` or `mcp-local start --all` — launches and registers with OpenCode, Cursor, and Claude Desktop (configurable).
- `mcp-local stop <service>` / `stop --all`, `restart <service>` / `restart --all`.
- `mcp-local stop --deregister` (default) — removes agent entries when stopping.
- `mcp-local rebuild <service>` — runs `build_command` and re-registers.
- `mcp-local status` — bubbletea live table; `mcp-local status --plain` for scripting.
- stdio services detected via PID file (not just port).

### Agent registration
On `start` / `restart` / `register`, services are registered with enabled agents:
- **OpenCode**: `~/.config/opencode/opencode.jsonc` — HTTP as `remote`, stdio as `local`.
- **Cursor**: `~/.cursor/mcp.json` — HTTP as `{ "url": "..." }`, stdio as `{ "command": "...", "args": [...], "env": {...} }`.
- **Claude Desktop**: `claude_desktop_config.json` — same shape as Cursor.

```yaml
agents:
  opencode: true
  cursor: true
  claude: true
```

### Tool & Tier Management (ast-context-cache)
- `mcp-local tools sync <service>` — fetch tool list from a running MCP server via `tools/list`.
- `mcp-local tools` — TUI to enable/disable tools, change tiers, edit descriptions.
- `mcp-local tools apply <service>` — write overrides to `~/.astcache/tools.json` (or `tools_config_path`).
- `mcp-local json tools <service>` — print the tools.json overrides that would be written.
- On start: `AST_MCP_TIER`, optional `AST_MCP_CODE_MODE`, and `AST_MCP_TOOLS_CONFIG` when tools are configured.
- **Restart** the service after changing tool settings (ast-mcp reads `tools.json` at startup only).

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
mcp-local tools apply ast-context-cache   # Write tools.json overrides
mcp-local json tools ast-context-cache  # Preview tools.json overrides
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
| `no_code_mode` | Injected as `AST_MCP_CODE_MODE=false` when true |
| `tools_config_path` | Path to tools overrides JSON (default `~/.astcache/tools.json`) |
| `tools` | Per-tool overrides (name, enabled, tier, description) → `tools.json` on start |

## Documentation

| Doc | Audience |
|-----|----------|
| [AGENTS.md](AGENTS.md) | AI coding assistants — commands, paths, workflows, limitations |
| [docs/goal-alignment-plan.md](docs/goal-alignment-plan.md) | Implementation history and future catalog TUI |
| [examples/](examples/) | Starter `config.yaml` and service snippets |

## Development

```bash
make        # build
make test   # run tests
make check  # test + vet + fmt
```

CI runs `go test ./...`, `go vet`, and `go build` on push/PR.

## License
MIT
