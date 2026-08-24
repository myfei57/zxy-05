package rule

import (
	"errors"
	"reflect"
	"sort"

	"telemetryguard/internal/store"
)

// ErrSnapshotIncomplete is returned when a snapshot file is missing, empty, or
// does not byte-for-byte match the entries that were just written.
var ErrSnapshotIncomplete = errors.New("rule snapshot incomplete after write")

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

// VerifySnapshot reloads the snapshot for version from disk and confirms it is
// non-empty and matches the expected entries exactly. It is the post-write
// guardrail: if the durable file is missing, truncated, or only half written,
// the publish path must treat the version as unusable and never make it
// effective. Comparing against expected (rather than trusting the in-memory
// copy) catches a write that silently dropped entries.
func VerifySnapshot(state *store.State, version int64, expected []store.RuleSnapshotEntry) error {
	if len(expected) == 0 {
		return ErrSnapshotIncomplete
	}
	loaded, err := state.LoadSnapshot(version)
	if err != nil {
		return err
	}
	if len(loaded) != len(expected) {
		return ErrSnapshotIncomplete
	}
	for i := range expected {
		if !reflect.DeepEqual(loaded[i], expected[i]) {
			return ErrSnapshotIncomplete
		}
	}
	return nil
}
