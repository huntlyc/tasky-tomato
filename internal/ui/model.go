package ui

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/huntlyc/tasky-tomato/internal/config"
	"github.com/huntlyc/tasky-tomato/internal/db"
	"github.com/huntlyc/tasky-tomato/internal/models"
)

type TaskFormState struct {
	Title       textinput.Model
	Desc        textarea.Model
	Focus       int
	Edit        bool
	ID          int
	ButtonFocus int // 0 = save, 1 = discard
}

type SettingsFormState struct {
	WorkMin    textinput.Model
	ShortBreak textinput.Model
	LongBreak  textinput.Model
	Sessions   textinput.Model
	Focus      int
}

type Model struct {
	db           *sql.DB
	tasks        []models.Task
	settings     models.Settings
	pomo         models.PomoState
	currentCol   int
	selected     [3]int
	mode         int
	message      string
	messageTime  time.Time
	width        int
	height       int
	showHelp     bool
	taskForm     TaskFormState
	settingsForm SettingsFormState
	err          error
}

type tickMsg struct{}

func TickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg { return tickMsg{} })
}

func NewModel(db *sql.DB) Model {
	m := Model{db: db, mode: config.ModeBoard, showHelp: true}
	if err := m.loadSettings(); err != nil {
		m.err = err
	}
	if err := m.loadTasks(); err != nil {
		m.err = err
	}
	m.settingsForm = newSettingsForm(m.settings)
	m.taskForm = newTaskForm(false, models.Task{})
	return m
}

func (m *Model) loadTasks() error {
	tasks, err := db.ListTasks(m.db)
	if err != nil {
		return err
	}
	m.tasks = tasks
	return nil
}

func (m *Model) loadSettings() error {
	settings, err := db.GetSettings(m.db)
	if err != nil {
		return err
	}
	m.settings = settings
	return nil
}

func (m *Model) SaveSettings() error {
	return db.SaveSettings(m.db, m.settings)
}

func (m *Model) AddTask(title, desc string) error {
	return db.AddTask(m.db, title, desc)
}

func (m *Model) UpdateTask(t models.Task) error {
	return db.UpdateTask(m.db, t)
}

func (m *Model) DeleteTask(id int) error {
	return db.DeleteTask(m.db, id)
}

func (m Model) saveTask() (Model, tea.Cmd) {
	if strings.TrimSpace(m.taskForm.Title.Value()) == "" {
		m.message = "title is required"
		m.messageTime = time.Now()
		return m, nil
	}
	if m.taskForm.Edit {
		t := models.Task{ID: m.taskForm.ID, Title: strings.TrimSpace(m.taskForm.Title.Value()), Description: strings.TrimSpace(m.taskForm.Desc.Value())}
		for _, task := range m.tasks {
			if task.ID == t.ID {
				t.Status = task.Status
				break
			}
		}
		if err := m.UpdateTask(t); err != nil {
			m.err = err
			return m, nil
		}
	} else {
		if err := m.AddTask(strings.TrimSpace(m.taskForm.Title.Value()), strings.TrimSpace(m.taskForm.Desc.Value())); err != nil {
			m.err = err
			return m, nil
		}
	}
	if err := m.loadTasks(); err != nil {
		m.err = err
		return m, nil
	}
	m.mode = config.ModeBoard
	m.message = "saved"
	m.messageTime = time.Now()
	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		return m.updateTick()
	case tea.KeyMsg:
		switch m.mode {
		case config.ModeBoard:
			return m.updateBoard(msg)
		case config.ModeTaskForm:
			return m.updateTaskForm(msg)
		case config.ModeSettingsForm:
			return m.updateSettingsForm(msg)
		case config.ModeConfirmDelete:
			return m.updateDeleteConfirm(msg)
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return errorView(m.err)
	}

	switch m.mode {
	case config.ModeTaskForm:
		return m.renderTaskForm()
	case config.ModeSettingsForm:
		return m.renderSettingsForm()
	case config.ModeConfirmDelete:
		return m.renderDeleteConfirm()
	default:
		return m.renderBoard()
	}
}

