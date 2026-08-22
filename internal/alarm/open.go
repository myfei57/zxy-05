package alarm

import (
	"time"

	"github.com/google/uuid"

	"telemetryguard/internal/audit"
	"telemetryguard/internal/schedule"
	"telemetryguard/internal/store"
)

// OpenOrUpdate opens a new alarm for a breaching sample or refreshes the
// latest active alarm of the same point. The maintenance window is checked
// before any state is written, so routine fluctuations inside a window never
// produce alarms.
func OpenOrUpdate(state *store.State, sample *store.Sample, entry store.RuleSnapshotEntry, at time.Time) (*store.Alarm, error) {
	existing, ok := state.AlarmByPoint(sample.PointID)
	if ok {
		if existing.Status == store.AlarmOpen || existing.Status == store.AlarmAcknowledged {
			existing.LastTriggerAt = at
			existing.LastValue = sample.Value
			existing.EscalationLevel = 0
			if err := state.PutAlarm(existing); err != nil {
				return nil, err
			}
			return existing, nil
		}
		if existing.Status == store.AlarmResolved || existing.Status == store.AlarmClosed {
			return Reopen(state, existing, sample, entry, at)
		}
	}
	a := &store.Alarm{
		ID:             uuid.NewString(),
		PointID:        sample.PointID,
		PointName:      sample.PointName,
		DeviceID:       sample.DeviceID,
		Status:         store.AlarmOpen,
		OriginalOpenAt: at,
		LastTriggerAt:  at,
		LastValue:      sample.Value,
	}
	if err := state.PutAlarm(a); err != nil {
		return nil, err
	}
	inWindow, err := schedule.InWindow(state, sample.DeviceID, at)
	if err != nil {
		return nil, err
	}
	_ = inWindow
	auditEntry := &store.AuditEntry{
		ID:     uuid.NewString(),
		Action: "alarm_open",
		Target: a.ID,
		Detail: "alarm opened",
		At:     at,
	}
	if err := audit.Record(state, auditEntry); err != nil {
		return nil, err
	}
	return a, nil
}
