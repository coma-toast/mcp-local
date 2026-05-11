# mcp-local 🚀

`mcp-local` is a unified local control plane for managing Model Context Protocol (MCP) servers. Instead of managing multiple binaries, ports, and agent configurations manually, `mcp-local` provides a single point of control for your local AI infrastructure.

## 🌟 Why mcp-local?

When building an advanced AI agent setup, you often end up with a fragmented ecosystem:
- Multiple MCP servers running on different ports.
- Complex environment variables (like `ONNXRUNTIME_LIB`) needed for specific binaries.
- Manual updates to `opencode.jsonc` or `mcp.json` every time a server changes.
- Difficulty syncing your "local AI stack" across multiple machines.

`mcp-local` solves this by treating your MCP infrastructure as code. You define your servers in a single portable config, and `mcp-local` handles the rest.

## ✨ Key Features

### 🛠️ Unified Lifecycle Management
Stop juggling terminal tabs. Manage servers with:
- `mcp-local start <service>` or `mcp-local start --all` — launches and registers with OpenCode (`~/.config/opencode/opencode.json` / `.jsonc`).
- `mcp-local stop <service>` / `stop --all`, `restart <service>` / `restart --all`.
- `mcp-local rebuild <service>` — runs `build_cmd` / `build_command` from config.
- `mcp-local status` — bubbletea live table; `mcp-local status --plain` for scripting.

### ⚙️ Tool & Tier Management (TUI)
Included is a professional Terminal User Interface (TUI) to manage the specific capabilities of your servers:
- **Enable/Disable**: Turn specific tools on or off without editing code.
- **Tier Control**: Assign tools to different access levels (`core`, `extended`, `complete`).
- **Custom Descriptions**: Override tool descriptions to help your agent understand when to use them.

### 🔄 Automatic Agent Registration
When you run `start`, `mcp-local` merges running services into OpenCode’s top-level **`mcp`** block under **`~/.config/opencode/opencode.jsonc`** (or **`opencode.json`**). Use **`register`** / **`deregister`** for explicit MCP URL updates.

### 📄 Portable Configuration
All settings are stored in `~/.mcp-local/config.yaml`. Sync this file across your machines (via Dropbox, iCloud, or Git) to have an identical AI toolset everywhere.

## 🤝 Supported Ecosystem

### Supported Agents
- **OpenCode**: Full integration with `opencode.jsonc`.
- **Cursor**: Support for `mcp.json` registration.
- **Any MCP-compliant client**: Since it manages standard MCP servers, it works with any client that supports HTTP or stdio transports.

### Supported MCP Servers
While it can manage *any* binary, it is optimized for:
- **ast-context-cache**: Full lifecycle and TUI tool management.
- **model-proxy**: Process management and registration.
- **Storybook**: Dev server lifecycle.
- **Custom Servers**: Simply add any binary to `config.yaml`.

## 🚀 Getting Started

### Installation
```bash
git clone https://github.com/coma-toast/mcp-local.git
cd mcp-local
make          # writes ./bin/mcp-local
# Optional: install system-wide
# sudo make install
```

Add `./bin` to your PATH (e.g. Fish: `fish_add_path -m ~/git/mcp-local/bin` after first build).

### Configuration
Create your config at `~/.mcp-local/config.yaml`.

**ast-context-cache:** use the ready-made starter ([`examples/config.starter.yaml`](examples/config.starter.yaml)) or merge the service fragment from [`examples/ast-context-cache.service.yaml`](examples/ast-context-cache.service.yaml). Adjust `ONNXRUNTIME_LIB` if needed (Apple Silicon Homebrew: `/opt/homebrew/lib/libonnxruntime.dylib`; Intel macOS often `/usr/local/lib/libonnxruntime.dylib`).

```yaml
services:
  - name: ast-context-cache
    command: ${HOME}/git/ast-context-cache/ast-mcp
    port: 7821
    type: http
    mcp_url: http://localhost:7821/mcp
    health_url: http://localhost:7821/health
    dashboard_url: http://localhost:7830
    log: ${HOME}/.mcp-local/ast-context-cache.log
    env:
      ONNXRUNTIME_LIB: /opt/homebrew/lib/libonnxruntime.dylib
    build_command: "cd ${HOME}/git/ast-context-cache && make build"
    deps:
      - ${HOME}/git/ast-context-cache/model/model.onnx

  - name: model-proxy
    command: ${HOME}/git/model-proxy/model-proxy
    type: stdio
```

### Usage
```bash
mcp-local start --all          # Start and register every service
mcp-local status               # Live table (q to quit)
mcp-local status --plain       # One-shot text output
mcp-local tools                # Tool tier TUI
mcp-local stop --all
mcp-local register             # Push mcp_url entries to OpenCode config
```

## 📜 License
MIT
