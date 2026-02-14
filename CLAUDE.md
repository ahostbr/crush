# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is Crush?

Crush is a terminal-based AI coding assistant built by Charm. It supports multiple LLM providers (Anthropic, OpenAI, Google, Bedrock, etc.) via OpenAI- and Anthropic-compatible APIs, with features like session management, LSP integration, MCP server support, and tool execution. Written in Go, built on the Charm ecosystem (Bubble Tea, Lip Gloss, Glamour).

## Build/Test/Lint Commands

- **Build**: `go build .` or `task build`
- **Run**: `task run` or `go run .`
- **Test**: `task test` or `go test -race -failfast ./...`
- **Run single test**: `go test ./internal/agent -run TestCoderAgent`
- **Update golden files**: `go test ./... -update` (regenerates `.golden` files)
- **Update specific golden**: `go test ./internal/ui/diffview -update`
- **Lint**: `task lint` (includes log capitalization check)
- **Lint with autofix**: `task lint:fix`
- **Format**: `task fmt` (runs `gofumpt -w .`)
- **Modernize**: `task modernize`
- **Record VCR cassettes**: `task test:record` (re-records all agent test cassettes)
- **Dev with profiling**: `task dev` (sets `CRUSH_PROFILE=true`, pprof on :6060)
- **Generate DB code**: `sqlc generate` (from `sqlc.yaml`)
- **Generate schema**: `task schema`

Environment: `CGO_ENABLED=0`, `GOEXPERIMENT=greenteagc`.

## Architecture

### Entry Point & CLI (`main.go`, `internal/cmd/`)

`main.go` calls `cmd.Execute()`. CLI uses Cobra with subcommands: root (interactive TUI), `run` (non-interactive), `dirs`, `projects`, `update-providers`, `logs`, `schema`, `login`, `stats`.

### Application Layer (`internal/app/`)

`App` wires together all services: sessions, messages, permissions, file tracking, agent coordinator, LSP manager. It initializes the DB connection, creates services, and manages the Bubble Tea program lifecycle. This is the central orchestration point.

### Agent System (`internal/agent/`)

- **`Coordinator`** interface manages running agents across sessions with queueing, cancellation, and summarization.
- **`SessionAgent`** interface runs a single agent call with a model/tools against a session.
- **Prompts** are Go templates embedded via `//go:embed` from `templates/` (e.g., `coder.md.tpl`, `task.md.tpl`, `initialize.md.tpl`).
- **`fantasy`** (`charm.land/fantasy`) is the LLM abstraction layer supporting all provider types (Anthropic, OpenAI, Bedrock, Google, Azure, Vercel, OpenRouter, etc.).
- **`catwalk`** (`charm.land/catwalk`) is the community-managed provider/model registry.

### Tools (`internal/agent/tools/`)

Each tool is a separate file (e.g., `bash.go`, `edit.go`, `grep.go`, `view.go`, `write.go`). Tool descriptions are in adjacent `.md` or `.tpl` files. MCP tools are in `tools/mcp/`. Tools receive context with session/message IDs.

### TUI (`internal/ui/`)

Built on Bubble Tea v2. Read `internal/ui/AGENTS.md` before working on UI code. Key principles:
- Main model in `model/ui.go` handles routing, focus, layout.
- Components are "dumb": expose methods, return `tea.Cmd`, render via `Render(width int) string`.
- Chat items in `chat/` are simple renderers with caching.
- Styles in `styles/styles.go` via `*common.Common`.
- Never do IO in `Update`; always use `tea.Cmd`.
- Use `github.com/charmbracelet/x/ansi` for ANSI string manipulation.

### Database (`internal/db/`)

SQLite via sqlc. Migrations in `internal/db/migrations/`. SQL queries in `internal/db/sql/`. Generated code in `internal/db/*.sql.go` — do NOT edit generated files. Two SQLite drivers available: `ncruces` and `modernc` (selected via build tags in `connect_*.go`).

### Config (`internal/config/`)

JSON config loaded from `.crush.json` / `crush.json` / `~/.config/crush/crush.json` (priority order). Config includes providers, LSP settings, MCP servers, permissions, and options. Reads context from `AGENTS.md`, `CLAUDE.md`, `.cursorrules`, etc. automatically.

### Key Supporting Packages

- `internal/session/` — Session CRUD, todo tracking, pub/sub notifications.
- `internal/message/` — Message content types and attachments.
- `internal/permission/` — Tool permission management (allowed/denied tools).
- `internal/lsp/` — LSP client manager for code intelligence.
- `internal/shell/` — Background shell command execution.
- `internal/event/` — Metrics/telemetry event types and logging.
- `internal/pubsub/` — Generic typed pub/sub broker.
- `internal/csync/` — Concurrent data structures (maps, slices, values).

## Code Style

- **Formatting**: Always use `gofumpt`. Fallback: `goimports`, then `gofmt`.
- **Imports**: Group stdlib, external, internal packages.
- **Testing**: Use `testify/require`, `t.Parallel()`, `t.SetEnv()`, `t.TempDir()`.
- **Mock providers**: Set `config.UseMockProviders = true` and call `config.ResetProviders()` in tests.
- **Log messages**: Must start with a capital letter (enforced by `task lint:log`).
- **Comments**: Own-line comments start with capitals, end with periods. Wrap at 78 columns.
- **JSON tags**: Use `snake_case`.
- **File permissions**: Use octal notation (`0o755`, `0o644`).
- **Commits**: Use semantic commits (`fix:`, `feat:`, `chore:`, `refactor:`, `docs:`, `sec:`). Keep to one line.
- **Stats HTML/CSS/JS**: Format with `prettier` (`task fmt:html`).
