package channel

import (
	"errors"

	"telemetryguard/internal/store"
)

// ErrNoChannel is returned when the route table is empty.
var ErrNoChannel = errors.New("no notification channel registered")

// CurrentRoute resolves the live route from channels sorted by priority.
func CurrentRoute(state *store.State) (string, error) {
	for _, c := range state.ChannelByPriority() {
		if c.Enabled && c.Reachable {
			if err := state.SetCurrentRoute(c.ID); err != nil {
				return "", err
			}
			return c.ID, nil
		}
	}
	return "", nil
}

// CachedRoute returns the persisted route pointer without revalidation.
func CachedRoute(state *store.State) string {
	return state.CurrentRoute()
}

// Failover marks a channel unreachable and re-resolves the current route.
func Failover(state *store.State, channelID string) error {
	c, ok := state.Channel(channelID)
	if !ok {
		return errors.New("channel does not exist")
	}
	c.Reachable = false
	return state.PutChannel(c)
}
