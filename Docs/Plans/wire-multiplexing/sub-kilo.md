# Sub-Plan: kilo TypeScript/SolidJS Multiplexing

**Master:** [master.md](master.md)
**Workstream:** WS2 — kilo
**Dependencies:** None (parallel with opencode)
**Phase:** 1
**Repo:** `E:\SAS\REPO_CLONES\kilo`
**Reference:** `E:\SAS\REPO_CLONES\litecode-v2` (conceptual reference — different stack)

## Context

kilo's TUI is built on SolidJS + @opentui/core + @opentui/solid. The session route
at `packages/opencode/src/cli/cmd/tui/routes/session/index.tsx` is the main UI entry point.

Existing building blocks already in kilo:
- `bun-pty: 0.4.8` — PTY spawning
- `SplitBorder` component at `@tui/component/border`
- `useKeyboard`, `useRenderer`, `useTerminalDimensions` from `@opentui/solid`
- `useKeybind` from `@tui/context/keybind` — the right extension point for tab keybindings
- `BoxRenderable`, `ScrollBoxRenderable` from `@opentui/core` — layout primitives

**Before implementing:** Read the following in kilo:
1. `packages/opencode/src/cli/cmd/tui/routes/session/index.tsx` — full session layout
2. `packages/opencode/src/cli/cmd/tui/context/keybind.ts` — existing keybind system
3. `packages/opencode/src/cli/cmd/tui/app.tsx` — top-level TUI
4. `packages/opencode/src/cli/cmd/tui/component/border.tsx` — SplitBorder usage

**From litecode-v2, read conceptually (for behavior/logic, not syntax):**
- `internal/ui/model/tabs.go` — tab key dispatch logic
- `internal/ui/model/pane_instance.go` — pane abstraction
- `internal/ui/tabsidebar/tabsidebar.go` — sidebar rendering approach

## Tasks

### T1: Read bun-pty API
- Read `node_modules/bun-pty/index.d.ts` (or equivalent types file) to understand the API
- Specifically: how to spawn a PTY, write to it, read output, resize, close
- Document the key API calls before writing any PTY code — this prevents implementation drift
- **Files:** (read-only scouting task — no files created)

### T2: Create TypeScript tab manager
- Create `packages/opencode/src/cli/cmd/tui/tabmgr/index.ts`
- A SolidJS-friendly tab manager using signals/stores
- Mirrors litecode-v2's tabmgr logic in TypeScript:
  - `createTabManager()` — factory returning a reactive store
  - Tab CRUD: `addTab(type: 'session' | 'shell', cwd: string)`, `closeTab(id)`, `switchTab(id)`
  - Pane tree per tab: `splitH(paneId)`, `splitV(paneId)`, `closePane(paneId)`
  - Focus tracking: `focusNext()`, `focusPrev()`
  - Persistence: `saveLayout()`, `loadLayout()` — sessions.json
- Use SolidJS `createStore` for reactivity so components re-render on tab changes
- **Files:** `packages/opencode/src/cli/cmd/tui/tabmgr/index.ts`

### T3: Create PTY pane component
- Create `packages/opencode/src/cli/cmd/tui/component/pane-pty.tsx`
- Uses bun-pty to spawn a shell (PowerShell on Windows, `$SHELL || bash` on Unix)
- Receives: `{ width: number, height: number, focused: boolean, pty: BunPty }`
- Renders PTY output into a `BoxRenderable` / direct @opentui/core render call
- Handles: keyboard input forwarding when focused, resize (resize PTY on dimension change),
  graceful close
- Reference: litecode-v2's `pane_shell.go` for the behavioral model
- **Files:** `packages/opencode/src/cli/cmd/tui/component/pane-pty.tsx`
- **Depends on:** T1

### T4: Create tab sidebar component
- Create `packages/opencode/src/cli/cmd/tui/component/tab-sidebar.tsx`
- Vertical sidebar showing: tab list with icons (session = chat icon, shell = terminal icon),
  current tab highlighted, scrollable if many tabs
