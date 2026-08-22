package point

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"telemetryguard/internal/store"
)

// ErrMissingDevice is returned when a point references an unknown device.
var ErrMissingDevice = errors.New("point device does not exist")

// ErrMissingPoint is returned when a point is not registered.
var ErrMissingPoint = errors.New("point does not exist")

// RegisterDevice creates a new active device.
func RegisterDevice(state *store.State, name string) (*store.Device, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("device name is empty")
	}
	d := &store.Device{
		ID:        uuid.NewString(),
		Name:      name,
		Disabled:  false,
		CreatedAt: time.Now().UTC(),
	}
	if err := state.PutDevice(d); err != nil {
		return nil, err
	}
	return d, nil
}

// RegisterPoint creates an enabled telemetry point on a device.
func RegisterPoint(state *store.State, deviceID, name, unit string) (*store.Point, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("point name is empty")
	}
	if _, ok := state.Device(deviceID); !ok {
		return nil, ErrMissingDevice
	}
	p := &store.Point{
		ID:       uuid.NewString(),
		DeviceID: deviceID,
		Name:     name,
		Unit:     unit,
		Enabled:  true,
	}
	if err := state.PutPoint(p); err != nil {
		return nil, err
	}
	return p, nil
}

// AdvanceLastSeq durably records the latest accepted sequence and value of a point.
func AdvanceLastSeq(state *store.State, pointID string, seq int64, value float64, at time.Time) (*store.Point, error) {
	p, ok := state.Point(pointID)
	if !ok {
		return nil, ErrMissingPoint
	}
	p.LastSeq = seq
	p.LastValue = value
	p.LastAt = at
	if err := state.PutPoint(p); err != nil {
		return nil, err
	}
	return p, nil
}
