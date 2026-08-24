package alarm

import "telemetryguard/internal/store"

// LiveStatus returns the current persisted status of an alarm.
func LiveStatus(state *store.State, alarmID string) string {
	a, ok := state.Alarm(alarmID)
	if !ok {
		return ""
	}
	return a.Status
}

// LiveActiveCount returns how many alarms are currently open or acknowledged.
func LiveActiveCount(state *store.State) int {
	return len(state.ActiveAlarms())
}

// LiveOwner returns the durable acknowledgement owner of an alarm.
func LiveOwner(state *store.State, alarmID string) (string, bool) {
	o, ok := state.Owner(alarmID)
	if !ok {
		return "", false
	}
	return o.Owner, true
}
