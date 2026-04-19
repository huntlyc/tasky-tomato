package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	MoveUp    key.Binding
	MoveDown  key.Binding
	MoveLeft  key.Binding
	MoveRight key.Binding
	New       key.Binding
	Edit      key.Binding
	Delete    key.Binding
	Advance   key.Binding
	Pomo      key.Binding
	Settings  key.Binding
	Save      key.Binding
	Discard   key.Binding
	Cancel    key.Binding
	Tab       key.Binding
	ShiftTab  key.Binding
	Quit      key.Binding
	Help      key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/up", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/down", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h/left", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l/right", "move right"),
	),
	MoveUp:    key.NewBinding(key.WithKeys("shift+k", "K"), key.WithHelp("⇧k", "move task up")),
	MoveDown:  key.NewBinding(key.WithKeys("shift+j", "J"), key.WithHelp("⇧j", "move task down")),
	MoveLeft:  key.NewBinding(key.WithKeys("shift+h", "H"), key.WithHelp("⇧h", "move task left")),
	MoveRight: key.NewBinding(key.WithKeys("shift+l", "L"), key.WithHelp("⇧l", "move task right")),
	New:       key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new task")),
	Edit:      key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit task")),
	Delete:    key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete task")),
	Advance:   key.NewBinding(key.WithKeys("space", "enter"), key.WithHelp("space", "advance task")),
	Pomo:      key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "start pomo")),
	Settings:  key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "settings")),
	Save:      key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
	Discard:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "discard")),
	Cancel:    key.NewBinding(key.WithKeys("ctrl+["), key.WithHelp("esc", "cancel")),
	Tab:       key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
	ShiftTab:  key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev field")),
	Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
}
