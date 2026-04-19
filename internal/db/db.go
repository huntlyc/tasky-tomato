package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/huntlyc/tasky-tomato/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func Init(db *sql.DB) error {
	stmts := []string{
		`PRAGMA foreign_keys = ON;`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL CHECK(status IN ('todo','doing','done')) DEFAULT 'todo',
			created_at TEXT NOT NULL DEFAULT (datetime('now'))
		);`,
		`CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			work_min INTEGER NOT NULL DEFAULT 25,
			short_break_min INTEGER NOT NULL DEFAULT 5,
			long_break_min INTEGER NOT NULL DEFAULT 15,
			sessions_before_long INTEGER NOT NULL DEFAULT 4
		);`,
		`INSERT OR IGNORE INTO settings (id, work_min, short_break_min, long_break_min, sessions_before_long)
		 VALUES (1, 25, 5, 15, 4);`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func ListTasks(db *sql.DB) ([]models.Task, error) {
	rows, err := db.Query(`SELECT id, title, description, status, created_at FROM tasks ORDER BY created_at ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func GetSettings(db *sql.DB) (models.Settings, error) {
	var s models.Settings
	row := db.QueryRow(`SELECT work_min, short_break_min, long_break_min, sessions_before_long FROM settings WHERE id = 1`)
	err := row.Scan(&s.WorkMin, &s.ShortBreakMin, &s.LongBreakMin, &s.SessionsBeforeLong)
	return s, err
}

func SaveSettings(db *sql.DB, s models.Settings) error {
	_, err := db.Exec(`UPDATE settings SET work_min=?, short_break_min=?, long_break_min=?, sessions_before_long=? WHERE id = 1`,
		s.WorkMin, s.ShortBreakMin, s.LongBreakMin, s.SessionsBeforeLong)
	return err
}

func AddTask(db *sql.DB, title, desc string) error {
	_, err := db.Exec(`INSERT INTO tasks (title, description, status) VALUES (?, ?, 'todo')`, title, desc)
	return err
}

func UpdateTask(db *sql.DB, t models.Task) error {
	_, err := db.Exec(`UPDATE tasks SET title=?, description=?, status=? WHERE id=?`, t.Title, t.Description, t.Status, t.ID)
	return err
}

func DeleteTask(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM tasks WHERE id=?`, id)
	return err
}

func MoveTask(db *sql.DB, id int, status, createdAt string) error {
	_, err := db.Exec(`UPDATE tasks SET status = ?, created_at = ? WHERE id = ?`, status, createdAt, id)
	return err
}

func ReorderTask(db *sql.DB, id int, createdAt string) error {
	_, err := db.Exec(`UPDATE tasks SET created_at = ? WHERE id = ?`, createdAt, id)
	return err
}

var ErrNotFound = fmt.Errorf("not found")
