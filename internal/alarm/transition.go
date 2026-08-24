package alarm

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"telemetryguard/internal/audit"
	"telemetryguard/internal/store"
)

// IsActive reports whether an alarm still needs operator attention.
func IsActive(a *store.Alarm) bool {
	return a != nil && (a.Status == store.AlarmOpen || a.Status == store.AlarmAcknowledged)
}

// ShouldEscalate reports whether an alarm still needs escalation attention.
// Only freshly open alarms escalate; acknowledged or resolved alarms do not.
func ShouldEscalate(a *store.Alarm) bool {
	return a != nil && a.Status == store.AlarmOpen
}

// CanAcknowledge reports whether an alarm may be acknowledged.
func CanAcknowledge(a *store.Alarm) bool {
	return a != nil && a.Status == store.AlarmOpen
}

// CanResolve reports whether an alarm may be resolved.
func CanResolve(a *store.Alarm) bool {
	return a != nil && (a.Status == store.AlarmOpen || a.Status == store.AlarmAcknowledged)
}

// Close closes an alarm record after it was resolved.
func Close(state *store.State, alarmID string, at time.Time) error {
	a, ok := state.Alarm(alarmID)
	if !ok {
		return errors.New("alarm does not exist")
	}
	a.Status = store.AlarmClosed
	if err := state.PutAlarm(a); err != nil {
		return err
	}
	entry := &store.AuditEntry{
		ID:     uuid.NewString(),
		Action: "alarm_close",
		Target: alarmID,
		Detail: "alarm closed",
		At:     at,
	}
	return audit.Record(state, entry)
}
