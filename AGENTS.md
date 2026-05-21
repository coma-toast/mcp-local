# AGENTS.md — mcp-local

Guidance for AI coding assistants operating on or through **mcp-local**, the local MCP control plane.

## What this tool is

**mcp-local** supervises local MCP server processes and keeps agent editor configs in sync. It is **not** an MCP server and does **not** implement MCP JSON-RPC (`tools/list`, etc.)—those belong to the servers it starts (e.g. ast-context-cache).

| Responsibility | mcp-local | MCP server (e.g. ast-mcp) |
|----------------|-----------|---------------------------|
| Start/stop/restart binaries | Yes | N/A |
| Register `mcp.json` / `opencode.jsonc` | Yes | N/A |
| Tool tier policy file (`tools.json`) | Writes from config | Reads at startup |
| Serve MCP over HTTP/stdio | No | Yes |

## Paths (defaults)

| Path | Purpose |
|------|---------|
| `~/.mcp-local/config.yaml` | Source of truth for services and tool overrides |
| `~/.mcp-local/pids/<name>.pid` | stdio / process tracking |
| `~/.astcache/tools.json` | Per-tool `enabled` / `tier` / `description` for ast-context-cache |
| `~/.cursor/mcp.json` | Cursor MCP entries (`mcpServers`) |
| `~/.config/opencode/opencode.jsonc` (or `.json`) | OpenCode MCP entries (`mcp`) |
| Claude Desktop config | OS-specific; see `internal/mgr/claudedesktop` |

Build output: `./bin/mcp-local` after `make` in this repo.

## When to use mcp-local (agent workflow)

**Prefer mcp-local** when the user’s stack is already defined in `config.yaml` or they want lifecycle + registration in one place.

**Do not** hand-edit `~/.cursor/mcp.json` or OpenCode MCP blocks for services that mcp-local manages—changes will be overwritten on the next `start`, `restart`, or `register`.

**Do** use shell commands below; use `mcp-local status --plain` and `mcp-local validate` before assuming a service is up.

### First-time / missing config

If `~/.mcp-local/config.yaml` is missing, many commands launch the **config wizard** TUI. For non-interactive setup, copy [examples/config.starter.yaml](examples/config.starter.yaml) and edit paths (`command`, `env`, `deps`).

```bash
mkdir -p ~/.mcp-local
cp examples/config.starter.yaml ~/.mcp-local/config.yaml
mcp-local validate
```

## Command reference (by task)

### Lifecycle

```bash
mcp-local list
mcp-local start <service> | mcp-local start --all
mcp-local stop <service> | mcp-local stop --all
mcp-local restart <service> | mcp-local restart --all
mcp-local rebuild <service>          # build_command + re-register
```

- **`stop --deregister`** defaults to **true** (removes agent entries). Use `--deregister=false` to leave MCP JSON in place.
- **`no_build_on_start: true`** on a service skips `build_command` on `start`.
- HTTP services: running = port in use. stdio services: running = PID file under `~/.mcp-local/pids/`.

### Agent registration (without restart)

```bash
mcp-local register [service]       # --dry-run to preview
mcp-local deregister [service]     # skips running services unless name given
mcp-local registered [service]     # ✅/❌ per agent
mcp-local import [--source opencode|cursor|claude]  # pull entries into config.yaml
```

Enable agents in config:

```yaml
agents:
  opencode: true
  cursor: true
  claude: true
```

If all three are `false`, mcp-local still defaults to registering OpenCode + Cursor + Claude (legacy behavior).

### ast-context-cache tool tiers

ast-mcp reads **`~/.astcache/tools.json`** and **`AST_MCP_TIER`** only at **process start**. After changing tools in config or TUI, **restart** the service.

```bash
mcp-local start ast-context-cache          # must be running for sync
mcp-local tools sync ast-context-cache     # MCP tools/list → config.yaml
mcp-local tools                            # TUI: enable/tier/description
mcp-local tools apply ast-context-cache    # write tools.json (no restart)
mcp-local restart ast-context-cache        # load tools.json + env
mcp-local json tools ast-context-cache     # preview overrides JSON
```

Config fields (see [README.md](README.md#config-reference)):

- `active_tier` → `AST_MCP_TIER` (`core` | `extended` | `complete`)
- `code_mode` / `no_code_mode` → `AST_MCP_CODE_MODE`
- `tools_config_path` → `AST_MCP_TOOLS_CONFIG` (non-default path)
- `tools:` list → written to tools.json on start (via `asttools.ApplyStartEnv`)

**Do not** set `AST_MCP_DISABLED_TOOLS`—ast-context-cache does not read it.

### Configuration editing

```bash
mcp-local config                 # $EDITOR on config.yaml
mcp-local add <name> [flags]     # CLI add/update
mcp-local edit <name>            # service TUI
mcp-local remove <name>          # --deregister default true
mcp-local validate
```

### Diagnostics

```bash
mcp-local status                 # live TUI (q to quit)
mcp-local status --plain         # scriptable one-shot
mcp-local health <service>
mcp-local log <service>          # tail one log
mcp-local logs                   # tail all (Ctrl+C)
mcp-local open <service>         # dashboard_url in browser
```

## Common agent scenarios

### User says ast MCP tools are missing

1. `mcp-local status --plain` — is `ast-context-cache` running?
2. Check `active_tier` in config; suggest `complete` if they need `execute_code`.
3. Check `~/.astcache/tools.json` for `"enabled": false`.
4. `mcp-local restart ast-context-cache` after any tier/tools change.

### User wants Cursor to see the server

1. Ensure `agents.cursor: true` and service has `mcp_url` (HTTP) or `command`+`args` (stdio).
2. `mcp-local register ast-context-cache` or `mcp-local restart ast-context-cache`.
3. Tell the user Cursor may need **MCP reload** or restart; mcp-local only writes `mcp.json`.

### User changed ast-context-cache code

```bash
mcp-local rebuild ast-context-cache
# or
mcp-local restart ast-context-cache   # runs build_command unless no_build_on_start
```

### Add a new stdio MCP server

```bash
mcp-local add my-server \
  --command "${HOME}/path/to/binary" \
  --type stdio \
  --args "--flag" \
  --env "KEY=value"
mcp-local start my-server
mcp-local registered my-server
```

## Limitations (do not promise otherwise)

- **OpenCode JSONC:** reads strip `//` comments; writes plain JSON (comments may be lost).
- **Catalog browse:** `mcp-local browse` / MCP Registry install TUI is **not** implemented yet (see [docs/goal-alignment-plan.md](docs/goal-alignment-plan.md)).
- **Remote OAuth MCPs:** registration may work in agent JSON; mcp-local does not run OAuth flows.
- **Cursor env for HTTP tier:** HTTP entries use `url` only; tier for Cursor may require manual `env` in `mcp.json` if not using stdio.

## Developing this repository

```bash
make          # ./bin/mcp-local
make test
make check    # test + vet + fmt
```

Layout:

- `cmd/mcp-local/` — CLI (lifecycle, tools, register)
- `internal/mgr/` — config, process, agents (opencode, cursor, claudedesktop), asttools
- `internal/tui/` — config wizard, tool manager
- `examples/` — starter and service snippets

**Git:** Only commit when the user asks. Do not push unless asked.

## Related projects

- **[ast-context-cache](https://github.com/coma-toast/ast-context-cache)** — MCP server; tool policy schema in `skills/tools.json.example` and that repo’s `AGENTS.md`.
- Human-oriented overview: [README.md](README.md).
