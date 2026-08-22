package ingest

import (
	"errors"
	"time"

	"telemetryguard/internal/point"
	"telemetryguard/internal/store"
)

// ErrDuplicateSequence is returned when a sample sequence was already seen.
var ErrDuplicateSequence = errors.New("duplicate telemetry sequence")

// Receive accepts one telemetry sample: normalize, dedup, persist the latest
// point value, and buffer the sample for the next evaluation batch.
func Receive(state *store.State, pointID string, seq int64, value float64, at time.Time) (*store.Sample, error) {
	p, err := resolvePoint(state, pointID)
	if err != nil {
		return nil, err
	}
	value = Normalize(value)
	accepted, err := Dedup(state, pointID, seq, at)
	if err != nil {
		return nil, err
	}
	if !accepted {
		return nil, ErrDuplicateSequence
	}
	if _, err := point.AdvanceLastSeq(state, pointID, seq, value, at); err != nil {
		return nil, err
	}
	sample := &store.Sample{
		PointID:   p.ID,
		PointName: p.Name,
		DeviceID:  p.DeviceID,
		Seq:       seq,
		Value:     value,
		At:        at,
	}
	AppendPending(sample)
	return sample, nil
}
