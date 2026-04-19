package models

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	CreatedAt   string
}

type Settings struct {
	WorkMin            int
	ShortBreakMin      int
	LongBreakMin       int
	SessionsBeforeLong int
}

type PomoPhase string

const (
	PhaseWork       PomoPhase = "work"
	PhaseShortBreak PomoPhase = "short_break"
	PhaseLongBreak  PomoPhase = "long_break"
)

type PomoState struct {
	Active     bool
	Phase      PomoPhase
	Remaining  time.Duration
	CycleCount int
	Paused     bool
}
