package cycle

import (
	"time"

	"telemetryguard/internal/escalate"
	"telemetryguard/internal/eval"
	"telemetryguard/internal/ingest"
	"telemetryguard/internal/store"
)

// HandleSample ingests one sample, evaluates it immediately, and starts any
// required escalation timer.
func HandleSample(state *store.State, pointID string, seq int64, value float64, at time.Time) ([]*store.Alarm, error) {
	sample, err := ingest.Receive(state, pointID, seq, value, at)
	if err != nil {
		return nil, err
	}
	opened, err := eval.EvaluateSamples(state, []*store.Sample{sample}, at)
	if err != nil {
		return nil, err
	}
	for _, a := range opened {
		if err := escalate.EnsureTimer(state, a, at); err != nil {
			return nil, err
		}
	}
	return opened, nil
}
