# AGENTS.md — mcp-local

Guidance for AI coding assistants operating on or through **mcp-local**, the local MCP control plane.

## What this tool is

**mcp-local** supervises local MCP server processes and keeps agent editor configs in sync. It is **not** an MCP server—tools (`tools/list`, etc.) belong to the servers it starts (e.g. ast-context-cache). The one transport piece it implements is `mcp-local bridge <url>`, a stdio ↔ Streamable HTTP relay for hosts that can only launch stdio servers (Claude Desktop).

| Responsibility | mcp-local | MCP server (e.g. ast-mcp) |
|----------------|-----------|---------------------------|
| Start/stop/restart binaries | Yes | N/A |
| Register `mcp.json` / `opencode.jsonc` / Claude Desktop | Yes (delegates to `ast-mcp install` for ast-context-cache v4+) | `ast-mcp install` |
| stdio → HTTP relay (`mcp-local bridge`) | Yes | Serves Streamable HTTP |
| Tool tier policy file (`tools.json`) | Writes from config | Reads at startup |
| Serve MCP over HTTP/stdio | No | Yes |

## Paths (defaults)

| Path | Purpose |
|------|---------|
| `~/.mcp-local/config.yaml` | Source of truth for services and tool overrides |
| `~/.mcp-local/pids/<name>.pid` | PID plus start time of each process mcp-local started (JSON; legacy plain-PID files are still read) |
| `~/.mcp-local/logs/<name>.log` | Default service log (0600) when `log:` is not set |
| `~/.astcache/tools.json` | Per-tool `enabled` / `tier` / `description` for ast-context-cache |
| `~/.cursor/mcp.json` | Cursor MCP entries (`mcpServers`) |
| `~/.config/opencode/opencode.jsonc` (or `.json`) | OpenCode MCP entries (`mcp`) |
| Claude Desktop config | OS-specific `claude_desktop_config.json`; see `internal/mgr/claudedesktop` |
| `~/.mcp-local/backups/<YYYYMMDD-HHMMSS>/` | Original agent config before each mcp-local edit (path `/` → `%`, last 5 kept per file) |

Build output: `./bin/mcp-local` after `make` in this repo.

## When to use mcp-local (agent workflow)

**Prefer mcp-local** when the user’s stack is already defined in `config.yaml` or they want lifecycle + registration in one place.

**Do not** hand-edit the managed keys (`url`, `headers`, `command`, `args`, `env`/`environment`, `type`) of mcp-local's entries in `~/.cursor/mcp.json`, OpenCode's `mcp` block, or Claude Desktop's `mcpServers`—they are rewritten on the next `start`, `restart`, or `register`. Other keys in those entries (`enabled`, `disabled`, `timeout`, …) and everything else in those files (other servers, comments, formatting) are left untouched: mcp-local patches only `<block>/<service name>` and writes atomically after a backup.

**Do** use shell commands below; use `mcp-local status --plain` and `mcp-local validate` before assuming a service is up.

### First-time / missing config

If `~/.mcp-local/config.yaml` is missing and stdin is a terminal, many commands launch the **config wizard** TUI. Without a terminal, or when the file exists but can't be read or parsed, they print the error and path to stderr and exit 1 (the file is never overwritten). Saves keep comments, unknown keys and `${HOME}`/`~` paths, write through a symlinked config atomically, and back up the previous file under `~/.mcp-local/backups/`. For non-interactive setup, copy [examples/config.starter.yaml](examples/config.starter.yaml) and edit paths (`command`, `env`, `deps`).

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
- HTTP services: running = port in use, or the PID file still names the process mcp-local started (PID + start time, so a reused PID is never signalled).
- **stdio services are registered only, never daemonized**: each agent launches its own copy on demand. `start` runs their build and registers them; `status` shows them as `on-demand`. Remote-only services (`mcp_url`, no `command`) are likewise registered only (`remote`).
- `stop` signals the whole process group (wrapper children of `npx`/`uv`/scripts too) and reports ✅ only once the process is gone and the port is free. Without a verified PID file it stops the port owner only if it runs the service's binary.
- `restart` builds first and leaves the running copy alone if the build fails; then stop, then start without a second build. Builds run from the `command`'s directory with a 10-minute timeout.
- `start`, `stop`, `restart` and `health` exit non-zero when any service fails. `status` is read-only (it never deletes PID files).

### Agent registration (without restart)

