#!/bin/bash
# =============================================================================
# Kuroryuu Rebranding Script
# Converts Crush (charmbracelet) branding to Kuroryuu (ahostbr)
# =============================================================================
set -eu

REPO="/mnt/e/SAS/REPO_CLONES/crush"
cd "$REPO"

# Verify we're in the right repo
if [ ! -f "go.mod" ] || ! grep -q "github.com/ahostbr/crush" go.mod; then
    echo "ERROR: Not in the expected crush repo or already rebranded"
    exit 1
fi

echo "=== Kuroryuu Rebranding Script ==="
echo "Working directory: $(pwd)"
echo ""

# Helper: find text files, excluding .git, go.sum, binaries
find_text() {
    find . \
        -not -path "./.git/*" \
        -not -name "go.sum" \
        -not -name "*.exe" \
        -not -name "*.db" \
        -not -name "*.png" \
        -not -name "*.jpg" \
        -not -name "*.gif" \
        -not -name "*.ico" \
        -not -name "*.woff" \
        -not -name "*.woff2" \
        -type f \
        "$@"
}

# ===========================================================================
# PHASE 1a: Go module path (MOST SPECIFIC — must come first)
# ===========================================================================
echo "Phase 1a: Go module path (github.com/ahostbr/crush -> github.com/ahostbr/crush)"
find_text \( -name "*.go" -o -name "go.mod" -o -name "*.yml" -o -name "*.yaml" -o -name "*.json" -o -name "*.md" -o -name "*.tpl" -o -name "*.sh" \) \
    -exec sed -i 's|github\.com/charmbracelet/crush|github.com/ahostbr/crush|g' {} +
echo "  Done."

# ===========================================================================
# PHASE 1b: Product-specific URLs (before general domain changes)
# ===========================================================================
echo "Phase 1b: Product-specific URLs"

# Schema URL
find_text \( -name "*.go" -o -name "*.json" -o -name "*.md" -o -name "*.tpl" \) \
    -exec sed -i 's|https://charm\.land/crush\.json|https://kuroryuu.com/kuroryuu.json|g' {} +

# Email
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" -o -name "*.yml" -o -name "*.yaml" \) \
    -exec sed -i 's|crush@charm\.land|cli@kuroryuu.com|g' {} +

# Telemetry endpoint
find_text \( -name "*.go" \) \
    -exec sed -i 's|https://data\.charm\.land|https://data.kuroryuu.com|g' {} +

# Catwalk endpoint
find_text \( -name "*.go" \) \
    -exec sed -i 's|https://catwalk\.charm\.sh|https://catwalk.kuroryuu.com|g' {} +

# Hyper endpoint
find_text \( -name "*.go" -o -name "*.json" \) \
    -exec sed -i 's|https://hyper\.charm\.land|https://hyper.kuroryuu.com|g' {} +

# charm.sh/crush homepage
find_text \( -name "*.go" -o -name "*.yml" -o -name "*.yaml" -o -name "*.md" \) \
    -exec sed -i 's|https://charm\.sh/crush|https://kuroryuu.com|g' {} +
find_text \( -name "*.go" -o -name "*.yml" -o -name "*.yaml" -o -name "*.md" \) \
    -exec sed -i 's|https://charm\.land|https://kuroryuu.com|g' {} +

# charm.sh specific URLs (repo.charm.sh etc) in docs only
find_text \( -name "*.md" \) \
    -exec sed -i 's|https://repo\.charm\.sh|https://repo.kuroryuu.com|g' {} +
find_text \( -name "*.md" \) \
    -exec sed -i 's|https://stuff\.charm\.sh|https://stuff.kuroryuu.com|g' {} +

# GitHub releases URL for update checker
find_text \( -name "*.go" \) \
    -exec sed -i 's|api\.github\.com/repos/charmbracelet/crush|api.github.com/repos/ahostbr/crush|g' {} +

echo "  Done."

