package ingest

import (
	"time"

	"telemetryguard/internal/store"
)

// Dedup rejects any sample whose sequence is not newer than the point's
// durably persisted last sequence, so late old-sequence samples are never
// re-evaluated even after the in-memory window rolls over.
func Dedup(state *store.State, pointID string, seq int64, at time.Time) (bool, error) {
	p, ok := state.Point(pointID)
	if ok && seq <= p.LastSeq {
		return false, nil
	}
	return true, nil
}
