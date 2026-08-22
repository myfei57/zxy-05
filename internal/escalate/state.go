package escalate

import (
	"time"

	"telemetryguard/internal/store"
)

// EnsureTimer creates or refreshes the escalation timer for an alarm.
func EnsureTimer(state *store.State, a *store.Alarm, at time.Time) error {
	e := &store.Escalation{
		AlarmID: a.ID,
		Level:   a.EscalationLevel,
		DueAt:   DueAt(a, a.EscalationLevel),
		State:   store.EscalationPending,
	}
	return state.PutEscalation(e)
}

// OwnerFor returns the durable acknowledgement owner of an alarm.
func OwnerFor(state *store.State, alarmID string) (string, bool) {
	o, ok := state.Owner(alarmID)
	if !ok {
		return "", false
	}
	return o.Owner, true
}
