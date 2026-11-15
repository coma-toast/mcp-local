# Popular MCP Server Example Configs

Ready-to-use `mcp-local` configs for 30 of the most popular MCP servers.

## Two Ways to Add Any Server

1. **1-liner CLI (recommended):** Run the `mcp-local add ...` command shown in the table below (or in the top of each `.service.yaml` file), then register with your agents:
   ```bash
   mcp-local register <name>
   ```
   > **Note:** `</dev/null` on `stdio` commands skips the interactive port/log prompts so you can copy-paste without pressing Enter.

2. **Copy-paste YAML:** Copy the `- name: ...` block from any `*.service.yaml` file into `~/.mcp-local/config.yaml` under `services:`, adjust any placeholder keys or paths, then run:
   ```bash
   mcp-local validate && mcp-local register
   ```

## Prerequisites by Transport

| Transport / Runner | Requirement |
|---|---|
| `npx` (`stdio`) | Node.js 18+ (22+ for `obsidian`, `sentry`) |
| `uvx` (`stdio`) | [`uv`](https://docs.astral.sh/uv/) (`brew install uv` or `curl -LsSf https://astral.sh/uv/install.sh \| sh`) |
| `docker` (`stdio`) | Docker daemon running |
| `http` (remote) | Agent/client supports Streamable HTTP + OAuth (or `headers:` PAT in YAML) |

---

## Catalog (Top 25 + 5 Extras)

### Core, Filesystem & Reasoning

| # | Name | Transport | Runner / Endpoint | Auth / Notes | Config |
|---|------|-----------|-------------------|--------------|--------|
| 1 | `filesystem` | `stdio` | `npx @modelcontextprotocol/server-filesystem` | Pass allowed directory path(s) | [filesystem.service.yaml](filesystem.service.yaml) |
| 2 | `github` | `http` | `https://api.githubcopilot.com/mcp/` | OAuth or `Authorization: Bearer` header (Docker option in YAML) | [github.service.yaml](github.service.yaml) |
| 3 | `git` | `stdio` | `uvx mcp-server-git` | Optional `--repository` | [git.service.yaml](git.service.yaml) |
| 4 | `sequential-thinking` | `stdio` | `npx @modelcontextprotocol/server-sequential-thinking` | None | [sequential-thinking.service.yaml](sequential-thinking.service.yaml) |
| 5 | `memory` | `stdio` | `npx @modelcontextprotocol/server-memory` | None (optional `MEMORY_FILE_PATH`) | [memory.service.yaml](memory.service.yaml) |
| 6 | `desktop-commander` | `stdio` | `npx @wonderwhy-er/desktop-commander` | None | [desktop-commander.service.yaml](desktop-commander.service.yaml) |
| 7 | `time` | `stdio` | `uvx mcp-server-time` | None | [time.service.yaml](time.service.yaml) |

### Browser Automation & DevTools

| # | Name | Transport | Runner / Endpoint | Auth / Notes | Config |
|---|------|-----------|-------------------|--------------|--------|
| 8 | `playwright` | `stdio` | `npx @playwright/mcp` | None | [playwright.service.yaml](playwright.service.yaml) |
| 9 | `chrome-devtools` | `stdio` | `npx chrome-devtools-mcp` | Local Chrome installed | [chrome-devtools.service.yaml](chrome-devtools.service.yaml) |

### Search, Scraping & Live Docs

| # | Name | Transport | Runner / Endpoint | Auth / Notes | Config |
|---|------|-----------|-------------------|--------------|--------|
| 10 | `context7` | `stdio` | `npx @upstash/context7-mcp` | Optional `--api-key` | [context7.service.yaml](context7.service.yaml) |
| 11 | `fetch` | `stdio` | `uvx mcp-server-fetch` | None | [fetch.service.yaml](fetch.service.yaml) |
| 12 | `brave-search` | `stdio` | `npx @brave/brave-search-mcp-server` | `BRAVE_API_KEY` | [brave-search.service.yaml](brave-search.service.yaml) |
| 13 | `firecrawl` | `stdio` | `npx firecrawl-mcp` | `FIRECRAWL_API_KEY` | [firecrawl.service.yaml](firecrawl.service.yaml) |
| 14 | `tavily` | `stdio` | `npx tavily-mcp` | `TAVILY_API_KEY` | [tavily.service.yaml](tavily.service.yaml) |
| 15 | `exa` | `http` | `https://mcp.exa.ai/mcp` | Optional key (local `stdio` option in YAML) | [exa.service.yaml](exa.service.yaml) |

### Databases & Backend

| # | Name | Transport | Runner / Endpoint | Auth / Notes | Config |
|---|------|-----------|-------------------|--------------|--------|
| 16 | `postgres` | `stdio` | `uvx postgres-mcp` | `DATABASE_URI` (`--access-mode=restricted` default) | [postgres.service.yaml](postgres.service.yaml) |
| 17 | `sqlite` | `stdio` | `uvx mcp-server-sqlite` | `--db-path` to `.db` file | [sqlite.service.yaml](sqlite.service.yaml) |
| 18 | `supabase` | `http` | `https://mcp.supabase.com/mcp` | OAuth | [supabase.service.yaml](supabase.service.yaml) |

### Productivity, Design & Monitoring

| # | Name | Transport | Runner / Endpoint | Auth / Notes | Config |
|---|------|-----------|-------------------|--------------|--------|
| 19 | `notion` | `http` | `https://mcp.notion.com/mcp` | OAuth (legacy `stdio` option in YAML) | [notion.service.yaml](notion.service.yaml) |
| 20 | `linear` | `http` | `https://mcp.linear.app/mcp` | OAuth | [linear.service.yaml](linear.service.yaml) |
| 21 | `sentry` | `stdio` | `npx @sentry/mcp-server` | `SENTRY_AUTH_TOKEN` (or `auth login`) | [sentry.service.yaml](sentry.service.yaml) |
| 22 | `slack` | `stdio` | `npx @zencoderai/slack-mcp-server` | `SLACK_BOT_TOKEN`, `SLACK_TEAM_ID` | [slack.service.yaml](slack.service.yaml) |
| 23 | `figma` | `stdio` | `npx figma-developer-mcp` | `FIGMA_API_KEY` (or remote `https://mcp.figma.com/mcp`) | [figma.service.yaml](figma.service.yaml) |
| 24 | `kubernetes` | `stdio` | `npx kubernetes-mcp-server` | Uses `~/.kube/config` | [kubernetes.service.yaml](kubernetes.service.yaml) |
| 25 | `grafana` | `stdio` | `uvx mcp-grafana` | `GRAFANA_URL`, `GRAFANA_SERVICE_ACCOUNT_TOKEN` | [grafana.service.yaml](grafana.service.yaml) |

### Bonus Popular Servers (26–30)

| # | Name | Transport | Runner / Endpoint | Auth / Notes | Config |
|---|------|-----------|-------------------|--------------|--------|
| 26 | `atlassian` | `http` | `https://mcp.atlassian.com/v1/mcp` | OAuth (Jira & Confluence) | [atlassian.service.yaml](atlassian.service.yaml) |
| 27 | `cloudflare` | `http` | `https://docs.mcp.cloudflare.com/mcp` | None for Docs (Wrangler API option in YAML) | [cloudflare.service.yaml](cloudflare.service.yaml) |
| 28 | `n8n` | `stdio` | `npx n8n-mcp` | Optional `N8N_API_URL` + `N8N_API_KEY` | [n8n.service.yaml](n8n.service.yaml) |
| 29 | `obsidian` | `stdio` | `npx obsidian-mcp@2` | Vault path (`--vault name=/path`) | [obsidian.service.yaml](obsidian.service.yaml) |
| 30 | `stripe` | `http` | `https://mcp.stripe.com` | OAuth | [stripe.service.yaml](stripe.service.yaml) |
