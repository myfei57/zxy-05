package rule

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"telemetryguard/internal/store"
)

// newTestState returns a State backed by a fresh temp dir.
func newTestState(t *testing.T) *store.State {
	t.Helper()
	return store.NewState(t.TempDir())
}

// addDraftRule creates a draft threshold rule for the given point name.
func addDraftRule(t *testing.T, state *store.State, name, point string) {
	t.Helper()
	if _, err := CreateRule(state, name, point, ">", 80.0, 1); err != nil {
		t.Fatalf("CreateRule: %v", err)
	}
}

func TestPublishVersion_Success_SnapshotCompleteBeforeEffective(t *testing.T) {
	state := newTestState(t)
	addDraftRule(t, state, "temp-high", "motor-temp")
	addDraftRule(t, state, "press-low", "pump-pressure")

	version, err := PublishVersion(state)
	if err != nil {
		t.Fatalf("PublishVersion: %v", err)
	}
	if version != state.Effective() {
		t.Fatalf("effective = %d, want published %d", state.Effective(), version)
	}

	// The effective snapshot must be loadable and complete, and it must match
	// exactly what BuildSnapshot produced — no truncation, no dropped entries.
	loaded, err := state.LoadSnapshot(version)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	want := BuildSnapshot(state)
	if len(loaded) != len(want) {
		t.Fatalf("snapshot entries = %d, want %d", len(loaded), len(want))
	}
	for i := range want {
		if loaded[i] != want[i] {
			t.Fatalf("entry %d = %+v, want %+v", i, loaded[i], want[i])
		}
	}
	if err := VerifySnapshot(state, version, want); err != nil {
		t.Fatalf("VerifySnapshot after publish: %v", err)
	}

	// All published rules carry the new version and published state.
	for _, r := range state.Rules() {
		if r.Version != version || r.State != store.RulePublished {
			t.Fatalf("rule %s version=%d state=%s, want version=%d published", r.Name, r.Version, r.State, version)
		}
	}
}

func TestPublishVersion_SnapshotWriteFails_EffectiveUnchanged(t *testing.T) {
	state := newTestState(t)
	addDraftRule(t, state, "temp-high", "motor-temp")

	// Poison the snapshots directory path so SaveSnapshot (which calls MkdirAll
	// on <root>/snapshots) fails with "not a directory" — mirroring a write
	// failure such as disk full or the historical snapshot write failure.
	snapDir := filepath.Join(state.Root(), "snapshots")
	if err := os.WriteFile(snapDir, []byte("poison"), 0o644); err != nil {
		t.Fatalf("poison snapshots dir: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(snapDir) })

	prevEffective := state.Effective()
	version, err := PublishVersion(state)
	if err == nil {
		t.Fatalf("PublishVersion unexpectedly succeeded, version=%d", version)
	}
	// The effective pointer must not have moved to the un-snapshotted version.
	if got := state.Effective(); got != prevEffective {
		t.Fatalf("effective = %d, want unchanged %d", got, prevEffective)
	}
	// Drafts must be restored so the publish can be retried.
	for _, r := range state.Rules() {
		if r.State != store.RuleDraft {
			t.Fatalf("rule %s state=%s, want draft after failed publish", r.Name, r.State)
		}
		if r.Version != 0 {
			t.Fatalf("rule %s version=%d, want 0 after failed publish", r.Name, r.Version)
		}
	}
	// No snapshot file should exist for the never-effective version.
	if got, _ := state.LoadSnapshot(version); len(got) != 0 {
		t.Fatalf("snapshot for failed version %d should be absent, got %d entries", version, len(got))
	}
}

func TestVerifySnapshot_RejectsTruncatedFile(t *testing.T) {
	state := newTestState(t)
	addDraftRule(t, state, "temp-high", "motor-temp")

	version, err := PublishVersion(state)
	if err != nil {
		t.Fatalf("PublishVersion: %v", err)
	}
	want := BuildSnapshot(state)

	// Overwrite the snapshot with truncated content and confirm VerifySnapshot
	// rejects it: this is the guardrail that keeps a half-written file from
	// being treated as effective.
	snapPath := filepath.Join(state.Root(), "snapshots", strconv.FormatInt(version, 10)+".json")
	if err := os.WriteFile(snapPath, []byte("[]"), 0o644); err != nil {
		t.Fatalf("truncate snapshot: %v", err)
	}
	if err := VerifySnapshot(state, version, want); err != ErrSnapshotIncomplete {
		t.Fatalf("VerifySnapshot truncated = %v, want ErrSnapshotIncomplete", err)
	}
}
