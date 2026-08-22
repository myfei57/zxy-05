package audit

import "telemetryguard/internal/store"

// CountActive returns the total number of alarm records, including closed ones.
func CountActive(state *store.State) int {
	return len(state.Alarms())
}
