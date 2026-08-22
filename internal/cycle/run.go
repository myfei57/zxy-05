package cycle

import (
	"time"

	"telemetryguard/internal/escalate"
	"telemetryguard/internal/eval"
	"telemetryguard/internal/ingest"
	"telemetryguard/internal/schedule"
	"telemetryguard/internal/store"
)

// RunCycle flushes pending samples, evaluates them against the effective rule
// snapshot, commits the batch durably, and only then advances the cursor.
func RunCycle(state *store.State, at time.Time) ([]*store.Alarm, error) {
	samples := ingest.FlushPending()
	if len(samples) == 0 {
		return nil, nil
	}
	batch, err := schedule.NewBatch(state, len(samples))
	if err != nil {
		return nil, err
	}
	opened, err := eval.EvaluateSamples(state, samples, at)
	if err != nil {
		return nil, err
	}
	batch.State = store.BatchCommitted
	batch.CommittedAt = at
	if err := state.PutBatch(batch); err != nil {
		return nil, err
	}
	if err := schedule.AdvanceCursor(state, batch); err != nil {
		return nil, err
	}
	for _, a := range opened {
		if err := escalate.EnsureTimer(state, a, at); err != nil {
			return nil, err
		}
	}
	return opened, nil
}