# ===========================================================================
# PHASE 1c: Environment variable prefix KURORYUU_ -> KURORYUU_
# ===========================================================================
echo "Phase 1c: Environment variable prefix"
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" -o -name "*.yml" -o -name "*.yaml" -o -name "*.sh" -o -name "*.json" -o -name "*.env*" -o -name "*.sample" \) \
    -exec sed -i 's/KURORYUU_/KURORYUU_/g' {} +
echo "  Done."

# ===========================================================================
# PHASE 1d: Data/config paths and constants
# ===========================================================================
echo "Phase 1d: Data/config paths and constants"

# .crushignore -> .kuroryuuignore
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" \) \
    -exec sed -i 's/\.crushignore/.kuroryuuignore/g' {} +

# appName constant
find_text -name "*.go" \
    -exec sed -i 's/appName              = "crush"/appName              = "kuroryuu"/g' {} +

# defaultDataDirectory constant
find_text -name "*.go" \
    -exec sed -i 's/defaultDataDirectory = "\.crush"/defaultDataDirectory = ".kuroryuu"/g' {} +

# .crush as directory path in strings (careful with quotes and slashes)
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" -o -name "*.yml" -o -name "*.yaml" \) \
    -exec sed -i 's|/\.crush/|/.kuroryuu/|g' {} +
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" \) \
    -exec sed -i 's|"\.crush"|".kuroryuu"|g' {} +

# Context path strings in Go code
find_text -name "*.go" \
    -exec sed -i 's/"crush\.md"/"kuroryuu.md"/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/"crush\.local\.md"/"kuroryuu.local.md"/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/"Crush\.md"/"Kuroryuu.md"/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/"Crush\.local\.md"/"Kuroryuu.local.md"/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/"CRUSH\.md"/"KURORYUU.md"/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/"CRUSH\.local\.md"/"KURORYUU.local.md"/g' {} +

# jsonschema examples
find_text -name "*.go" \
    -exec sed -i 's/example=CRUSH\.md/example=KURORYUU.md/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/example=\.crush/example=.kuroryuu/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/default=\.crush/default=.kuroryuu/g' {} +

echo "  Done."

# ===========================================================================
# PHASE 1e: User-Agent strings
# ===========================================================================
echo "Phase 1e: User-Agent strings"
find_text -name "*.go" \
    -exec sed -i 's|"crush/1\.0"|"kuroryuu/1.0"|g' {} +
echo "  Done."

# ===========================================================================
# PHASE 1f: CLI/binary name
# ===========================================================================
echo "Phase 1f: CLI/binary name"
find_text -name "*.go" \
    -exec sed -i 's/Use:   "crush"/Use:   "Kuroryuu-cli-v0.5"/g' {} +
echo "  Done."

# ===========================================================================
# PHASE 1g: Attribution strings
# ===========================================================================
echo "Phase 1g: Attribution strings"
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" \) \
    -exec sed -i 's/Generated with Crush/Generated with Kuroryuu/g' {} +
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" \) \
    -exec sed -i 's/via Crush/via Kuroryuu/g' {} +
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" \) \
    -exec sed -i 's/Co-Authored-By: Crush/Co-Authored-By: Kuroryuu/g' {} +
echo "  Done."

# ===========================================================================
# PHASE 1h: Template/prompt branding
# ===========================================================================
echo "Phase 1h: Template and prompt branding"
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" \) \
    -exec sed -i 's/You are Crush/You are Kuroryuu/g' {} +
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" \) \
    -exec sed -i 's/agent for Crush/agent for Kuroryuu/g' {} +
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" \) \
    -exec sed -i 's/Crush crashed/Kuroryuu crashed/g' {} +
find_text \( -name "*.go" -o -name "*.md" -o -name "*.tpl" \) \
    -exec sed -i 's/Crush Usage Statistics/Kuroryuu Usage Statistics/g' {} +
echo "  Done."

# ===========================================================================
# PHASE 1i: Build config (.goreleaser.yml, Taskfile.yaml)
# ===========================================================================
echo "Phase 1i: Build configuration"

# goreleaser: project_name
sed -i 's/^project_name: crush$/project_name: kuroryuu/' .goreleaser.yml

