package channel

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"telemetryguard/internal/store"
)

// ErrChannelConflict is returned when a channel name is already registered.
var ErrChannelConflict = errors.New("channel name already exists")

// RegisterChannel creates an enabled, reachable channel with a priority.
func RegisterChannel(state *store.State, name string, priority int) (*store.ChannelConfig, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("channel name is empty")
	}
	for _, c := range state.Channels() {
		if c.Name == name {
			return nil, ErrChannelConflict
		}
	}
	c := &store.ChannelConfig{
		ID:        uuid.NewString(),
		Name:      name,
		Enabled:   true,
		Reachable: true,
		Priority:  priority,
		CreatedAt: time.Now().UTC(),
	}
	if err := state.PutChannel(c); err != nil {
		return nil, err
	}
	return c, nil
}

// SetReachable updates the reachability of a channel.
func SetReachable(state *store.State, channelID string, reachable bool) error {
	c, ok := state.Channel(channelID)
	if !ok {
		return errors.New("channel does not exist")
	}
	c.Reachable = reachable
	return state.PutChannel(c)
}
