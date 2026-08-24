package escalate

import (
	"time"

	"telemetryguard/internal/store"
)

// LevelMinutes is the SLA window added per escalation level.
const LevelMinutes = 30

// DueAt computes the escalation deadline from the alarm's original open time.
// Anchoring on the original open time keeps the SLA countdown stable across
// recover-and-retrigger cycles, so a flapping alarm cannot keep pushing its
// escalation out indefinitely.
func DueAt(a *store.Alarm, level int) time.Time {
	return a.OriginalOpenAt.Add(time.Duration(level*LevelMinutes) * time.Minute)
}
