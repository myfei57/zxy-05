package notify

import (
	"time"

	"telemetryguard/internal/channel"
	"telemetryguard/internal/store"
)

// RetryFailed re-dispatches failed notifications through the current route.
func RetryFailed(state *store.State, at time.Time) int {
	requeued := 0
	for _, n := range state.FailedNotifications() {
		routeID, err := channel.CurrentRoute(state)
		if err != nil || routeID == "" {
			continue
		}
		n.Attempts++
		n.ChannelID = routeID
		n.Status = store.NotificationSent
		n.Error = ""
		if err := state.PutNotification(n); err != nil {
			continue
		}
		requeued++
	}
	return requeued
}
