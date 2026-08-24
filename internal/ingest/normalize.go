package ingest

import (
	"errors"
	"math"

	"telemetryguard/internal/store"
)

// ErrUnknownPoint is returned when a sample references a missing point.
var ErrUnknownPoint = errors.New("sample point does not exist")

// Normalize rounds a raw reading to three decimals.
func Normalize(value float64) float64 {
	return math.Round(value*1000) / 1000
}

// resolvePoint loads the point and verifies it and its device are enabled.
func resolvePoint(state *store.State, pointID string) (*store.Point, error) {
	p, ok := state.Point(pointID)
	if !ok {
		return nil, ErrUnknownPoint
	}
	if !p.Enabled {
		return nil, errors.New("point is disabled")
	}
	d, ok := state.Device(p.DeviceID)
	if !ok || d.Disabled {
		return nil, errors.New("device is disabled")
	}
	return p, nil
}
