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
Stop juggling terminal tabs. Manage all your servers with simple commands:
- `mcp-local start`: Launches all configured servers and automatically registers them with your agent.
- `mcp-local stop`: Gracefully shuts down all managed services.
- `mcp-local rebuild`: Stops, rebuilds from source (via `make` or custom scripts), and restarts.
- `mcp-local status`: A live overview of which services are healthy and running.

### ⚙️ Tool & Tier Management (TUI)
Included is a professional Terminal User Interface (TUI) to manage the specific capabilities of your servers:
- **Enable/Disable**: Turn specific tools on or off without editing code.
- **Tier Control**: Assign tools to different access levels (`core`, `extended`, `complete`).
- **Custom Descriptions**: Override tool descriptions to help your agent understand when to use them.

### 🔄 Automatic Agent Registration
No more manual JSON editing. When you run `start`, `mcp-local` automatically syncs your running services into:
- **OpenCode**: `~/.config/opencode/opencode.jsonc`
- **Cursor**: `.cursor/mcp.json` (via configuration)

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
make build
sudo mv mcp-local /usr/local/bin/
```

### Configuration
Create your config at `~/.mcp-local/config.yaml`:

```yaml
services:
  - name: ast-context-cache
    command: ${HOME}/git/ast-context-cache/ast-mcp
    port: 7821
    type: http
    env:
      ONNXRUNTIME_LIB: /usr/local/lib/libonnxruntime.dylib
    build_cmd: "cd ${HOME}/git/ast-context-cache && make build"
    deps:
      - ${HOME}/git/ast-context-cache/model/model.onnx

  - name: model-proxy
    command: ${HOME}/git/model-proxy/model-proxy
    type: stdio
```

### Usage
```bash
mcp-local start    # Start and register everything
mcp-local tools    # Manage tools/tiers via TUI
mcp-local status   # Check health
mcp-local stop     # Shut down everything
```

## 📜 License
MIT
