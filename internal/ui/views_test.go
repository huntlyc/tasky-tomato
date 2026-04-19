package ui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestColorByStatus(t *testing.T) {
	tests := []struct {
		status string
	}{
		{"todo"},
		{"doing"},
		{"done"},
		{"unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := ColorByStatus(tt.status)
			assert.NotNil(t, got)
		})
	}
}

func TestTaskCard(t *testing.T) {
	tests := []struct {
		name  string
		title string
		desc  string
	}{
		{"with description", "Test", "Desc"},
		{"without description", "Test", ""},
		{"multi-line", "Test", "Line1\nLine2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TaskCard("todo", tt.title, tt.desc)
			assert.NotEmpty(t, got)
		})
	}
}

func TestTaskCardSelected(t *testing.T) {
	got := TaskCardSelected("todo", "Test", "Desc")
	assert.NotEmpty(t, got)
}

func TestCountPill(t *testing.T) {
	tests := []struct {
		status string
		count  string
	}{
		{"todo", "3"},
		{"doing", "1"},
		{"done", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := CountPill(tt.status, tt.count)
			assert.NotEmpty(t, got)
			assert.Contains(t, got, tt.count)
		})
	}
}

func TestStatusPill(t *testing.T) {
	tests := []struct {
		status string
		count  string
	}{
		{"todo", "3"},
		{"doing", "1"},
		{"done", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := StatusPill(tt.status, tt.count)
			assert.NotEmpty(t, got)
			assert.Contains(t, got, tt.status)
			assert.Contains(t, got, tt.count)
		})
	}
}

func TestTruncateLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		max   int
		want  string
	}{
		{"short", "short", 50, "short"},
		{"truncate", "this is a very long string", 15, "this is a ve..."},
		{"multi-line", "line one\nline two", 20, "line one\nline two"},
		{"empty", "", 50, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateLines(tt.input, tt.max)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParsePositiveInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantVal int
		wantOk  bool
	}{
		{"valid", "25", 25, true},
		{"zero", "0", 0, false},
		{"invalid", "abc", 0, false},
		{"negative", "-5", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parsePositiveInt(tt.input)
			assert.Equal(t, tt.wantOk, ok)
			if ok {
				assert.Equal(t, tt.wantVal, got)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name  string
		input time.Duration
		want  string
	}{
		{"zero", 0, "00:00"},
		{"minutes", 5 * time.Minute, "05:00"},
		{"seconds", 30 * time.Second, "00:30"},
		{"mixed", 5*time.Minute + 30*time.Second, "05:30"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
