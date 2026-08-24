package verifycase

import (
	"os"
	"path/filepath"
	"testing"

	"telemetryguard/internal/rule"
	"telemetryguard/internal/store"
)

func TestRulePublishWaitsForSnapshotDurable(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "snapshots"), []byte("blocker"), 0o644); err != nil {
		t.Fatal(err)
	}
	state := store.NewState(dir)
	if _, err := rule.CreateRule(state, "high-temp", "temp_main", ">", 80, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := rule.PublishVersion(state); err == nil {
		t.Fatal("expected the snapshot write to fail")
	}
	if state.Effective() != 0 {
		t.Fatalf("effective advanced to %d although snapshot was not durable", state.Effective())
	}
}
