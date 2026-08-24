package rule

import (
	"sort"

	"telemetryguard/internal/store"
)

// BuildSnapshot converts the published rule set into snapshot entries.
func BuildSnapshot(state *store.State) []store.RuleSnapshotEntry {
	rules := state.PublishedRules()
	sort.Slice(rules, func(i, j int) bool { return rules[i].Name < rules[j].Name })
	entries := make([]store.RuleSnapshotEntry, 0, len(rules))
	for _, r := range rules {
		entries = append(entries, store.RuleSnapshotEntry{
			RuleID:    r.ID,
			PointName: r.PointName,
			Op:        r.Op,
			Threshold: r.Threshold,
			Level:     r.Level,
			State:     r.State,
		})
	}
	return entries
}

// SnapshotComplete reports whether a snapshot for the version exists and is not empty.
func SnapshotComplete(state *store.State, version int64) bool {
	entries, err := state.LoadSnapshot(version)
	if err != nil {
		return false
	}
	return len(entries) > 0
}