# goreleaser: archive/binary name templates
sed -i 's/crush_/kuroryuu_/g' .goreleaser.yml

# goreleaser: completion file names
sed -i 's/crush\.bash/kuroryuu.bash/g' .goreleaser.yml
sed -i 's/crush\.zsh/kuroryuu.zsh/g' .goreleaser.yml
sed -i 's/crush\.fish/kuroryuu.fish/g' .goreleaser.yml
sed -i 's/crush\.1\.gz/kuroryuu.1.gz/g' .goreleaser.yml

# goreleaser: AUR/nfpm package paths
sed -i 's|/doc/crush/|/doc/kuroryuu/|g' .goreleaser.yml
sed -i 's|/licenses/crush/|/licenses/kuroryuu/|g' .goreleaser.yml
sed -i 's|crush\.git|kuroryuu.git|g' .goreleaser.yml
sed -i 's|crush-bin\.git|kuroryuu-bin.git|g' .goreleaser.yml

# goreleaser: npm scope
sed -i 's|@charmland/crush|@ahostbr/kuroryuu|g' .goreleaser.yml

# goreleaser: hook commands that reference ./crush binary
sed -i 's|\./crush completion|./kuroryuu completion|g' .goreleaser.yml
sed -i 's|\./crush man|./kuroryuu man|g' .goreleaser.yml
sed -i 's|"go run \. completion bash >./completions/crush|"go run . completion bash >./completions/kuroryuu|g' .goreleaser.yml
sed -i 's|"go run \. completion zsh >./completions/crush|"go run . completion zsh >./completions/kuroryuu|g' .goreleaser.yml
sed -i 's|"go run \. completion fish >./completions/crush|"go run . completion fish >./completions/kuroryuu|g' .goreleaser.yml
sed -i "s|go run \. completion bash >./completions/crush|go run . completion bash >./completions/kuroryuu|g" .goreleaser.yml
sed -i "s|go run \. completion zsh >./completions/crush|go run . completion zsh >./completions/kuroryuu|g" .goreleaser.yml
sed -i "s|go run \. completion fish >./completions/crush|go run . completion fish >./completions/kuroryuu|g" .goreleaser.yml

# goreleaser: install paths with /crush" at end
sed -i 's|/crush"|/kuroryuu"|g' .goreleaser.yml
sed -i 's|_crush"|_kuroryuu"|g' .goreleaser.yml

# goreleaser: binary install reference
sed -i 's|"./crush"|"./kuroryuu"|g' .goreleaser.yml

# goreleaser: commit author and publisher
sed -i 's/name: "Charm"/name: "Kuroryuu"/g' .goreleaser.yml
sed -i 's/publisher: charmbracelet/publisher: ahostbr/g' .goreleaser.yml
sed -i 's/copyright: Charmbracelet, Inc/copyright: ahostbr/g' .goreleaser.yml
sed -i 's|charmcli@users\.noreply\.github\.com|cli@kuroryuu.com|g' .goreleaser.yml

# goreleaser: maintainer emails
sed -i 's|kujtimiihoxha <kujtim@charm\.sh>|maintainer@kuroryuu.com|g' .goreleaser.yml
sed -i 's|caarlos0 <carlos@charm\.sh>.*|cli@kuroryuu.com|g' .goreleaser.yml

# goreleaser: winget
sed -i 's|publisher_url: https://charm\.land|publisher_url: https://kuroryuu.com|g' .goreleaser.yml

# goreleaser: remaining "crush" as binary name in package sections
sed -i 's|provides:\n      - crush|provides:\n      - kuroryuu|g' .goreleaser.yml
sed -i 's|conflicts:\n      - crush|conflicts:\n      - kuroryuu|g' .goreleaser.yml

# goreleaser: charmbracelet org in repo owners (for homebrew, scoop, nur, winget)
sed -i 's/owner: charmbracelet/owner: ahostbr/g' .goreleaser.yml
sed -i 's/owner: "charmbracelet"/owner: "ahostbr"/g' .goreleaser.yml

