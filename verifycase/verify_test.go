package verifycase

import (
	"testing"
	"time"

	"telemetryguard/internal/alarm"
	"telemetryguard/internal/audit"
	"telemetryguard/internal/store"
)

func TestActiveCountExcludesClosedAlarms(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	t0 := time.Now().UTC()
	sample := &store.Sample{PointID: "p1", PointName: "temp_main", DeviceID: "d1", Seq: 1, Value: 95, At: t0}
	entry := store.RuleSnapshotEntry{RuleID: "r1", PointName: "temp_main", Op: ">", Threshold: 80, Level: 1, State: store.RulePublished}
	a, err := alarm.OpenOrUpdate(state, sample, entry, t0)
	if err != nil {
		t.Fatal(err)
	}
	if err := alarm.Close(state, a.ID, t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if got := audit.CountActive(state); got != 0 {
		t.Fatalf("active count = %d, want 0 after close", got)
	}
}
