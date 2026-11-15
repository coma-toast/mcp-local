# mcp-local

`mcp-local` is a unified local control plane for Model Context Protocol (MCP) servers. Define your stdio, HTTP, remote, and embedded filesystem servers once in `~/.mcp-local/config.yaml`, and `mcp-local` handles process lifecycles, agent config syncing, and tool-tier overrides.

> **AI agents:** See [AGENTS.md](AGENTS.md) for command workflows, file paths, and `ast-context-cache` tier details.

---

## Quick Start (One-Shot Agent & IDE Setup)

Install `mcp-local`, add its built-in `mcp` server, and integrate with OpenCode, Cursor, and Claude Desktop in one task:

```bash
# 1. Build and install (or grab a release tarball from GitHub Releases)
git clone https://github.com/coma-toast/mcp-local.git && cd mcp-local
make install INSTALL_DIR=~/.local/bin

# 2. Initialize config.yaml if missing
mkdir -p ~/.mcp-local
[ -f ~/.mcp-local/config.yaml ] || cp examples/config.starter.yaml ~/.mcp-local/config.yaml

# 3. Add and register `mcp-local mcp` across enabled agent IDEs
~/.local/bin/mcp-local add mcp-local \
  --type stdio \
  --command "${HOME}/.local/bin/mcp-local" \
  --args mcp \
  --port 0 </dev/null
~/.local/bin/mcp-local register mcp-local
~/.local/bin/mcp-local registered mcp-local
```

For any other MCP-capable IDE (Antigravity, Windsurf, Zed, etc.), point a stdio MCP entry at `~/.local/bin/mcp-local mcp`.

<details>
<summary>Alternative install options (release binary, custom prefix, source build)</summary>

