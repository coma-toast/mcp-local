# mcp-local

A unified local management tool for MCP servers.

## Features
- Lifecycle management (start/stop/restart) for multiple MCP servers.
- Automatic registration with agents (opencode.jsonc).
- Portable configuration via `~/.mcp-local/config.yaml`.
- Unified log streaming.
- Dependency checks before startup.

## Installation
```bash
go build -o mcp-local ./cmd/mcp-local/main.go
mv mcp-local /usr/local/bin/
```
