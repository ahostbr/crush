# LiteHarness Integration Plan for Crush

## Goal

Crush should be able to receive LiteHarness inbox messages while an interactive Crush TUI session is already running. The integration should live in the Go codebase, not as a Python sidecar. It should poll the same LiteHarness maildir contract used by other agents, then inject received message bodies into the active Bubble Tea session as normal user prompts.

## Current Crush Architecture

Crush is a Go CLI with a Bubble Tea TUI.

Relevant entry points:

- `main.go` starts the CLI command tree.
- `internal/cmd/root.go` creates the interactive Bubble Tea program with `tea.NewProgram(...)`.
- `internal/app/app.go` owns app services and an `events chan tea.Msg`.
- `internal/app/app.go` `Subscribe(program)` forwards app events into Bubble Tea with `program.Send(msg)`.
- `internal/ui/model/ui.go` owns the main UI model and sends user prompts through `m.sendMessage(...)`.
- `internal/ui/model/ui.go` currently has an unexported `sendMessageMsg` used internally to ask the UI to send a prompt.

The useful in-process injection path is:

```text
LiteHarness watcher goroutine
  -> app.events or program.Send(custom tea.Msg)
  -> UI Update handles message
  -> m.sendMessage(body)
  -> AgentCoordinator.Run(...)
```

This path preserves the current session, avoids terminal focus issues, and avoids treating received text as keyboard paste.

## Skill Loading Findings

Crush implements the Agent Skills open-standard style in `internal/skills/skills.go`.

Skill discovery works like this:

1. `internal/config/config.go` exposes `options.skills_paths` as the configurable skill root list.
2. `internal/config/load.go` adds default global skill directories during config defaulting.
3. `internal/agent/prompt/prompt.go` expands `cfg.Options.SkillsPaths`, calls `skills.Discover(...)`, and converts results to prompt XML with `skills.ToPromptXML(...)`.
4. `internal/agent/templates/coder.md.tpl` injects the resulting `<available_skills>` block into the coder prompt.
5. The agent activates a skill by reading the listed `SKILL.md` path with the View tool. There is no slash-command skill registry for skills.

Default skill paths on Windows:

- If `KURORYUU_SKILLS_DIR` is set, that directory is used as the default skill directory.
- Otherwise Crush defaults to:
  - `%LOCALAPPDATA%\kuroryuu\skills`
  - `%LOCALAPPDATA%\agents\skills`

Additional directories can be configured in `options.skills_paths`.

The View tool receives `cfg.Options.SkillsPaths`, which allows skill files under those roots to be read without the normal outside-workdir permission flow.

## Custom Commands and MCP Findings

Custom markdown commands are separate from skills.

`internal/commands/commands.go` loads custom `.md` commands from:

- user config commands directory under the Kuroryuu config home,
- `~/.kuroryuu/commands`,
- the project data directory command folder, usually `.kuroryuu/commands`.

MCP support is config-driven through the `mcp` key in `internal/config/config.go`. Crush supports `stdio`, `sse`, and `http` MCP transports. `internal/agent/tools/mcp/init.go` initializes configured servers, lists their tools/prompts/resources, and `internal/agent/tools/mcp-tools.go` wraps MCP tools into agent tools named `mcp_<server>_<tool>`.

This integration does not need MCP. The watcher should be a direct in-process app feature.

## Maildir Polling Contract

The watcher should use the existing LiteHarness inbox layout:

```text
~/.liteharness/mailboxes/<agent-id>/new/
~/.liteharness/mailboxes/<agent-id>/cur/
```

Polling behavior:

1. Resolve the agent id for the running Crush session.
2. Poll `new/` every 2 seconds.
3. For each new message file, read the JSON payload.
4. Extract the message body text.
5. Inject the body into the active Bubble Tea session through an internal `tea.Msg` path.
6. Move the processed file from `new/` to `cur/` after the message has been accepted for delivery.
7. Log parse, read, move, and injection failures visibly enough for debugging.

The watcher should process files in stable order, such as filename sort order, to preserve message order as much as possible.

Expected message payload fields include at least:

```json
{
  "from": "sender-agent-id",
  "type": "notification",
  "body": "message text"
}
```

The implementation should tolerate extra fields and should skip or quarantine malformed files rather than crashing the TUI.

## Proposed Go Design

Add a small LiteHarness package under a Crush-owned internal path:

