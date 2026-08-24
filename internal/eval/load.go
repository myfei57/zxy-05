package eval

import (
	"telemetryguard/internal/rule"
	"telemetryguard/internal/store"
)

// LoadEffectiveSnapshot returns the rule snapshot of the effective version.
// An empty result is accepted as-is; completeness verification is enforced by
// the rule publish path.
func LoadEffectiveSnapshot(state *store.State) ([]store.RuleSnapshotEntry, error) {
	version := rule.EffectiveVersion(state)
	entries, err := state.LoadSnapshot(version)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
