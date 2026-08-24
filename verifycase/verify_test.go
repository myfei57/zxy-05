package verifycase

import (
	"testing"
	"time"

	"telemetryguard/internal/alarm"
	"telemetryguard/internal/escalate"
	"telemetryguard/internal/store"
)

func TestReopenKeepsOriginalOpenTime(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	t0 := time.Now().UTC()
	t1 := t0.Add(2 * time.Hour)
	existing := &store.Alarm{ID: "a1", PointID: "p1", PointName: "temp_main", DeviceID: "d1", Status: store.AlarmResolved, OriginalOpenAt: t0, LastTriggerAt: t0}
	if err := state.PutAlarm(existing); err != nil {
		t.Fatal(err)
	}
	sample := &store.Sample{PointID: "p1", PointName: "temp_main", DeviceID: "d1", Seq: 2, Value: 95, At: t1}
	entry := store.RuleSnapshotEntry{RuleID: "r1", PointName: "temp_main", Op: ">", Threshold: 80, Level: 1, State: store.RulePublished}
	a, err := alarm.Reopen(state, existing, sample, entry, t1)
	if err != nil {
		t.Fatal(err)
	}
	if !a.OriginalOpenAt.Equal(t0) {
		t.Fatalf("original open time reset to %v, want %v", a.OriginalOpenAt, t0)
	}
	if !escalate.DueAt(a, 1).Equal(t0.Add(30 * time.Minute)) {
		t.Fatalf("SLA deadline %v, want %v", escalate.DueAt(a, 1), t0.Add(30*time.Minute))
	}
}
