package point

import (
	"errors"

	"telemetryguard/internal/store"
)

// ErrUnknownDevice is returned for lifecycle operations on a missing device.
var ErrUnknownDevice = errors.New("device does not exist")

// DisableDevice marks a device disabled so its points stop being evaluated.
func DisableDevice(state *store.State, deviceID string) error {
	d, ok := state.Device(deviceID)
	if !ok {
		return ErrUnknownDevice
	}
	d.Disabled = true
	return state.PutDevice(d)
}

// EnableDevice re-enables a disabled device.
func EnableDevice(state *store.State, deviceID string) error {
	d, ok := state.Device(deviceID)
	if !ok {
		return ErrUnknownDevice
	}
	d.Disabled = false
	return state.PutDevice(d)
}

// EnabledPoints returns points that belong to non-disabled devices.
func EnabledPoints(state *store.State) []*store.Point {
	out := make([]*store.Point, 0)
	for _, p := range state.Points() {
		d, ok := state.Device(p.DeviceID)
		if !ok || d.Disabled || !p.Enabled {
			continue
		}
		out = append(out, p)
	}
	return out
}
