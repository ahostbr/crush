# Plan: Wire Multiplexing into OpenCode + Kilo

## Why This Exists
LiteSuite users across multiple CLIs need terminal multiplexing — tabs, split panes, and PTY
shell access — as first-class features. litecode-v2 already ships a fully working multiplexer.
This plan ports that system to opencode (Go) and kilo (TypeScript/SolidJS) so all three CLIs
in LiteSuite have the capability.

## Task Description
Port the working multiplexing system from `E:\SAS\REPO_CLONES\litecode-v2` (the reference
implementation) into two separate repos:
- **opencode** (`E:\SAS\REPO_CLONES\opencode`) — Go/BubbleTea port
- **kilo** (`E:\SAS\REPO_CLONES\kilo`) — TypeScript/SolidJS/@opentui equivalent

Both workstreams run in parallel. litecode-v2 is read-only — the reference, not a target.

## Objective
Each CLI ships with:
1. **Tab bar** — create (session tab / shell tab), switch, close, reorder
2. **Split panes** — horizontal + vertical splits, navigate between panes
3. **PTY terminal pane** — live interactive shell (PowerShell on Windows, bash on Unix)
4. **Chat/session pane** — existing chat UI hosted in a pane slot (multiple sessions side-by-side)
5. **Persistent layout** — tab/pane layout saved and restored (sessions.json)

## Reference Implementation
litecode-v2 is the authoritative source. Sub-agents MUST read these files before implementing:

| Component | litecode-v2 path |
|-----------|-----------------|
| Tab keybindings | `internal/ui/model/keys.go` — `Tabs` section |
| Tab key dispatch | `internal/ui/model/tabs.go` |
| Pane abstraction | `internal/ui/model/pane_instance.go` |
| PTY shell pane | `internal/ui/model/pane_shell.go` |
| Tab sidebar | `internal/ui/tabsidebar/tabsidebar.go` |
| Git watcher | `internal/gitwatcher/watcher.go` |
| Split/tabmgr/pty | `internal/split/`, `internal/tabmgr/`, `internal/pty/` |
| Main UI model | `internal/ui/model/ui.go` (see how everything wires together) |

## Fact Dependencies

| Fact | Confidence | Workstream | Impact if Wrong |
|------|-----------|------------|-----------------|
| opencode split/tabmgr/pty packages compile cleanly | HIGH | opencode | Low — verified by go build |
| litecode-v2 gitwatcher is needed in opencode | HIGH | opencode | Medium — tab CWD tracking breaks without it |
| kilo's @opentui/core supports BoxRenderable layout patterns | HIGH | kilo | Low — verified in session/index.tsx |
| kilo's useKeybind is the right extension point for tab keys | HIGH | kilo | Medium — wrong hook means keybindings don't work |
| alt+* keybindings don't conflict with opencode's existing ctrl+* keys | HIGH | opencode | Low — opencode uses only ctrl+{key}, litecode-v2 style uses alt+{key} |
| bun-pty v0.4.8 is sufficient for PTY terminal pane in kilo | MED | kilo | Medium — may need API investigation |

## Keybindings Reference (from litecode-v2)

Use these bindings as the default for both opencode and kilo. They use `alt+` prefix
which avoids all conflicts with existing app-level `ctrl+` bindings.

| Action | Key |
|--------|-----|
| New session tab | alt+t |
| New shell tab | alt+shift+t |
| Close tab | alt+w |
| Next tab | alt+] |
| Prev tab | alt+[ |
| Move tab right | alt+shift+] |
| Move tab left | alt+shift+[ |
| Split horizontal | alt+shift+h |
| Split vertical | alt+shift+v |
| Close pane | alt+shift+w |
| Focus next pane | alt+right / alt+down |
| Focus prev pane | alt+left / alt+up |
| Toggle tab sidebar | alt+b |

## Workstreams

| ID | Name | Description | Sub-Plan | Phase |
|----|------|-------------|----------|-------|
| WS1 | opencode | Port litecode-v2 multiplexing to Go/BubbleTea opencode | [sub-opencode.md](sub-opencode.md) | 1 |
| WS2 | kilo | Build TypeScript/SolidJS multiplexing equivalent | [sub-kilo.md](sub-kilo.md) | 1 |

## Orchestration DAG

Phase 1 (parallel): WS1, WS2
Phase 2 (after both): Integration smoke test + review

Both workstreams are fully independent — different repos, different stacks. Run in parallel.

## Acceptance Criteria

**opencode:**
- `go build ./...` passes
- `go test ./...` passes (no new failures)
- Can create a session tab (alt+t), a shell tab (alt+shift+t), split a pane (alt+shift+v),
  navigate between panes (alt+right/left), close a tab (alt+w)
- PTY shell pane renders live output from PowerShell/bash
- Tab layout persists across restarts (sessions.json)

**kilo:**
- `bun typecheck` passes
- Tab sidebar renders with correct current-tab highlighting
- Can create session tab, shell tab, split panes, navigate panes
- PTY shell pane renders live output
- Keybindings don't conflict with existing kilo keybindings

## Validation Commands

**opencode:**
```bash
cd E:\SAS\REPO_CLONES\opencode
go build ./...
go test ./...
```

**kilo:**
```bash
cd E:\SAS\REPO_CLONES\kilo\packages\opencode
bun typecheck
bun test
```

## Remaining Uncertainties

- bun-pty v0.4.8 API surface: sub-agent should read its types before implementing the PTY pane
  component. If the API is insufficient, may need a newer version.
- opencode's `appModel` (page-based) vs litecode-v2's `UI` model (pane-based): the port is
  an adaptation, not a 1:1 copy. The sub-agent needs to decide where `TabManager` and `panes`
  live in the page-based architecture (likely: wrap the chat page to become pane-aware, or
  add a new `MultiplexPage` that hosts panes and delegates to the existing chat page per pane).

## Execution Workflow
1. Worktree (`superpowers:using-git-worktrees`) — isolate before touching code
2. Tests (`superpowers:test-driven-development`) — tests before implementation
3. Implement — dispatch sub-plans to agents via `/max-subagents-parallel`
4. Debug (`superpowers:systematic-debugging`) — when tests fail
5. Verify (`superpowers:verification-before-completion`) — every task verified
6. Review (`superpowers:requesting-code-review`) — before merging
7. Finish (`superpowers:finishing-a-development-branch`) — structured merge/PR/cleanup

## Execution Echo
After implementing this plan, revisit:
- Did the appModel adaptation for opencode work as expected, or was a different integration
  point (e.g., new MultiplexPage) needed?
- Did bun-pty v0.4.8 have the API needed for the kilo PTY pane?
- Were the alt+* keybindings the right choice for kilo's opentui context?
- Was the workstream split (parallel) correct, or did one block the other?

## Notes
- DO NOT modify litecode-v2. It is the reference, not a target.
- Crush is excluded from this plan — its context needs verification before touching it.
- Session repo (where this plan lives): `E:\SAS\REPO_CLONES\crush`