- Reference: litecode-v2's `tabsidebar/tabsidebar.go` for rendering logic; adapt to
  @opentui/solid JSX patterns (same visual behavior, different syntax)
- Props: `{ tabManager: TabManager, visible: boolean, width: number, height: number }`
- **Files:** `packages/opencode/src/cli/cmd/tui/component/tab-sidebar.tsx`
- **Depends on:** T2

### T5: Extend keybind context with tab/pane bindings
- Read `packages/opencode/src/cli/cmd/tui/context/keybind.ts` to understand the existing
  keybind system
- Extend it (or create a companion `tab-keybinds.ts`) with the same alt+* bindings as
  litecode-v2:
  - alt+t: new session tab
  - alt+shift+t: new shell tab
  - alt+w: close tab
  - alt+]: next tab
  - alt+[: prev tab
  - alt+shift+]: move tab right
  - alt+shift+[: move tab left
  - alt+shift+h: split horizontal
  - alt+shift+v: split vertical
  - alt+shift+w: close pane
  - alt+right / alt+down: focus next pane
  - alt+left / alt+up: focus prev pane
  - alt+b: toggle sidebar
- Verify these don't conflict with existing kilo keybindings (read current keybind.ts)
- **Files:** `packages/opencode/src/cli/cmd/tui/context/keybind.ts` (extended) or
  `packages/opencode/src/cli/cmd/tui/context/tab-keybinds.ts` (new)
- **Depends on:** T2

### T6: Create pane layout component
- Create `packages/opencode/src/cli/cmd/tui/component/pane-layout.tsx`
- Given a split tree from the tab manager, compute pane rectangles (use the same
  geometry logic as litecode-v2's `internal/split/layout.go`, ported to TypeScript)
- For each leaf pane: render either the session view (existing `index.tsx` content wrapped
  in a pane slot) or the PTY pane component (T3)
- Renders dividers between panes (can use SplitBorder from existing border component)
- Props: `{ tab: Tab, width: number, height: number, focusedPaneId: string }`
- **Files:** `packages/opencode/src/cli/cmd/tui/component/pane-layout.tsx`
- **Depends on:** T3, T4

### T7: Wire into app.tsx and session route
- In `packages/opencode/src/cli/cmd/tui/app.tsx`:
  - Initialize tab manager (load from sessions.json via T2's loadLayout)
  - Pass tab manager via SolidJS context or props
- In `packages/opencode/src/cli/cmd/tui/routes/session/index.tsx`:
  - Wrap existing session content in a pane slot (it becomes the "session pane" content)
  - Add tab sidebar (T4) on the left side (match existing app chrome — sidebar position)
  - Add pane layout (T6) as the main content area
  - Wire keyboard handler from T5 to tab manager actions
- **Files:** `packages/opencode/src/cli/cmd/tui/app.tsx`,
  `packages/opencode/src/cli/cmd/tui/routes/session/index.tsx`
- **Depends on:** T5, T6

### T8: Validate
- `bun typecheck` — zero errors
- `bun test` — no regressions
- Manual smoke: launch kilo TUI, create session tab (alt+t), create shell tab (alt+shift+t),
  split pane (alt+shift+v), navigate (alt+right/left), close tab (alt+w)
- **Depends on:** T7

## Validation Commands
```bash
cd E:\SAS\REPO_CLONES\kilo\packages\opencode
bun typecheck
bun test
```

## Acceptance Criteria
- TypeScript compiles with zero errors
- All existing bun tests pass
- Tab sidebar renders correctly with session/shell tab types distinguished
- Pane layout renders split panes with dividers
- PTY shell opens PowerShell on Windows, bash (or $SHELL) on Unix
- Keybindings work and don't conflict with existing kilo bindings
- Sessions.json persists and restores layout