```text
internal/kuroryuu/liteharness/
  watcher.go
```

Suggested responsibilities:

- `watcher.go`: resolve inbox paths, list `new/`, parse JSON, move files to `cur/`, run a context-bound polling loop, and publish received bodies to a callback.

Use `context.Context` so the watcher exits cleanly when Crush exits.

Pseudo-shape:

```go
type Message struct {
    From string `json:"from"`
    Type string `json:"type"`
    Body string `json:"body"`
}

type Watcher struct {
    AgentID string
    Root string
    Interval time.Duration
    Log *slog.Logger
}

type Deliver func(context.Context, Message) error

func (w Watcher) Run(ctx context.Context, deliver Deliver) error
```

The watcher should not import UI code. It should only call a delivery callback.

## Bubble Tea Injection Design

Add an exported or cross-package message type for LiteHarness delivery. The current `sendMessageMsg` in `internal/ui/model/ui.go` is unexported, so there are two reasonable options:

1. Define a new exported UI message type, such as `ui.LiteHarnessMessageMsg`, and handle it in `UI.Update` by calling `m.sendMessage(msg.Body)`.
2. Define an app-level event type in a shared internal package and let the UI model handle that type.

The watcher should start near interactive program setup in `internal/cmd/root.go` or app initialization, after the `App` exists and before or after `app.Subscribe(program)` is started.

Preferred flow:

```text
root.go creates App and tea.Program
root.go starts app.Subscribe(program)
root.go starts liteharness watcher goroutine with app context
watcher deliver callback sends ui/app LiteHarness message through app.events or program.Send
UI receives message and calls m.sendMessage(body)
```

`app.events` remains internal to `App`, so add a narrow app method such as:

```go
func (app *App) Send(ctx context.Context, msg tea.Msg) error {
    // enqueue msg or return ctx/timeout error
}
```

That keeps all external event injection consistent with the existing `Subscribe(program)` path.

## Agent ID Resolution

The watcher needs a stable agent id. Preferred options, in order:

The initial implementation uses an explicit environment variable:

- `LITEHARNESS_AGENT_ID`: enables the watcher and selects the mailbox id.
- `LITEHARNESS_ROOT`: optional override for the LiteHarness root directory. Defaults to `~/.liteharness`.

If `LITEHARNESS_AGENT_ID` is not set, the watcher is disabled and Crush behaves as it did before.

For interoperability with Sentinel and other agents, the selected id should be printed or logged at startup and documented in any future Crush LiteHarness skill instructions.

## Why Terminal Injection Should Be Fallback Only

External terminal automation is less reliable for Crush because it depends on focus and TUI state.

Important caveats found in the TUI:

- `internal/ui/model/keys.go` binds Enter to sending the textarea contents.
- `internal/ui/model/ui.go` handles `tea.PasteMsg` specially.
- `handlePasteMsg` converts pastes with more than 10 newlines into text attachments instead of prompt text.
- `internal/cmd/root.go` `MaybePrependStdin(...)` only affects startup or non-interactive run mode, not an already-running TUI session.

Because LiteHarness messages can be multi-line, external paste could silently become an attachment or land in the wrong UI state. In-process `tea.Msg` delivery is the correct integration path.

## Implementation Steps

1. Add `internal/kuroryuu/liteharness` with maildir polling and message parsing.
2. Add `ui.LiteHarnessMessageMsg` for inbound LiteHarness messages.
3. Add `app.Send(ctx, msg)` as a narrow app event enqueue method.
4. Start the watcher from interactive startup when `LITEHARNESS_AGENT_ID` is set.
5. On inbound message, deliver the body to `m.sendMessage(...)` through Bubble Tea update handling.
6. Move processed files from `new/` to `cur/` after successful enqueue into the Bubble Tea event path.
7. Quarantine malformed and empty messages by moving them to `cur/` after logging.
8. Add focused tests for maildir parsing, file movement, failed delivery retry, malformed message quarantine, and env gating.
9. Add a manual smoke test: send a LiteHarness message to the Crush agent id while the TUI is running and verify the active session receives it as a prompt.

## Open Decisions

- Whether to add config-file support in addition to the current environment-variable opt-in.
- Whether to persist a generated agent id if no explicit id is supplied.
- Whether inbound messages should be wrapped with sender metadata before sending to the agent, for example `From <id>: <body>`.
- Whether malformed messages should eventually move to a dedicated error directory instead of `cur/`.
