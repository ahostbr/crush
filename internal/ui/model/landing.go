package model

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ahostbr/crush/internal/agent"
	"github.com/ahostbr/crush/internal/ui/common"
	"github.com/ahostbr/crush/internal/ui/logo"
	"github.com/charmbracelet/ultraviolet/layout"
)

// selectedLargeModel returns the currently selected large language model from
// the agent coordinator, if one exists.
func (m *UI) selectedLargeModel() *agent.Model {
	if m.com.App.AgentCoordinator != nil {
		model := m.com.App.AgentCoordinator.Model()
		return &model
	}
	return nil
}

// landingView renders the landing page view showing the current working
// directory, model information, and LSP/MCP status in a two-column layout.
func (m *UI) landingView() string {
	t := m.com.Styles
	width := m.layout.main.Dx()
	cwd := common.PrettyPath(t, m.com.Config().WorkingDir(), width)

	parts := []string{
		cwd,
	}

	parts = append(parts, "", m.modelInfo(width))
	infoSection := lipgloss.JoinVertical(lipgloss.Left, parts...)

	_, remainingHeightArea := layout.SplitVertical(m.layout.main, layout.Fixed(lipgloss.Height(infoSection)+1))

	mcpLspSectionWidth := min(30, (width-1)/2)

	lspSection := m.lspInfo(mcpLspSectionWidth, max(1, remainingHeightArea.Dy()), false)
	mcpSection := m.mcpInfo(mcpLspSectionWidth, max(1, remainingHeightArea.Dy()), false)

	statusContent := lipgloss.JoinHorizontal(lipgloss.Left, lspSection, " ", mcpSection)

	// Render dragon art if there's enough space
	mainHeight := m.layout.main.Dy() - 1
	usedHeight := lipgloss.Height(infoSection) + 1 + lipgloss.Height(statusContent) + 1
	dragon := logo.DragonArt()
	dragonLines := strings.Split(dragon, "\n")
	dragonHeight := len(dragonLines)
	availableHeight := mainHeight - usedHeight

	var dragonRendered string
	if availableHeight >= 10 && width >= 60 {
		// Trim dragon if it's taller than available space
		if dragonHeight > availableHeight {
			dragonLines = dragonLines[:availableHeight]
			dragon = strings.Join(dragonLines, "\n")
		}
		dragonStyle := lipgloss.NewStyle().Foreground(t.Primary)
		dragonRendered = dragonStyle.Render(dragon)
	}

	var sections []string
	sections = append(sections, infoSection, "", statusContent)
	if dragonRendered != "" {
		sections = append(sections, "", dragonRendered)
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(mainHeight).
		PaddingTop(1).
		Render(
			lipgloss.JoinVertical(lipgloss.Left, sections...),
		)
}
