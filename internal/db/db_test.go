package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/huntlyc/tasky-tomato/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	f, err := os.CreateTemp("", "test-*.db")
	require.NoError(t, err)
	f.Close()

	db, err := Open(f.Name())
	require.NoError(t, err)

	err = Init(db)
	require.NoError(t, err)

	return db
}

func TestInit(t *testing.T) {
	f, err := os.CreateTemp("", "test-*.db")
	assert.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())

	db, err := Open(f.Name())
	assert.NoError(t, err)
	defer db.Close()

	err = Init(db)
	assert.NoError(t, err)
}

func TestListTasks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tasks, err := ListTasks(db)
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestAddTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := AddTask(db, "Test Task", "Test description")
	assert.NoError(t, err)

	tasks, err := ListTasks(db)
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Test Task", tasks[0].Title)
	assert.Equal(t, "Test description", tasks[0].Description)
	assert.Equal(t, "todo", tasks[0].Status)
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := AddTask(db, "Original", "")
	assert.NoError(t, err)

	tasks, err := ListTasks(db)
	task := tasks[0]
	task.Title = "Updated"
	task.Status = "doing"

	err = UpdateTask(db, task)
	assert.NoError(t, err)

	tasks, err = ListTasks(db)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", tasks[0].Title)
	assert.Equal(t, "doing", tasks[0].Status)
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := AddTask(db, "To Delete", "")
	assert.NoError(t, err)

	tasks, err := ListTasks(db)
	id := tasks[0].ID

	err = DeleteTask(db, id)
	assert.NoError(t, err)

	tasks, err = ListTasks(db)
	assert.NoError(t, err)
	assert.Len(t, tasks, 0)
}

func TestGetSettings(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	s, err := GetSettings(db)
	assert.NoError(t, err)
	assert.Equal(t, 25, s.WorkMin)
	assert.Equal(t, 5, s.ShortBreakMin)
	assert.Equal(t, 15, s.LongBreakMin)
	assert.Equal(t, 4, s.SessionsBeforeLong)
}

func TestSaveSettings(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	s := models.Settings{
		WorkMin:            30,
		ShortBreakMin:      10,
		LongBreakMin:       20,
		SessionsBeforeLong: 3,
	}
	err := SaveSettings(db, s)
	assert.NoError(t, err)

	loaded, err := GetSettings(db)
	assert.NoError(t, err)
	assert.Equal(t, 30, loaded.WorkMin)
	assert.Equal(t, 10, loaded.ShortBreakMin)
	assert.Equal(t, 20, loaded.LongBreakMin)
	assert.Equal(t, 3, loaded.SessionsBeforeLong)
}
