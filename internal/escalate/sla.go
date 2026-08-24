package escalate

import (
	"time"

	"telemetryguard/internal/store"
)

// LevelMinutes is the SLA window added per escalation level.
const LevelMinutes = 30

// DueAt computes the escalation deadline from the alarm's original open time.
func DueAt(a *store.Alarm, level int) time.Time {
	return a.OriginalOpenAt.Add(time.Duration(level*LevelMinutes) * time.Minute)
}
