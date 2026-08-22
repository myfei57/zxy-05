package rule

import "telemetryguard/internal/store"

// PublishVersion publishes all draft rules as a new version and makes it
// effective. The snapshot is durably stored and the version recorded before
// the effective pointer moves, so a failed write never leaves an effective
// version without a complete snapshot.
func PublishVersion(state *store.State) (int64, error) {
	version := NextVersion(state)
	drafts := state.DraftRules()
	for _, r := range drafts {
		r.State = store.RulePublished
		r.Version = version
		if err := state.PutRule(r); err != nil {
			return 0, err
		}
	}
	entries := BuildSnapshot(state)
	if err := state.SetEffective(version); err != nil {
		return 0, err
	}
	if err := state.SaveSnapshot(version, entries); err != nil {
		return version, err
	}
	if err := RecordVersion(state, version); err != nil {
		return version, err
	}
	return version, nil
}
