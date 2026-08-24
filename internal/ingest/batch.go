package ingest

import (
	"time"

	"telemetryguard/internal/store"
)

// pending holds samples waiting for the next scheduled evaluation batch.
var pending []*store.Sample

// FlushPending returns and clears the pending sample buffer.
func FlushPending() []*store.Sample {
	out := pending
	pending = nil
	return out
}

// AppendPending queues a sample for the next batch.
func AppendPending(s *store.Sample) {
	pending = append(pending, s)
}

// PendingCount reports how many samples are buffered.
func PendingCount() int {
	return len(pending)
}

// NowUTC is the platform clock used by the control plane.
func NowUTC() time.Time {
	return time.Now().UTC()
}
