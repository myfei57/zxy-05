package console

import (
	"telemetryguard/internal/ingest"
	"telemetryguard/internal/point"
	"telemetryguard/internal/store"
)

func pointCounts(state *store.State) (int, int, int) {
	return point.Counts(state)
}

func ingestPendingCount() int {
	return ingest.PendingCount()
}

func notifyFailed(state *store.State) []*store.Notification {
	return state.FailedNotifications()
}
