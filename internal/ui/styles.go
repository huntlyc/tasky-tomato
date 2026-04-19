package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	baseColor     = lipgloss.Color("#1e1e2e")
	surface2Color = lipgloss.Color("#585b70")
	textColor     = lipgloss.Color("#cdd6f4")
	subtextColor  = lipgloss.Color("#a6adc8")
	blueColor     = lipgloss.Color("#89b4fa")
	lavenderColor = lipgloss.Color("#b4befe")
	greenColor    = lipgloss.Color("#a6e3a1")
	yellowColor   = lipgloss.Color("#f9e2af")
	peachColor    = lipgloss.Color("#fab387")
	redColor      = lipgloss.Color("#f38ba8")
	mauveColor    = lipgloss.Color("#cba6f7")
	tealColor     = lipgloss.Color("#94e2d5")
	roseColor     = lipgloss.Color("#f5e0dc")

	ScreenStyle                 = lipgloss.NewStyle().Foreground(textColor).Margin(0, 2)
	AppStyle                    = lipgloss.NewStyle().Foreground(textColor).Padding(1, 2)
	ColumnFrameStyle            = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(surface2Color)
	CurrentHeaderStyle          = lipgloss.NewStyle().Bold(true).Foreground(lavenderColor)
	ColumnHeaderStyle           = lipgloss.NewStyle().Bold(true).Foreground(subtextColor)
	HeaderStyle                 = lipgloss.NewStyle().Bold(true).Foreground(roseColor)
	TitleStyle                  = lipgloss.NewStyle().Bold(true).Foreground(roseColor).MarginBottom(2).Inline(true)
	SubtitleStyle               = lipgloss.NewStyle().Foreground(subtextColor).Italic(true).PaddingLeft(1)
	MessageStyle                = lipgloss.NewStyle().Foreground(yellowColor).PaddingTop(1)
	HelpStyle                   = lipgloss.NewStyle().Foreground(subtextColor).PaddingTop(1)
	HelpCollapsedStyle          = lipgloss.NewStyle().Foreground(surface2Color).Italic(true).PaddingTop(1)
	PomoIdleStyle               = lipgloss.NewStyle().Foreground(subtextColor).Bold(true).Padding(0, 1).MarginTop(1)
	PomoWorkStyle               = lipgloss.NewStyle().Foreground(greenColor).Bold(true).Padding(0, 1).MarginTop(1)
	PomoBreakStyle              = lipgloss.NewStyle().Foreground(tealColor).Bold(true).Padding(0, 1).MarginTop(1)
	PomoLongBreakStyle          = lipgloss.NewStyle().Foreground(peachColor).Bold(true).Padding(0, 1).MarginTop(1)
	TaskTitleStyle              = lipgloss.NewStyle().Bold(true).Foreground(textColor).Margin(0).Padding(0)
	TaskDescStyle               = lipgloss.NewStyle().Foreground(subtextColor).Margin(0).Padding(0)
	EmptyTaskStyle              = lipgloss.NewStyle().Foreground(surface2Color).Italic(true).Padding(1, 0)
	LabelStyle                  = lipgloss.NewStyle().Bold(true).Foreground(subtextColor).PaddingTop(1)
	FormHintStyle               = lipgloss.NewStyle().Foreground(surface2Color).PaddingBottom(1)
	FocusedFieldStyle           = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(mauveColor).Padding(0, 1)
	FormPanelStyle              = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(surface2Color).Padding(1, 2)
	ButtonPrimaryStyle          = lipgloss.NewStyle().Foreground(greenColor).Border(lipgloss.RoundedBorder()).BorderForeground(greenColor).Bold(true).Padding(0, 2).MarginRight(1)
	ButtonSecondaryStyle        = lipgloss.NewStyle().Foreground(textColor).Border(lipgloss.RoundedBorder()).BorderForeground(surface2Color).Padding(0, 2)
	ButtonPrimaryFocusedStyle   = lipgloss.NewStyle().Foreground(greenColor).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("white")).Bold(true).Padding(0, 2).MarginRight(1)
	ButtonSecondaryFocusedStyle = lipgloss.NewStyle().Foreground(textColor).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("white")).Padding(0, 2)
)

func ColorByStatus(status string) lipgloss.Color {
	switch status {
	case "todo":
		return mauveColor
	case "doing":
		return peachColor
	case "done":
		return greenColor
	default:
		return blueColor
	}
}

func TaskCard(status, title, desc string) string {
	borderc := ColorByStatus(status)
	lines := []string{title}
	if desc != "" {
		lines = append(lines, desc)
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderc).
		Padding(0, 1).
		Foreground(textColor).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func TaskCardSelected(_, title, desc string) string {
	lines := []string{title}
	if desc != "" {
		lines = append(lines, desc)
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("white")).
		Padding(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func CountPill(status, count string) string {
	return lipgloss.NewStyle().
		Foreground(ColorByStatus(status)).
		Padding(0, 1).
		Render(count)
}

func StatusPill(status, count string) string {
	label := map[string]string{"todo": "todo", "doing": "doing", "done": "done"}[status]
	return lipgloss.NewStyle().
		Foreground(ColorByStatus(status)).
		Padding(0, 1).
		MarginRight(1).
		Render(fmt.Sprintf("%s %s", label, count))
}