# goreleaser: license
sed -i 's/license: fsl11Mit/license: agpl3Only/g' .goreleaser.yml

# Taskfile: binary references
sed -i 's|crush{{exeExt}}|kuroryuu{{exeExt}}|g' Taskfile.yaml
sed -i 's|\./crush|./kuroryuu|g' Taskfile.yaml

echo "  Done."

# ===========================================================================
# PHASE 1j: CI/GitHub files
# ===========================================================================
echo "Phase 1j: CI/GitHub files"

# CODEOWNERS
find_text -name "CODEOWNERS" \
    -exec sed -i 's/@charmbracelet/@ahostbr/g' {} +

# CLA workflow and signatures
find_text -path "./.github/*" \( -name "*.yml" -o -name "*.yaml" -o -name "*.json" \) \
    -exec sed -i 's/charmbracelet\/crush/ahostbr\/crush/g' {} +
find_text -path "./.github/*" \( -name "*.yml" -o -name "*.yaml" \) \
    -exec sed -i 's/charmcli@users\.noreply\.github\.com/cli@kuroryuu.com/g' {} +

echo "  Done."

# ===========================================================================
# PHASE 1k: Logo/UI branding
# ===========================================================================
echo "Phase 1k: Logo and UI branding"

# Trademark text
find_text -name "logo.go" \
    -exec sed -i 's/Charm™/Kuroryuu/g' {} +

# "Crush" as brand string in SmallRender
find_text -name "logo.go" \
    -exec sed -i 's/"Crush"/"Kuroryuu"/g' {} +

# Comment references in logo
find_text -name "logo.go" \
    -exec sed -i 's/Crush wordmark/Kuroryuu wordmark/g' {} +
find_text -name "logo.go" \
    -exec sed -i 's/Crush logo/Kuroryuu logo/g' {} +
find_text -name "logo.go" \
    -exec sed -i 's/Crush title/Kuroryuu title/g' {} +

# Chroma theme name (used for syntax highlighting)
find_text -name "*.go" \
    -exec sed -i 's/"crush"/"kuroryuu"/g' {} +

# Style field names referencing Charm
find_text -name "styles.go" \
    -exec sed -i 's/LogoCharmColor/LogoBrandColor/g' {} +

echo "  Done."
echo "  WARNING: internal/ui/logo/logo.go letterforms spell C-R-U-S-H and need manual rework for KURORYUU"

# ===========================================================================
# PHASE 1l: Telemetry
# ===========================================================================
echo "Phase 1l: Telemetry"
# Disable PostHog key (it points to Charm's analytics)
find_text -name "event.go" -path "*/event/*" \
    -exec sed -i 's/phc_4zt4VgDWLqbYnJYEwLRxFoaTL2noNrQij0C6E8k3I0V/DISABLED_REPLACE_WITH_YOUR_OWN_KEY/g' {} +
echo "  Done."

# ===========================================================================
# PHASE 1m: General prose replacements in docs (LAST)
# ===========================================================================
echo "Phase 1m: General prose in docs"

# README.md — Crush as product name
sed -i 's/\bCrush\b/Kuroryuu/g' README.md 2>/dev/null || true
sed -i 's/\bcrush\b/kuroryuu/g' README.md 2>/dev/null || true

# Social/community links (Charm-specific)
sed -i 's|twitter\.com/charmcli|kuroryuu.com|g' README.md 2>/dev/null || true
sed -i 's|mastodon\.social/@charmcli|kuroryuu.com|g' README.md 2>/dev/null || true
sed -i 's|bsky\.app/profile/charm\.land|kuroryuu.com|g' README.md 2>/dev/null || true

# CLI description strings in Go code
find_text -name "root.go" -path "*/cmd/*" \
    -exec sed -i 's/crush configuration/kuroryuu configuration/g' {} +
find_text -name "schema.go" -path "*/cmd/*" \
    -exec sed -i 's/crush configuration/kuroryuu configuration/g' {} +