- **Prebuilt binaries:** Every merge to `main` publishes macOS/Linux (`amd64`/`arm64`) tarballs + `checksums.txt` on [GitHub Releases](https://github.com/coma-toast/mcp-local/releases):
  ```bash
  tar -xzf mcp-local_<version>_darwin_arm64.tar.gz
  install -m 0755 mcp-local ~/.local/bin/
  ```
- **Custom install path:** `make install PREFIX=~/.local` (default `/usr/local/bin`; use `sudo make install` if needed).
- **Self-update:** `mcp-local update --check`, `mcp-local update` (verifies SHA-256 against `checksums.txt`), or `mcp-local update --version v0.1.9 --force`. Set `GH_TOKEN`/`GITHUB_TOKEN` if rate-limited.

</details>

---

## Features

- **Lifecycle & Diagnostics:** `start`, `stop` (`--deregister` default `true`, verifies PID + start time and process group), `restart` (builds before stopping), `rebuild` (runs in `command` dir, 10m timeout), `status` (TUI or `--plain`), `doctor [--fix]` (orphans, port conflicts, stale PIDs, missing deps/registrations, low disk), `health`, `log`/`logs`, `open` (`http`/`https` `dashboard_url`).
- **Multi-Agent Registration:** Syncs **OpenCode** (`~/.config/opencode/opencode.jsonc`), **Cursor** (`~/.cursor/mcp.json`), and **Claude Desktop** (`claude_desktop_config.json`, via `mcp-local bridge` for HTTP). Preserves JSONC comments, trailing commas, key order, UTF-8 BOM, CRLF, and user keys; backs up originals to `~/.mcp-local/backups/` (last 5 kept); refuses files with duplicate keys. Delegates `ast-context-cache` (`ast-mcp` v4+) to `ast-mcp install|uninstall` unless `installer: native`.
- **Control-Plane MCP (`mcp-local mcp`):** Exposes 12 tools over stdio JSON-RPC 2.0 so agents can manage your stack directly:
  - *Services:* `list_services`, `service_status`, `start_service`, `stop_service`, `restart_service`, `add_service`, `remove_service`
  - *Agents:* `register_service`, `deregister_service`
  - *Tools:* `list_tools` (config or `live: true`), `sync_tools`, `configure_tools`
- **stdio ↔ HTTP Bridge (`mcp-local bridge <url>`):** Spec `2025-06-18` relay with session/protocol headers, SSE + GET notification stream reconnect (`Last-Event-ID`), transparent `404` session re-initialization, `--header`, and idle `--timeout` (`60s` default).
- **Embedded Filesystem MCP (`type: filesystem`):** Serves `path` over Streamable HTTP (`/mcp`, `/health`) bound to `127.0.0.1:<port>` with 6 tools (`read_file`, `write_file`, `list_directory`, `search_files`, `get_file_info`, `create_directory`). Enforces localhost `Host`/`Origin`, `application/json` `POST`s, symlink-safe path resolution, 10 MiB body/read caps, and a 1,000-match search cap.
- **Tool & Tier Management (`ast-context-cache`):** `tools sync`, `tools` TUI, `tools apply` (`~/.astcache/tools.json`), `json tools`, plus `AST_MCP_TIER` / `AST_MCP_CODE_MODE` / `AST_MCP_TOOLS_CONFIG` injection on start (restart after changes).
- **Config Management:** First-run setup wizard TUI, `add`, `edit` (TUI), `remove`, `config` (`$EDITOR`/`$VISUAL` with flag support), `validate`, `import --from opencode|cursor|claude`, `list`.

---

## Top 25 Popular MCP Tools (Copy-Paste 1-Liners)

Copy any 1-liner below to add a popular MCP server to `~/.mcp-local/config.yaml`, then run `mcp-local register <name>` (or `mcp-local register` to sync all) for OpenCode, Cursor, and Claude Desktop. Full YAML templates and 5 bonus configs live in [examples/popular/](examples/popular/).

> `</dev/null` skips interactive port/log prompts on `stdio` servers. Replace `YOUR_*` placeholders before registering.

### Core, Filesystem & Reasoning

```bash
# 1. Filesystem — read/write/search allowed directories (or use built-in: mcp-local add my-files --type filesystem --path ~/projects --port 4000 --health http://localhost:4000/health)
mcp-local add filesystem --type stdio --command npx --args -y --args @modelcontextprotocol/server-filesystem --args "${HOME}/projects" </dev/null

# 2. GitHub — official remote MCP (OAuth; see examples/popular/github.service.yaml for PAT or Docker)
mcp-local add github --type http --mcp-url https://api.githubcopilot.com/mcp/

# 3. Git — status, diff, log, and repo tools (requires uv)
mcp-local add git --type stdio --command uvx --args mcp-server-git </dev/null

# 4. Sequential Thinking — structured step-by-step problem solving
mcp-local add sequential-thinking --type stdio --command npx --args -y --args @modelcontextprotocol/server-sequential-thinking </dev/null

# 5. Memory — persistent knowledge-graph memory across chats
mcp-local add memory --type stdio --command npx --args -y --args @modelcontextprotocol/server-memory </dev/null

# 6. Desktop Commander — terminal commands, processes, and file editing
mcp-local add desktop-commander --type stdio --command npx --args -y --args @wonderwhy-er/desktop-commander </dev/null

# 7. Time — current time and IANA timezone conversion (requires uv)
mcp-local add time --type stdio --command uvx --args mcp-server-time </dev/null
```

### Browser Automation & DevTools

```bash
# 8. Playwright — browser automation via accessibility snapshots
mcp-local add playwright --type stdio --command npx --args -y --args @playwright/mcp </dev/null

# 9. Chrome DevTools — live Chrome console, network, DOM, and performance traces
mcp-local add chrome-devtools --type stdio --command npx --args -y --args chrome-devtools-mcp </dev/null
```

### Search, Scraping & Live Docs

```bash
# 10. Context7 — up-to-date, version-specific library docs
mcp-local add context7 --type stdio --command npx --args -y --args @upstash/context7-mcp </dev/null

# 11. Fetch — fetch web pages as Markdown (requires uv)
mcp-local add fetch --type stdio --command uvx --args mcp-server-fetch </dev/null

# 12. Brave Search — web, local, image, news, and video search
mcp-local add brave-search --type stdio --command npx --args -y --args @brave/brave-search-mcp-server --env BRAVE_API_KEY=YOUR_KEY </dev/null

# 13. Firecrawl — scrape, crawl, map, and search web pages to Markdown
mcp-local add firecrawl --type stdio --command npx --args -y --args firecrawl-mcp --env FIRECRAWL_API_KEY=YOUR_KEY </dev/null

# 14. Tavily — AI-native web search, extraction, and crawling
mcp-local add tavily --type stdio --command npx --args -y --args tavily-mcp --env TAVILY_API_KEY=YOUR_KEY </dev/null

# 15. Exa — semantic web search and research (hosted Streamable HTTP)
mcp-local add exa --type http --mcp-url https://mcp.exa.ai/mcp
```

### Databases & Backend

```bash
# 16. Postgres — schema inspection, explain plans, index tuning, read-only SQL (requires uv)
mcp-local add postgres --type stdio --command uvx --args postgres-mcp --args --access-mode=restricted --env DATABASE_URI=postgresql://user:pass@localhost:5432/mydb </dev/null

# 17. SQLite — query and inspect a local SQLite database (requires uv)
mcp-local add sqlite --type stdio --command uvx --args mcp-server-sqlite --args --db-path --args "${HOME}/data.db" </dev/null

# 18. Supabase — SQL queries, migrations, logs, and project management (OAuth)
mcp-local add supabase --type http --mcp-url https://mcp.supabase.com/mcp
```

### Productivity, Design & Ops

```bash
# 19. Notion — workspace search and Markdown page read/edit (OAuth)
mcp-local add notion --type http --mcp-url https://mcp.notion.com/mcp

# 20. Linear — issues, projects, cycles, and teams (OAuth)
mcp-local add linear --type http --mcp-url https://mcp.linear.app/mcp

# 21. Sentry — issue triage, stack traces, and release analysis
mcp-local add sentry --type stdio --command npx --args -y --args @sentry/mcp-server --env SENTRY_AUTH_TOKEN=YOUR_TOKEN </dev/null

# 22. Slack — workspace channels, messages, threads, and reactions
mcp-local add slack --type stdio --command npx --args -y --args @zencoderai/slack-mcp-server --env SLACK_BOT_TOKEN=xoxb-YOUR_TOKEN --env SLACK_TEAM_ID=T01234567 </dev/null

# 23. Figma — design frames, layout, and tokens via Framelink
mcp-local add figma --type stdio --command npx --args -y --args figma-developer-mcp --args --stdio --env FIGMA_API_KEY=YOUR_KEY </dev/null

# 24. Kubernetes — cluster pods, logs, deployments, and manifests (~/.kube/config)
mcp-local add kubernetes --type stdio --command npx --args -y --args kubernetes-mcp-server </dev/null

# 25. Grafana — dashboards, PromQL/LogQL queries, and alerts (requires uv)
mcp-local add grafana --type stdio --command uvx --args mcp-grafana --env GRAFANA_URL=http://localhost:3000 --env GRAFANA_SERVICE_ACCOUNT_TOKEN=YOUR_TOKEN </dev/null
```

More examples (`atlassian`, `cloudflare`, `n8n`, `obsidian`, `stripe`) are in [examples/popular/](examples/popular/).

---

## Command Reference

```bash
# Lifecycle & Monitoring
mcp-local start [service | --all]           # Start & register with agents
mcp-local stop [service | --all]            # Stop & deregister (--deregister=false to keep)
mcp-local restart [service | --all]         # Build, stop, start, register
mcp-local rebuild <service>                 # Run build_command & re-register
mcp-local status [service] [--plain]        # Live TUI or one-shot table
mcp-local doctor [service] [--fix]          # Diagnose & apply safe fixes
mcp-local health [service]                  # Check health_url
mcp-local log <service> | logs              # Tail one or all service logs
mcp-local open <service>                    # Open dashboard_url in browser

# Config & Registration
mcp-local list | validate | config          # List, validate, or edit in $EDITOR
mcp-local add <name> [flags]                # Add/update service
mcp-local edit <service> | remove <service> # TUI editor or deregister & delete
mcp-local import [--from opencode|cursor|claude]
mcp-local register [service] [--dry-run]    # Sync entries to enabled agents
mcp-local deregister [service]              # Remove entries mcp-local wrote
mcp-local registered [service]              # Show registration status per agent

# MCP, Bridge & Tools
mcp-local mcp                               # Run control-plane stdio MCP server
mcp-local bridge <url> [--header k=v] [--timeout 60s]
mcp-local tools [sync|apply <service>]      # Tool tier TUI, sync, or write tools.json
mcp-local json tools <service>              # Preview tools.json overrides
mcp-local update [--check|--version vX.Y.Z|--force]
```

---

## Configuration (`~/.mcp-local/config.yaml`)

```yaml
agents: # omit to enable all; when present, only true hosts are managed
  opencode: true
  cursor: true
  claude: true

services:
  - name: ast-context-cache
    command: ${HOME}/git/ast-context-cache/ast-mcp
    port: 7821
    type: http
    mcp_url: http://localhost:7821/mcp
    health_url: http://localhost:7821/health
    dashboard_url: http://localhost:7830
    active_tier: complete
    build_command: "make build"
    env:
      ONNXRUNTIME_LIB: /opt/homebrew/lib/libonnxruntime.dylib

  - name: model-proxy
    command: ${HOME}/git/model-proxy/bin/model-proxy
    type: stdio
    args: ["--config", "${HOME}/.config/model-proxy/config.yaml"]

  - name: remote-api
    type: http
    mcp_url: https://api.example.com/mcp

  - name: my-files
    type: filesystem
    path: ~/projects
    port: 4000
```

<details>
<summary>All service config fields</summary>

| Field | Description |
|-------|-------------|
| `name` | Service identifier |
| `type` | `http`, `stdio`, or `filesystem` |
| `command`, `args` | Binary path (supports `${HOME}`, `~/`) and arguments |
| `port`, `extra_ports` | Primary TCP port and additional owned local ports |
| `mcp_url`, `headers` | MCP endpoint URL and HTTP headers (`bridge --header` for Claude) |
| `health_url`, `dashboard_url` | Health endpoint and web UI URL (`http`/`https`) |
| `log` | Log file path (default `~/.mcp-local/logs/<name>.log`) |
| `env`, `deps` | Environment variables and required file paths checked before start |
| `build_command`, `no_build_on_start` | Shell build command (runs in `command` dir, 10m timeout) and skip flag |
| `path` | Root directory for `type: filesystem` |
| `active_tier`, `code_mode`, `no_code_mode` | Injected as `AST_MCP_TIER` and `AST_MCP_CODE_MODE` |
| `tools_config_path`, `tools` | Path to `tools.json` (default `~/.astcache/tools.json`) and per-tool overrides |
| `installer` | `ast-mcp` (default for `ast-mcp` binary) or `native` |
| `description` | Human-readable description |

</details>

---

## Development & Docs

```bash
make          # build ./bin/mcp-local
make check    # go test + go vet + gofmt -l
make fmt      # gofmt -w .
```

`Makefile` sets `GOWORK=off` so parent `go.work` files don't affect builds. CI runs `gofmt`, `go vet`, `go test -race`, `staticcheck`, and `go build`.

- [AGENTS.md](AGENTS.md) — Agent workflow, paths, and limitations
- [docs/goal-alignment-plan.md](docs/goal-alignment-plan.md) — Architecture history and future catalog TUI
- [examples/](examples/) — Starter `config.yaml`, service snippets, and [popular MCP configs](examples/popular/)
- **License:** MIT
