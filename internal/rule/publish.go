package rule

import "telemetryguard/internal/store"

// PublishVersion publishes all draft rules as a new version and makes it
// effective. The snapshot is written and read back to verify completeness, and
// the version is recorded, all before the effective pointer moves; only the
// final step swings evaluation to the new version. If any pre-effective step
// fails, draft rules are restored so the publish can be retried and no
// effective version is ever left pointing at a missing or truncated snapshot.
func PublishVersion(state *store.State) (int64, error) {
	version := NextVersion(state)
	drafts := state.DraftRules()
	for _, r := range drafts {
		r.State = store.RulePublished
		r.Version = version
		if err := state.PutRule(r); err != nil {
			revertDrafts(state, drafts)
			return 0, err
		}
	}
	entries := BuildSnapshot(state)

	// Durably store the snapshot and read it back to confirm it is complete.
	// The effective pointer is NOT moved until this succeeds, so a failed or
	// truncated write never leaves evaluation pointing at a half-written file.
	if err := state.SaveSnapshot(version, entries); err != nil {
		revertDrafts(state, drafts)
		return 0, err
	}
	if err := VerifySnapshot(state, version, entries); err != nil {
		state.DeleteSnapshot(version)
		revertDrafts(state, drafts)
		return 0, err
	}
	if err := RecordVersion(state, version); err != nil {
		state.DeleteSnapshot(version)
		revertDrafts(state, drafts)
		return 0, err
	}

	// Last step: now that the snapshot is verified and the version recorded,
	// swing the effective pointer. Failure here leaves a fully written snapshot
	// and version but an unchanged effective pointer, which is safe — callers
	// retry or roll back.
	if err := state.SetEffective(version); err != nil {
		revertDrafts(state, drafts)
		return 0, err
	}
	return version, nil
}

// revertDrafts restores the given rules to draft state after a failed publish,
// so a retry re-publishes the same rules instead of leaving them stranded as
// published-but-never-effective.
func revertDrafts(state *store.State, drafts []*store.Rule) {
	for _, r := range drafts {
		r.State = store.RuleDraft
		r.Version = 0
		_ = state.PutRule(r)
	}
}
