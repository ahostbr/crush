# Sub-Plan: opencode Go/BubbleTea Multiplexing Port

**Master:** [master.md](master.md)
**Workstream:** WS1 — opencode
**Dependencies:** None (parallel with kilo)
**Phase:** 1
**Repo:** `E:\SAS\REPO_CLONES\opencode`
**Reference:** `E:\SAS\REPO_CLONES\litecode-v2` (read-only)

## Context

opencode uses a `pages map[PageID]tea.Model` architecture in `internal/tui/tui.go`. The
existing infrastructure packages are already in place and compile clean:
- `internal/split/` — binary split tree data layer
- `internal/tabmgr/` — tab manager + persistence
- `internal/pty/` — ConPTY (Windows) + creack/pty (Unix)

What's missing is all the UI layer: gitwatcher, tabsidebar, pane types, keybindings, and
rendering. litecode-v2 has all of this working — read it first.

**Critical architectural decision:** opencode's `appModel` is page-based (one full-screen
page at a time). litecode-v2's `UI` is pane-based (split tree of panes). Before writing any
code, decide on the integration approach:

Option A — Wrap existing chat page: The `ChatPage` becomes pane-aware. Multiple `ChatPage`
instances can be hosted in a split grid. `appModel` gains a `TabManager` and routes input
to the focused pane's page.

Option B — New MultiplexPage: Add a new `MultiplexPage` (alongside `ChatPage`, `LogsPage`)
that owns the split tree and hosts pane instances internally. `appModel` navigates to
`MultiplexPage` on startup (instead of `ChatPage`).

Read `internal/tui/tui.go`, `internal/tui/page/page.go`, and `internal/tui/page/chat.go`
before committing to an approach. Recommend Option B for cleaner separation.

## Tasks

### T1: Port gitwatcher
- Copy `E:\SAS\REPO_CLONES\litecode-v2\internal\gitwatcher\watcher.go` to
  `E:\SAS\REPO_CLONES\opencode\internal\gitwatcher\watcher.go`
- Copy `watcher_test.go` alongside
- Replace all `github.com/ahostbr/litecode-v2` imports with `github.com/opencode-ai/opencode`
- Run `go test ./internal/gitwatcher/...` — must pass
- **Files:** `internal/gitwatcher/watcher.go`, `internal/gitwatcher/watcher_test.go`

### T2: Port tabsidebar
- Copy `E:\SAS\REPO_CLONES\litecode-v2\internal\ui\tabsidebar\tabsidebar.go` to
  `E:\SAS\REPO_CLONES\opencode\internal\tui\tabsidebar\tabsidebar.go`
- Replace module path in imports
- Verify it compiles: `go build ./internal/tui/tabsidebar/...`
- **Files:** `internal/tui/tabsidebar/tabsidebar.go`
- **Depends on:** T1

### T3: Add Tabs keybindings to appModel keymap
- In `internal/tui/tui.go`, extend the `keyMap` struct with a `Tabs` section mirroring
  litecode-v2's `internal/ui/model/keys.go` Tabs section exactly
- Add DefaultTabKeys() initialization matching litecode-v2's alt+* bindings
- Verify no conflicts with existing keys (ctrl+l, ctrl+c, ctrl+_, ctrl+s, ctrl+k, ctrl+f,
  ctrl+o, ctrl+t, ctrl+n, ctrl+d — all ctrl+, so alt+* is safe)
- **Files:** `internal/tui/tui.go`

### T4: Create PaneInstance type
- Create `internal/tui/pane_instance.go`
- Reference: `E:\SAS\REPO_CLONES\litecode-v2\internal\ui\model\pane_instance.go`
- Adapt to opencode's imports (charmbracelet/bubbletea not charm.land/bubbletea/v2)
- PaneInstance holds: ID, Type (PaneSession/PaneShell), chat page state, shell state, CWD,
  Focused, Ready
- **Files:** `internal/tui/pane_instance.go`

### T5: Create ShellPane (PTY pane renderer)
- Create `internal/tui/pane_shell.go`
- Reference: `E:\SAS\REPO_CLONES\litecode-v2\internal\ui\model\pane_shell.go`
- Adapts litecode-v2's ShellScreen to opencode's BubbleTea version
- Spawns PTY via `internal/pty` package
- Shell: PowerShell on Windows (`pwsh` or `powershell`), `$SHELL` or `bash` on Unix
- Renders PTY output as lipgloss-styled string into its allocated rectangle
- **Files:** `internal/tui/pane_shell.go`
- **Depends on:** T4

### T6: Create MultiplexPage (integration point)
- Create `internal/tui/page/multiplex.go` — a new tea.Model implementing `page.PageID`
- This page owns: `tabManager *tabmgr.TabManager`, `panes map[string]*PaneInstance`,
  `gitWatcher *gitwatcher.Watcher`
- `Init()`: initialize TabManager from sessions.json (use tabmgr persistence), create
  default tab with one session pane
- `Update()`: route key events through tab key handler (port from litecode-v2 `tabs.go`)
- `View()`: query active tab's SplitTree → layout rectangles → render each pane into its
  rectangle → composite with tab sidebar
- Read litecode-v2's `internal/ui/model/ui.go` Draw() section and `tabs.go` handleTabKeys()
  as the primary reference for this file
- **Files:** `internal/tui/page/multiplex.go`
- **Depends on:** T2, T3, T4, T5

### T7: Register MultiplexPage in appModel
- In `internal/tui/tui.go`, add `page.MultiplexPageID` to the pages map
- On startup: navigate to MultiplexPage instead of (or in addition to) ChatPage
- Route global key events: tab keys go to MultiplexPage, other keys go to currentPage as before
- **Files:** `internal/tui/tui.go`
- **Depends on:** T6

### T8: Wire persistent layout
- In MultiplexPage.Init(), load layout from sessions.json via `tabmgr.LoadLayout()`
- On each tab/pane change, call `tabmgr.SaveLayout()` (debounced, same pattern as litecode-v2)
- **Files:** `internal/tui/page/multiplex.go` (extends T6)
- **Depends on:** T6

### T9: Validate
- `go build ./...` — clean build
- `go test ./...` — no regressions
- Manual smoke: launch opencode, create session tab (alt+t), split vertical (alt+shift+v),
  create shell tab (alt+shift+t), navigate panes (alt+right/left), close tab (alt+w)
- **Depends on:** T7, T8

## Validation Commands
```bash
cd E:\SAS\REPO_CLONES\opencode
go build ./...
go test ./...
go test ./internal/gitwatcher/...
go test ./internal/tui/...
```

## Acceptance Criteria
- Clean build
- All existing tests pass
- New gitwatcher tests pass
- Tab/pane navigation works at runtime
- PTY shell opens PowerShell on Windows, bash on Unix
- Layout persists to sessions.json and restores on restart
