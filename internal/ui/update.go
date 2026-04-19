package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/huntlyc/tasky-tomato/internal/config"
	"github.com/huntlyc/tasky-tomato/internal/db"
	"github.com/huntlyc/tasky-tomato/internal/models"
)

func (m Model) updateBoard(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, Keys.Help):
		m.showHelp = !m.showHelp
		return m, nil
	case key.Matches(msg, Keys.Left):
		if m.currentCol > 0 {
			m.currentCol--
		}
		return m, nil
	case key.Matches(msg, Keys.Right):
		if m.currentCol < 2 {
			m.currentCol++
		}
		return m, nil
	case key.Matches(msg, Keys.Up):
		if sel := m.selected[m.currentCol]; sel > 0 {
			m.selected[m.currentCol] = sel - 1
		}
		return m, nil
	case key.Matches(msg, Keys.Down):
		if sel := m.selected[m.currentCol]; sel < len(m.tasksInCol(m.currentCol))-1 {
			m.selected[m.currentCol] = sel + 1
		}
		return m, nil
	case key.Matches(msg, Keys.MoveUp):
		if err := m.moveSelectedTask(-1, 0); err != nil {
			m.err = err
		}
		return m, nil
	case key.Matches(msg, Keys.MoveDown):
		if err := m.moveSelectedTask(1, 0); err != nil {
			m.err = err
		}
		return m, nil
	case key.Matches(msg, Keys.MoveLeft):
		if err := m.moveSelectedTask(0, -1); err != nil {
			m.err = err
		}
		return m, nil
	case key.Matches(msg, Keys.MoveRight):
		if err := m.moveSelectedTask(0, 1); err != nil {
			m.err = err
		}
		return m, nil
	case key.Matches(msg, Keys.New):
		m.taskForm = newTaskForm(false, models.Task{})
		m.mode = config.ModeTaskForm
		return m, nil
	case key.Matches(msg, Keys.Edit):
		tasks := m.tasksInCol(m.currentCol)
		if len(tasks) == 0 {
			return m, nil
		}
		sel := m.selected[m.currentCol]
		if sel >= len(tasks) {
			sel = len(tasks) - 1
		}
		m.taskForm = newTaskForm(true, tasks[sel])
		m.mode = config.ModeTaskForm
		return m, nil
	case key.Matches(msg, Keys.Delete):
		tasks := m.tasksInCol(m.currentCol)
		if len(tasks) == 0 {
			return m, nil
		}
		sel := m.selected[m.currentCol]
		if sel >= len(tasks) {
			sel = len(tasks) - 1
		}
		m.taskForm.ID = tasks[sel].ID
		m.message = fmt.Sprintf("Delete %q? (y/n)", tasks[sel].Title)
		m.messageTime = time.Now()
		m.mode = config.ModeConfirmDelete
		return m, nil
	case key.Matches(msg, Keys.Advance):
		if err := m.advanceSelectedTask(); err != nil {
			m.err = err
		}
		return m, nil
	case key.Matches(msg, Keys.Pomo):
		if m.pomo.Active {
			if m.pomo.Paused {
				m.pomo.Paused = false
				return m, TickCmd()
			}
			m.pomo.Paused = true
			return m, nil
		}
		workMin := m.settings.WorkMin
		if workMin <= 0 {
			workMin = 25
		}
		m.pomo.Active = true
		m.pomo.Paused = false
		m.pomo.Phase = models.PhaseWork
		m.pomo.Remaining = time.Duration(workMin) * time.Minute
		m.pomo.CycleCount = 0
		return m, TickCmd()
	case key.Matches(msg, Keys.Settings):
		m.settingsForm = newSettingsForm(m.settings)
		m.mode = config.ModeSettingsForm
		return m, nil
	}
	return m, nil
}

func (m *Model) advanceSelectedTask() error {
	tasks := m.tasksInCol(m.currentCol)
	if len(tasks) == 0 {
		return nil
	}
	sel := m.selected[m.currentCol]
	if sel >= len(tasks) {
		sel = len(tasks) - 1
	}
	t := tasks[sel]
	switch t.Status {
	case "todo":
		t.Status = "doing"
	case "doing":
		t.Status = "done"
	case "done":
		t.Status = "todo"
	}
	if err := m.UpdateTask(t); err != nil {
		return err
	}
	if err := m.loadTasks(); err != nil {
		return err
	}
	if m.selected[m.currentCol] >= len(m.tasksInCol(m.currentCol)) && m.selected[m.currentCol] > 0 {
		m.selected[m.currentCol]--
	}
	return nil
}

