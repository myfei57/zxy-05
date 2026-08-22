package alarm

import (
	"errors"
	"time"

	"telemetryguard/internal/store"
)

// ErrAlarmNotActive is returned when resolving an alarm that is not open or acknowledged.
var ErrAlarmNotActive = errors.New("alarm is not active")

// Resolve marks an active alarm as resolved after the point recovers.
func Resolve(state *store.State, alarmID string, at time.Time) error {
	a, ok := state.Alarm(alarmID)
	if !ok {
		return errors.New("alarm does not exist")
	}
	if a.Status != store.AlarmOpen && a.Status != store.AlarmAcknowledged {
		return ErrAlarmNotActive
	}
	a.Status = store.AlarmResolved
	a.ResolvedAt = at
	return state.PutAlarm(a)
}

// ResolveIfOpen resolves the latest active alarm of a point when a sample is
// back inside the threshold.
func ResolveIfOpen(state *store.State, pointID string, at time.Time) error {
	latest, ok := state.AlarmByPoint(pointID)
	if !ok {
		return nil
	}
	if latest.Status != store.AlarmOpen && latest.Status != store.AlarmAcknowledged {
		return nil
	}
	return Resolve(state, latest.ID, at)
}
