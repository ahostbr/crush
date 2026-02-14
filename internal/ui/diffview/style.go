package diffview

import (
	"charm.land/lipgloss/v2"
)

// LineStyle defines the styles for a given line type in the diff view.
type LineStyle struct {
	LineNumber lipgloss.Style
	Symbol     lipgloss.Style
	Code       lipgloss.Style
}

// Style defines the overall style for the diff view, including styles for
// different line types such as divider, missing, equal, insert, and delete
// lines.
type Style struct {
	DividerLine LineStyle
	MissingLine LineStyle
	EqualLine   LineStyle
	InsertLine  LineStyle
	DeleteLine  LineStyle
}

// DefaultLightStyle provides a default light theme style for the diff view.
func DefaultLightStyle() Style {
	return Style{
		DividerLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#555555")).
				Background(lipgloss.Color("#333333")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#999999")).
				Background(lipgloss.Color("#222222")),
		},
		MissingLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Background(lipgloss.Color("#333333")),
			Code: lipgloss.NewStyle().
				Background(lipgloss.Color("#333333")),
		},
		EqualLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#333333")).
				Background(lipgloss.Color("#333333")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1A1A1A")).
				Background(lipgloss.Color("#E0E0E0")),
		},
		InsertLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4A8C4A")).
				Background(lipgloss.Color("#2b322a")),
			Symbol: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4A8C4A")).
				Background(lipgloss.Color("#323931")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1A1A1A")).
				Background(lipgloss.Color("#323931")),
		},
		DeleteLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CC3333")).
				Background(lipgloss.Color("#312929")),
			Symbol: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CC3333")).
				Background(lipgloss.Color("#383030")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1A1A1A")).
				Background(lipgloss.Color("#383030")),
		},
	}
}

// DefaultDarkStyle provides a default dark theme style for the diff view.
func DefaultDarkStyle() Style {
	return Style{
		DividerLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#C8C8C8")).
				Background(lipgloss.Color("#445577")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#C8C8C8")).
				Background(lipgloss.Color("#222222")),
		},
		MissingLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Background(lipgloss.Color("#333333")),
			Code: lipgloss.NewStyle().
				Background(lipgloss.Color("#333333")),
		},
		EqualLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#808080")).
				Background(lipgloss.Color("#1A1A1A")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E0E0E0")).
				Background(lipgloss.Color("#0D0D0D")),
		},
		InsertLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4A8C4A")).
				Background(lipgloss.Color("#293229")),
			Symbol: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4A8C4A")).
				Background(lipgloss.Color("#303a30")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E0E0E0")).
				Background(lipgloss.Color("#303a30")),
		},
		DeleteLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CC3333")).
				Background(lipgloss.Color("#332929")),
			Symbol: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CC3333")).
				Background(lipgloss.Color("#3a3030")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E0E0E0")).
				Background(lipgloss.Color("#3a3030")),
		},
	}
}
