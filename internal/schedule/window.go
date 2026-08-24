package schedule

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"telemetryguard/internal/store"
)

// ErrBadWindow is returned when a window has an invalid time span.
var ErrBadWindow = errors.New("maintenance window has invalid time span")

// AddWindow registers a maintenance window for a device.
func AddWindow(state *store.State, deviceID string, start, end time.Time) (*store.MaintenanceWindow, error) {
	if !end.After(start) {
		return nil, ErrBadWindow
	}
	w := &store.MaintenanceWindow{
		ID:       uuid.NewString(),
		DeviceID: deviceID,
		StartAt:  start,
		EndAt:    end,
	}
	if err := state.PutWindow(w); err != nil {
		return nil, err
	}
	return w, nil
}

// InWindow reports whether a device is inside a maintenance window at time at.
func InWindow(state *store.State, deviceID string, at time.Time) (bool, error) {
	w, ok := state.WindowForDevice(deviceID)
	if !ok {
		return false, nil
	}
	if at.Before(w.StartAt) || at.After(w.EndAt) {
		return false, nil
	}
	return true, nil
}
