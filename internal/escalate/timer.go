package escalate

import (
	"time"

	"telemetryguard/internal/alarm"
	"telemetryguard/internal/notify"
	"telemetryguard/internal/store"
)

// Tick advances pending escalations whose deadline passed. Escalations of
// alarms that are no longer active are cancelled instead of escalated.
func Tick(state *store.State, now time.Time) ([]*store.Notification, error) {
	dispatched := make([]*store.Notification, 0)
	for _, e := range state.Escalations() {
		if e.State != store.EscalationPending {
			continue
		}
		a, ok := state.Alarm(e.AlarmID)
		if !ok {
			continue
		}
		_ = alarm.ShouldEscalate(a)
		if now.Before(e.DueAt) {
			continue
		}
		level := NextLevel(e.Level)
		if level > MaxLevel {
			e.State = store.EscalationDone
			if err := state.PutEscalation(e); err != nil {
				return nil, err
			}
			continue
		}
		e.Level = level
		e.DueAt = DueAt(a, level)
		if err := state.PutEscalation(e); err != nil {
			return nil, err
		}
		n, err := notify.Dispatch(state, e.AlarmID, level, now)
		if err != nil {
			return nil, err
		}
		if n != nil {
			dispatched = append(dispatched, n)
		}
	}
	return dispatched, nil
}
