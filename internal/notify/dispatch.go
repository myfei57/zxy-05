package notify

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"telemetryguard/internal/channel"
	"telemetryguard/internal/store"
)

// ErrNoRoute is returned when no reachable channel exists.
var ErrNoRoute = errors.New("no reachable notification channel")

// Dispatch sends one notification for an alarm escalation through the live
// route table.
func Dispatch(state *store.State, alarmID string, level int, at time.Time) (*store.Notification, error) {
	routeID, err := channel.CurrentRoute(state)
	if err != nil {
		return nil, err
	}
	if routeID == "" {
		return nil, ErrNoRoute
	}
	c, ok := state.Channel(routeID)
	if !ok || !c.Reachable {
		return nil, ErrNoRoute
	}
	n := &store.Notification{
		ID:        uuid.NewString(),
		AlarmID:   alarmID,
		ChannelID: routeID,
		Status:    store.NotificationSent,
		Attempts:  1,
		CreatedAt: at,
	}
	if err := state.PutNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}
