# Tasky Tomato

A test of vibe coding. It's always a todo app, right?

This is a simple TUI Kanban task manager with built-in Pomodoro timer, using Charm's tooling.

This was done in opencode with 'big pickle'... .... ..

## Features

- Kanban board: Todo, Doing, Done
- CRUD operations on tasks
- Vim-like keyboard navigation
- Customizable Pomodoro settings (work duration, short break, long break)
- Pomodoro workflow: 4 work sessions + short break, then long break
- Data stored in SQLite

## Usage

### Navigation

- `h` / `left`: Move to previous column
- `l` / `right`: Move to next column
- `j` / `down`: Select next task
- `k` / `up`: Select previous task
- `space`: Move selected task to next status (todo -> doing -> done -> todo)

### Task Management

- `a`: Add new task
- `e`: Edit selected task
- `d`: Delete selected task (confirm with y/n)

### Pomodoro

- `s`: Start/stop Pomodoro (only when in "Doing" column and tasks present)
- Timer runs automatically, switching phases

### Settings

- `c`: Configure Pomodoro settings (work minutes, short break minutes)
- Long break is automatically set to 2x short break

### General

- `q` / `Ctrl+C`: Quit

## Building

```bash
go build
```

## Running

```bash
./tasky-tomato
```

Data is stored in `tasky.db`.
