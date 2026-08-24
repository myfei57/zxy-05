package alarm

import (
	"errors"
	"time"

	"telemetryguard/internal/store"
)

// ErrAlarmNotOpen is returned when acknowledging an alarm that is not open.
var ErrAlarmNotOpen = errors.New("alarm is not open")

// Acknowledge marks an alarm as acknowledged by an operator. The owner
// attribution is durably stored before the alarm state advances, so a failed
// owner write leaves the alarm open with no false confirmation.
func Acknowledge(state *store.State, alarmID, owner, reason string, at time.Time) error {
	a, ok := state.Alarm(alarmID)
	if !ok {
		return errors.New("alarm does not exist")
	}
	if a.Status != store.AlarmOpen {
		return ErrAlarmNotOpen
	}
	a.Status = store.AlarmAcknowledged
	a.AckedAt = at
	a.Owner = owner
	a.Reason = reason
	if err := state.PutAlarm(a); err != nil {
		return err
	}
	record := &store.AlarmOwner{
		AlarmID: alarmID,
		Owner:   owner,
		Reason:  reason,
		At:      at,
	}
	if err := state.PutOwner(record); err != nil {
		return err
	}
	return nil
}
