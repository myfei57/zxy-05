package point

import "telemetryguard/internal/store"

// Counts summarizes device and point totals.
func Counts(state *store.State) (devices int, points int, disabled int) {
	for _, d := range state.Devices() {
		devices++
		if d.Disabled {
			disabled++
		}
	}
	points = len(state.Points())
	return devices, points, disabled
}