```bash
mcp-local register [service]       # --dry-run to preview
mcp-local deregister [service]     # no name: stopped HTTP services with a port, only entries mcp-local wrote
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

If the `agents:` key is missing, mcp-local registers with OpenCode + Cursor + Claude (legacy behavior). If it is present, only hosts set to `true` are used (all `false` = none).

Entry shapes per host:

| Host | HTTP service | stdio service |
|------|--------------|---------------|
| OpenCode (`mcp`) | `{"type":"remote","url":…,"enabled":true,"timeout":30000}` | `{"type":"local","command":[…],"environment":{…}}` |
| Cursor (`mcpServers`) | `{"url":…}` | `{"command":…,"args":[…],"env":{…}}` |
| Claude Desktop (`mcpServers`) | `{"command":"<abs path>/mcp-local","args":["bridge",<mcp_url>]}` | `{"type":"stdio","command":…,"args":[…],"env":{…}}` |

Claude Desktop cannot connect to HTTP servers itself, so HTTP services run through `mcp-local bridge`.

**ast-context-cache delegation:** a service named `ast-context-cache` whose `command` basename is `ast-mcp` defaults to `installer: ast-mcp` (set `installer: native` to opt out); other names use mcp-local's writers with a warning. For those HTTP services, register/deregister run `<command> install|uninstall --target <opencode|cursor|claude_desktop> --component mcp --yes --json --mcp-url <mcp_url>` instead of mcp-local's writers. If `<command> --version` is below 4.0.0 (or fails), or the installer exits 4 (unsupported), mcp-local falls back to its own writer and prints a warning. Errors from any host are collected and reported together; other hosts still get updated.

### stdio bridge

```bash
mcp-local bridge http://localhost:7821/mcp                         # stdin/stdout JSON-RPC ↔ Streamable HTTP
mcp-local bridge https://host/mcp --header "Authorization=Bearer …" --timeout 120s
```

Newline-delimited JSON-RPC on stdin is POSTed to the URL (`Accept: application/json, text/event-stream`); JSON and SSE responses come back one message per stdout line. It keeps `Mcp-Session-Id` / `MCP-Protocol-Version`, forwards the server's GET SSE stream when offered, turns HTTP failures on requests into JSON-RPC errors (`code -32000`), exits 0 on stdin EOF (sending DELETE for the session), and logs only to stderr. On a 404 for its session it replays the cached `initialize` / `notifications/initialized` (the replayed response is not forwarded) and retries once. `--timeout` is a header and idle timeout, not a cap on the whole call. Batches get `-32600` per request, and invalid JSON gets `-32700`.

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
- `tools:` list → merged into tools.json on start (via `asttools.ApplyStartEnv`); entries for tools without a config override (e.g. set in the ast dashboard) are kept. A tool with no `enabled:` key is enabled; `tools sync` stores only your own descriptions, not the server's.

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
mcp-local status                 # live TUI (q to quit) — DETAIL column shows issues
mcp-local status --plain         # scriptable one-shot
mcp-local health <service>
mcp-local log <service>          # tail one log
mcp-local logs                   # tail all (Ctrl+C)
mcp-local open <service>         # dashboard_url in browser
mcp-local doctor                 # diagnose all services (orphaned ports, degraded health, port conflicts, stale PID files, agent registration)
mcp-local doctor my-server --fix # stop untracked copies of its own binary, drop stale PID files, restart if degraded, re-register; never stops other programs; exit 1 while problems remain
mcp-local update --check         # is a newer release out? (`mcp-local update` installs it, verified against checksums.txt)
```

- **Status states**: `running`, `degraded` (health check failed), `orphaned` (same binary on port, not managed), `port conflict` (different program on port), `stopped`.
- **Doctor checks**: PID file, owned ports (orphaned/conflict), health_url, deps, stale PID, agent registration, disk space.
- **Doctor --fix**: stops same-binary orphans, removes stale PID files, restarts degraded services. Refuses to kill different programs.

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

- **JSONC:** comments, trailing commas, key order, and formatting are preserved; only mcp-local's own entries are rewritten.
- **Entry names:** deregistration removes only an entry whose name exactly equals the service name; a bare `deregister` also skips entries mcp-local didn't write (recorded in `~/.mcp-local/written-entries.json`). Delegated `ast-mcp install` entries are named by ast-mcp, so `mcp-local registered` may not see them under a different service name.
- **Catalog browse:** `mcp-local browse` / MCP Registry install TUI is **not** implemented yet (see [docs/goal-alignment-plan.md](docs/goal-alignment-plan.md)).
- **Remote OAuth MCPs:** registration may work in agent JSON; mcp-local does not run OAuth flows.
- **Cursor env for HTTP tier:** HTTP entries use `url` only; tier for Cursor may require manual `env` in `mcp.json` if not using stdio.

## Developing this repository

```bash
make          # ./bin/mcp-local
make test
make check    # test + vet + gofmt -l (fails on unformatted code, rewrites nothing)
make fmt      # gofmt -w
make install INSTALL_DIR=~/.local/bin   # default $(PREFIX)/bin, PREFIX ?= /usr/local; creates the dir
```

The Makefile exports `GOWORK=off`, so run `go` commands outside make with `GOWORK=off` too if a parent directory has a `go.work`.

Layout:

- `cmd/mcp-local/` — CLI (lifecycle, tools, register)
- `internal/mgr/` — config, process, agents (opencode, cursor, claudedesktop), asttools
- `internal/tui/` — config wizard, tool manager
- `examples/` — starter and service snippets

**Git:** Only commit when the user asks. Do not push unless asked.

## Related projects

- **[ast-context-cache](https://github.com/coma-toast/ast-context-cache)** — MCP server; tool policy schema in `skills/tools.json.example` and that repo’s `AGENTS.md`.
- Human-oriented overview: [README.md](README.md).