# Remaining "crush" in Go source as string literals (careful — only quoted strings)
# This catches things like Short: "... crush ..." descriptions
find_text -name "*.go" \
    -exec sed -i 's/crush logs/kuroryuu logs/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/crush run/kuroryuu run/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/crush update-providers/kuroryuu update-providers/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/crush schema/kuroryuu schema/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/crush dirs/kuroryuu dirs/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/crush projects/kuroryuu projects/g' {} +

# Go function names with "Crush" in them
find_text -name "*.go" \
    -exec sed -i 's/createDotCrushDir/createDotKuroryuuDir/g' {} +
find_text -name "*.go" \
    -exec sed -i 's/PushPopCrushEnv/PushPopKuroryuuEnv/g' {} +

echo "  Done."

# ===========================================================================
# PHASE 2: File renames
# ===========================================================================
echo ""
echo "Phase 2: File renames"

if [ -f "crush.json" ]; then
    mv crush.json kuroryuu.json
    echo "  Renamed crush.json -> kuroryuu.json"
fi

if [ -f ".crush.json" ]; then
    mv .crush.json .kuroryuu.json
    echo "  Renamed .crush.json -> .kuroryuu.json"
fi

echo "  Done."

# ===========================================================================
# PHASE 3: Verification
# ===========================================================================
echo ""
echo "Phase 3: Verification"
echo ""

# Ensure module line is correct
go mod edit -module github.com/ahostbr/crush 2>/dev/null || echo "  (go mod edit skipped — go not in PATH or error)"

# Delete go.sum for regeneration
rm -f go.sum
echo "  Deleted go.sum (will be regenerated by go mod tidy)"

echo ""
echo "--- Checking upstream charm.land imports are intact ---"
echo "bubbletea imports:"
grep -r "charm.land/bubbletea" --include="*.go" -l 2>/dev/null | head -3 || echo "  NONE FOUND — problem!"
echo "lipgloss imports:"
grep -r "charm.land/lipgloss" --include="*.go" -l 2>/dev/null | head -3 || echo "  NONE FOUND — problem!"
echo "fantasy imports:"
grep -r "charm.land/fantasy" --include="*.go" -l 2>/dev/null | head -3 || echo "  NONE FOUND — problem!"
echo "catwalk imports:"
grep -r "charm.land/catwalk" --include="*.go" -l 2>/dev/null | head -3 || echo "  NONE FOUND — problem!"

echo ""
echo "--- Checking for remaining charmbracelet/crush references ---"
REMAINING=$(grep -rn "charmbracelet/crush" --include="*.go" --include="*.mod" 2>/dev/null | wc -l)
echo "  Remaining charmbracelet/crush in Go files: $REMAINING"
if [ "$REMAINING" -gt 0 ]; then
    grep -rn "charmbracelet/crush" --include="*.go" --include="*.mod" 2>/dev/null | head -10
fi

echo ""
echo "--- Remaining 'crush' occurrences for manual review ---"
grep -rn '\bcrush\b' --include="*.go" --include="*.md" --include="*.yml" --include="*.yaml" --include="*.json" --include="*.tpl" 2>/dev/null \
    | grep -v "charmbracelet" \
    | grep -v "charm.land" \
    | grep -v "rebrand.sh" \
    | grep -v "go.sum" \
    | head -30 || echo "  None found!"

echo ""
echo "=== Rebranding complete! ==="
echo ""
echo "Manual follow-up items:"
echo "  1. internal/ui/logo/logo.go — letterforms need manual rework for KURORYUU"
echo "  2. internal/event/event.go — PostHog key disabled, replace with your own"
echo "  3. Run 'go mod tidy' to regenerate go.sum"
echo "  4. Run 'go build .' to verify compilation"
echo "  5. Run 'go run main.go schema > schema.json' to regenerate schema"
echo "  6. Re-record VCR test cassettes: 'task test:record'"
echo "  7. Review CLA.md and LICENSE.md for legal updates"
echo "  8. Review SVG files in internal/cmd/stats/"
echo "  9. Full docs pass on README.md, AGENTS.md, CLAUDE.md"