func errorView(err error) string {
	return AppStyle.Render(lipgloss.NewStyle().Foreground(redColor).Render(fmt.Sprintf("Error: %v", err)))
}

func truncateLines(s string, maxChars int) string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, 2)
	for i, line := range lines {
		if i >= 2 {
			break
		}
		if len(line) > maxChars {
			line = line[:maxChars-3] + "..."
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func parsePositiveInt(s string) (int, bool) {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}

func (m Model) contentWidth() int {
	if m.width <= 0 {
		return 114
	}
	return m.width - 4
}

func (m Model) columnWidth() int {
	w := m.contentWidth() - 4 // account for margin
	col := w / 3
	if col < 20 {
		col = 20
	}
	return col - 2
}

func centerOverlay(view string, width int) string {
	if width <= 0 {
		width = 80
	}
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(view)
}

func initTextInput(placeholder string, value string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(value)
	ti.CharLimit = 120
	ti.Width = 40
	return ti
}

func newTaskForm(edit bool, t models.Task) TaskFormState {
	title := initTextInput("Task title", "")
	desc := textarea.New()
	desc.Placeholder = "Description"
	desc.SetHeight(4)
	desc.SetWidth(60)
	desc.ShowLineNumbers = false
	desc.FocusedStyle.Base = lipgloss.NewStyle().BorderForeground(lipgloss.Color("205"))
	desc.BlurredStyle.Base = lipgloss.NewStyle().BorderForeground(lipgloss.Color("240"))
	if edit {
		title.SetValue(t.Title)
		desc.SetValue(t.Description)
	}
	title.Focus()
	desc.Blur()
	return TaskFormState{Title: title, Desc: desc, Focus: 0, Edit: edit, ID: t.ID, ButtonFocus: 0}
}

func newSettingsForm(s models.Settings) SettingsFormState {
	work := initTextInput("work min", strconv.Itoa(s.WorkMin))
	short := initTextInput("short break min", strconv.Itoa(s.ShortBreakMin))
	long := initTextInput("long break min", strconv.Itoa(s.LongBreakMin))
	sessions := initTextInput("sessions before long", strconv.Itoa(s.SessionsBeforeLong))
	work.Focus()
	short.Blur()
	long.Blur()
	sessions.Blur()
	return SettingsFormState{WorkMin: work, ShortBreak: short, LongBreak: long, Sessions: sessions, Focus: 0}
}

func (m *Model) syncTaskFormFocus() {
	switch m.taskForm.Focus {
	case 0:
		m.taskForm.Title.Focus()
		m.taskForm.Desc.Blur()
	case 1:
		m.taskForm.Title.Blur()
		m.taskForm.Desc.Focus()
	}
}

func (m *Model) syncSettingsFormFocus() {
	m.settingsForm.WorkMin.Blur()
	m.settingsForm.ShortBreak.Blur()
	m.settingsForm.LongBreak.Blur()
	m.settingsForm.Sessions.Blur()
	switch m.settingsForm.Focus {
	case 0:
		m.settingsForm.WorkMin.Focus()
	case 1:
		m.settingsForm.ShortBreak.Focus()
	case 2:
		m.settingsForm.LongBreak.Focus()
	case 3:
		m.settingsForm.Sessions.Focus()
	}
}

func (m *Model) tasksByStatus(status string) []models.Task {
	out := make([]models.Task, 0)
	for _, t := range m.tasks {
		if t.Status == status {
			out = append(out, t)
		}
	}
	return out
}

func (m *Model) tasksInCol(col int) []models.Task {
	switch col {
	case config.ColTodo:
		return m.tasksByStatus("todo")
	case config.ColDoing:
		return m.tasksByStatus("doing")
	case config.ColDone:
		return m.tasksByStatus("done")
	default:
		return nil
	}
}
