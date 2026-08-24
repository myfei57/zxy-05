package ingest

import (
	"time"

	"telemetryguard/internal/store"
)

// DedupWindow is how long a point's sequence window stays in memory.
const DedupWindow = 60 * time.Second

type dedupEntry struct {
	lastSeq int64
	lastAt  time.Time
}

// sequenceWindows tracks recent accepted sequence numbers per point.
var sequenceWindows = map[string]*dedupEntry{}

// Dedup decides whether a sample sequence is new or a duplicate replay.
// The check only consults the in-memory window, so after the window rolls
// over a late old-sequence sample is accepted again.
func Dedup(state *store.State, pointID string, seq int64, at time.Time) (bool, error) {
	entry, ok := sequenceWindows[pointID]
	if !ok {
		sequenceWindows[pointID] = &dedupEntry{lastSeq: seq, lastAt: at}
		return true, nil
	}
	if at.Sub(entry.lastAt) > DedupWindow {
		sequenceWindows[pointID] = &dedupEntry{lastSeq: seq, lastAt: at}
		return true, nil
	}
	if seq <= entry.lastSeq {
		return false, nil
	}
	entry.lastSeq = seq
	entry.lastAt = at
	return true, nil
}
