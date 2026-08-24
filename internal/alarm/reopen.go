package alarm

import (
	"time"

	"github.com/google/uuid"

	"telemetryguard/internal/audit"
	"telemetryguard/internal/store"
)

// Reopen creates a fresh alarm record when a point breaches again after it
// recovered. The original open time is carried over from the previous record
// so the SLA countdown is never reset by recurring triggers.
func Reopen(state *store.State, existing *store.Alarm, sample *store.Sample, entry store.RuleSnapshotEntry, at time.Time) (*store.Alarm, error) {
	a := &store.Alarm{
		ID:              uuid.NewString(),
		PointID:         sample.PointID,
		PointName:       sample.PointName,
		DeviceID:        sample.DeviceID,
		Status:          store.AlarmOpen,
		OriginalOpenAt:  existing.OriginalOpenAt,
		LastTriggerAt:   at,
		LastValue:       sample.Value,
		EscalationLevel: 0,
	}
	if err := state.PutAlarm(a); err != nil {
		return nil, err
	}
	auditEntry := &store.AuditEntry{
		ID:     uuid.NewString(),
		Action: "alarm_open",
		Target: a.ID,
		Detail: "alarm reopened",
		At:     at,
	}
	if err := audit.Record(state, auditEntry); err != nil {
		return nil, err
	}
	return a, nil
}
