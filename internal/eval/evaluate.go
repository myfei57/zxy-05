package eval

import (
	"time"

	"telemetryguard/internal/alarm"
	"telemetryguard/internal/store"
)

// ThresholdMet reports whether a value breaches a rule threshold.
func ThresholdMet(entry store.RuleSnapshotEntry, value float64) bool {
	switch entry.Op {
	case ">":
		return value > entry.Threshold
	case ">=":
		return value >= entry.Threshold
	case "<":
		return value < entry.Threshold
	case "<=":
		return value <= entry.Threshold
	}
	return false
}

// EvaluateSamples checks a batch of samples against the effective rule snapshot
// and opens or resolves alarms. Draft entries are skipped so unpublished rules
// never affect evaluation.
func EvaluateSamples(state *store.State, samples []*store.Sample, at time.Time) ([]*store.Alarm, error) {
	entries, err := LoadEffectiveSnapshot(state)
	if err != nil {
		return nil, err
	}
	opened := make([]*store.Alarm, 0)
	for _, sample := range samples {
		entry, ok := matchRule(entries, sample.PointName)
		if !ok {
			continue
		}
		if entry.State == store.RuleDraft {
			continue
		}
		if ThresholdMet(entry, sample.Value) {
			a, err := alarm.OpenOrUpdate(state, sample, entry, at)
			if err != nil {
				return nil, err
			}
			opened = append(opened, a)
		} else {
			if err := alarm.ResolveIfOpen(state, sample.PointID, at); err != nil {
				return nil, err
			}
		}
	}
	return opened, nil
}

func matchRule(entries []store.RuleSnapshotEntry, pointName string) (store.RuleSnapshotEntry, bool) {
	for _, e := range entries {
		if e.PointName == pointName {
			return e, true
		}
	}
	return store.RuleSnapshotEntry{}, false
}
