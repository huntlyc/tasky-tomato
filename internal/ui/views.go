package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/huntlyc/tasky-tomato/internal/config"
	"github.com/huntlyc/tasky-tomato/internal/models"
)

func (m Model) renderBoard() string {
	top := m.renderTopBar()
	board := lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderColumn(config.ColTodo, "Todo", "todo"),
		m.renderColumn(config.ColDoing, "Doing", "doing"),
		m.renderColumn(config.ColDone, "Done", "done"),
	)

	pomo := m.renderPomo()
	footer := m.renderFooter()

	content := lipgloss.JoinVertical(lipgloss.Left, top, board, pomo, footer)
	return ScreenStyle.Width(m.contentWidth()).Render(content)
}

func (m Model) renderTopBar() string {
	title := TitleStyle.Render("tasky tomato")
	msg := ""
	if m.message != "" {
		msg = MessageStyle.Render(m.message)
	}
	headerRow := lipgloss.JoinHorizontal(lipgloss.Left, title, lipgloss.NewStyle().Width(m.contentWidth()-20).Align(lipgloss.Right).Render(msg))
	return headerRow
}

func (m Model) renderColumn(col int, title string, status string) string {
	tasks := m.tasksByStatus(status)
	selected := m.selected[col]
	isCurrent := m.currentCol == col

	count := fmt.Sprintf("%d", len(tasks))
	head := ColumnHeaderStyle.Render(fmt.Sprintf("%s  %s", title, CountPill(status, count)))
	if isCurrent {
		head = CurrentHeaderStyle.Render(fmt.Sprintf("%s  %s", title, CountPill(status, count)))
	}

	items := make([]string, 0, len(tasks)+1)
	if len(tasks) == 0 {
		items = append(items, EmptyTaskStyle.Render("no tasks"))
	}
	for i, t := range tasks {
		desc := truncateLines(t.Description, 50)
		card := TaskCard(status, t.Title, desc)
		if isCurrent && i == selected {
			card = TaskCardSelected(status, t.Title, desc)
		}
		items = append(items, card)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	return ColumnFrameStyle.Width(m.columnWidth()).Render(lipgloss.JoinVertical(lipgloss.Left, head, content))
}

func (m Model) renderPomo() string {
	if !m.pomo.Active {
		return ""
	}

	phaseName := strings.ReplaceAll(string(m.pomo.Phase), "_", " ")
	var total time.Duration
	var style lipgloss.Style
	switch m.pomo.Phase {
	case models.PhaseWork:
		total = time.Duration(m.settings.WorkMin) * time.Minute
		style = PomoWorkStyle
	case models.PhaseShortBreak:
		total = time.Duration(m.settings.ShortBreakMin) * time.Minute
		style = PomoBreakStyle
	case models.PhaseLongBreak:
		total = time.Duration(m.settings.LongBreakMin) * time.Minute
		style = PomoLongBreakStyle
	}

	progress := 0
	if total > 0 {
		progress = int(float64(m.pomo.Remaining) / float64(total) * 100)
	}
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	barWidth := 40
	filled := (barWidth * progress) / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	label := fmt.Sprintf("%s %s %s", phaseName, bar, formatDuration(m.pomo.Remaining))
	if m.pomo.Paused {
		label += " (paused)"
		style = PomoPausedStyle
	}
	label += fmt.Sprintf(" • cycle %d/%d", m.pomo.CycleCount+1, m.settings.SessionsBeforeLong)

	return style.Render(label)
}

func (m Model) renderFooter() string {
	if !m.showHelp {
		return HelpCollapsedStyle.Render("? for help")
	}

	pomoLabel := "p: start pomo"
	if m.pomo.Active {
		if m.pomo.Paused {
			pomoLabel = "p: resume"
		} else {
			pomoLabel = "p: pause"
		}
	}

	items := []string{
		"h/j/k/l: move", "shift+h/j/k/l: move task",
		"n: new", "e: edit", pomoLabel, "s: settings", "q: quit",
	}
	return HelpStyle.Render(strings.Join(items, "  •  "))
}

func (m Model) renderTaskForm() string {
	title := m.taskForm.Title.View()
	desc := m.taskForm.Desc.View()
	if m.taskForm.Focus == 0 {
		title = FocusedFieldStyle.Render(title)
	} else if m.taskForm.Focus == 1 {
		desc = FocusedFieldStyle.Render(desc)
	}

	var saveBtn, discardBtn string
	if m.taskForm.ButtonFocus == 0 {
		saveBtn = ButtonPrimaryFocusedStyle.Render("Save")
		discardBtn = ButtonSecondaryStyle.Render("Discard")
	} else {
		saveBtn = ButtonPrimaryStyle.Render("Save")
		discardBtn = ButtonSecondaryFocusedStyle.Render("Discard")
	}
	buttons := lipgloss.JoinHorizontal(lipgloss.Left, saveBtn, discardBtn)

	panel := FormPanelStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		HeaderStyle.Render(func() string {
			if m.taskForm.Edit {
				return "Edit task"
			}
			return "New task"
		}()),
		FormHintStyle.Render("Tab to switch fields • Enter in title → description • Ctrl+S save • Esc discard"),
		LabelStyle.Render("Title"),
		title,
		LabelStyle.Render("Description"),
		desc,
		buttons,
	))
	return centerOverlay(panel, m.contentWidth())
}

func (m Model) renderSettingsForm() string {
	fields := []string{
		m.settingsForm.WorkMin.View(),
		m.settingsForm.ShortBreak.View(),
		m.settingsForm.LongBreak.View(),
		m.settingsForm.Sessions.View(),
	}
	for i := range fields {
		if i == m.settingsForm.Focus {
			fields[i] = FocusedFieldStyle.Render(fields[i])
		}
	}
	panel := FormPanelStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		HeaderStyle.Render("Pomodoro settings"),
		FormHintStyle.Render("Tab to switch fields • Enter to save • Esc to cancel"),
		LabelStyle.Render("Work minutes"), fields[0],
		LabelStyle.Render("Short break minutes"), fields[1],
		LabelStyle.Render("Long break minutes"), fields[2],
		LabelStyle.Render("Sessions before long break"), fields[3],
	))
	return centerOverlay(panel, m.contentWidth())
}

func (m Model) renderDeleteConfirm() string {
	panel := FormPanelStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		HeaderStyle.Render("Delete task"),
		FormHintStyle.Render(m.message),
		FormHintStyle.Render("Press y to delete or n to cancel"),
	))
	return centerOverlay(panel, m.contentWidth())
}