func (m *Model) moveSelectedTask(vertical int, horizontal int) error {
	if vertical == 0 && horizontal == 0 {
		return nil
	}
	tasks := m.tasksInCol(m.currentCol)
	if len(tasks) == 0 {
		return nil
	}
	sel := m.selected[m.currentCol]
	if sel >= len(tasks) {
		sel = len(tasks) - 1
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	current := tasks[sel]

	if horizontal != 0 {
		targetCol := m.currentCol + horizontal
		if targetCol < config.ColTodo || targetCol > config.ColDone {
			return nil
		}
		targetStatus := []string{"todo", "doing", "done"}[targetCol]
		err := db.MoveTask(m.db, current.ID, targetStatus, now)
		if err != nil {
			return err
		}
		if err := m.loadTasks(); err != nil {
			return err
		}
		m.currentCol = targetCol
		m.selected[targetCol] = len(m.tasksInCol(targetCol)) - 1
		return nil
	}

	targetIdx := sel + vertical
	if targetIdx < 0 || targetIdx >= len(tasks) {
		return nil
	}
	neighbor := tasks[targetIdx]
	if err := db.ReorderTask(m.db, neighbor.ID, now); err != nil {
		return err
	}
	if err := db.ReorderTask(m.db, current.ID, neighbor.CreatedAt); err != nil {
		return err
	}
	if err := m.loadTasks(); err != nil {
		return err
	}
	m.selected[m.currentCol] = targetIdx
	return nil
}

func (m Model) updateTaskForm(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.Discard):
		m.mode = config.ModeBoard
		return m, nil
	case key.Matches(msg, Keys.Save):
		m.messageTime = time.Now()
		return m.saveTask()
	case key.Matches(msg, Keys.Tab):
		if m.taskForm.Focus < 2 {
			m.taskForm.Focus = m.taskForm.Focus + 1
		} else {
			m.taskForm.ButtonFocus = (m.taskForm.ButtonFocus + 1) % 2
		}
		m.syncTaskFormFocus()
		return m, nil
	case key.Matches(msg, Keys.ShiftTab):
		if m.taskForm.Focus > 0 {
			m.taskForm.Focus = m.taskForm.Focus - 1
		} else {
			m.taskForm.ButtonFocus = (m.taskForm.ButtonFocus + 1) % 2
		}
		m.syncTaskFormFocus()
		return m, nil
	case msg.String() == "enter" && m.taskForm.Focus >= 2:
		if m.taskForm.ButtonFocus == 0 {
			return m.saveTask()
		}
		m.mode = config.ModeBoard
		return m, nil
	}

	if m.taskForm.Focus == 0 {
		if msg.String() == "enter" {
			m.taskForm.Focus = 1
			m.syncTaskFormFocus()
			return m, nil
		}
		var cmd tea.Cmd
		m.taskForm.Title, cmd = m.taskForm.Title.Update(msg)
		return m, cmd
	}
	var cmd tea.Cmd
	m.taskForm.Desc, cmd = m.taskForm.Desc.Update(msg)
	return m, cmd
}

func (m Model) updateSettingsForm(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.Discard):
		m.mode = config.ModeBoard
		return m, nil
	case key.Matches(msg, Keys.Save):
		work, ok1 := parsePositiveInt(m.settingsForm.WorkMin.Value())
		shortB, ok2 := parsePositiveInt(m.settingsForm.ShortBreak.Value())
		longB, ok3 := parsePositiveInt(m.settingsForm.LongBreak.Value())
		sessions, ok4 := parsePositiveInt(m.settingsForm.Sessions.Value())
		if !ok1 || !ok2 || !ok3 || !ok4 || sessions < 1 {
			m.message = "enter valid positive integers"
			m.messageTime = time.Now()
			return m, nil
		}
		m.settings = models.Settings{WorkMin: work, ShortBreakMin: shortB, LongBreakMin: longB, SessionsBeforeLong: sessions}
		if err := m.SaveSettings(); err != nil {
			m.err = err
			return m, nil
		}
		m.mode = config.ModeBoard
		m.message = "settings saved"
		m.messageTime = time.Now()
		return m, nil
	case key.Matches(msg, Keys.Tab):
		m.settingsForm.Focus = (m.settingsForm.Focus + 1) % 4
		m.syncSettingsFormFocus()
		return m, nil
	case key.Matches(msg, Keys.ShiftTab):
		m.settingsForm.Focus = (m.settingsForm.Focus + 3) % 4
		m.syncSettingsFormFocus()
		return m, nil
	}

	var cmd tea.Cmd
	switch m.settingsForm.Focus {
	case 0:
		m.settingsForm.WorkMin, cmd = m.settingsForm.WorkMin.Update(msg)
	case 1:
		m.settingsForm.ShortBreak, cmd = m.settingsForm.ShortBreak.Update(msg)
	case 2:
		m.settingsForm.LongBreak, cmd = m.settingsForm.LongBreak.Update(msg)
	case 3:
		m.settingsForm.Sessions, cmd = m.settingsForm.Sessions.Update(msg)
	}
	return m, cmd
}

func (m Model) updateDeleteConfirm(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "y", "enter":
		if err := m.DeleteTask(m.taskForm.ID); err != nil {
			m.err = err
			return m, nil
		}
		if err := m.loadTasks(); err != nil {
			m.err = err
			return m, nil
		}
		m.mode = config.ModeBoard
		m.message = "deleted"
		m.messageTime = time.Now()
		return m, nil
	case "n", "esc":
		m.mode = config.ModeBoard
		return m, nil
	}
	return m, nil
}

func (m Model) updateTick() (Model, tea.Cmd) {
	if m.message != "" && time.Since(m.messageTime) > 3*time.Second {
		m.message = ""
	}

	if !m.pomo.Active || m.pomo.Paused {
		return m, TickCmd()
	}
	if m.pomo.Remaining > time.Second {
		m.pomo.Remaining -= time.Second
		return m, TickCmd()
	}
	switch m.pomo.Phase {
	case models.PhaseWork:
		m.pomo.CycleCount++
		if m.pomo.CycleCount >= m.settings.SessionsBeforeLong {
			m.pomo.Phase = models.PhaseLongBreak
			m.pomo.Remaining = time.Duration(m.settings.LongBreakMin) * time.Minute
			m.pomo.CycleCount = 0
		} else {
			m.pomo.Phase = models.PhaseShortBreak
			m.pomo.Remaining = time.Duration(m.settings.ShortBreakMin) * time.Minute
		}
	case models.PhaseShortBreak, models.PhaseLongBreak:
		m.pomo.Phase = models.PhaseWork
		m.pomo.Remaining = time.Duration(m.settings.WorkMin) * time.Minute
	}
	return m, TickCmd()
}
