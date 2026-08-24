package rule

import (
	"time"

	"telemetryguard/internal/store"
)

// NextVersion returns the next rule version number.
func NextVersion(state *store.State) int64 {
	next := int64(1)
	for _, v := range state.Versions() {
		if v.Version >= next {
			next = v.Version + 1
		}
	}
	return next
}

// RecordVersion marks a version number as published.
func RecordVersion(state *store.State, version int64) error {
	v := &store.RuleVersion{
		Version:   version,
		State:     store.RulePublished,
		CreatedAt: time.Now().UTC(),
	}
	return state.PutVersion(v)
}

// EffectiveVersion returns the currently effective version.
func EffectiveVersion(state *store.State) int64 {
	return state.Effective()
}
