package verifycase

import (
	"testing"
	"time"

	"telemetryguard/internal/alarm"
	"telemetryguard/internal/schedule"
	"telemetryguard/internal/store"
)

func TestWindowCheckedBeforeAlarmOpen(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	now := time.Now().UTC()
	if _, err := schedule.AddWindow(state, "d1", now.Add(-time.Hour), now.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	sample := &store.Sample{PointID: "p1", PointName: "temp_main", DeviceID: "d1", Seq: 1, Value: 95, At: now}
	entry := store.RuleSnapshotEntry{RuleID: "r1", PointName: "temp_main", Op: ">", Threshold: 80, Level: 1, State: store.RulePublished}
	a, err := alarm.OpenOrUpdate(state, sample, entry, now)
	if err != nil {
		t.Fatal(err)
	}
	if a != nil {
		t.Fatal("alarm returned inside maintenance window")
	}
	if len(state.Alarms()) != 0 {
		t.Fatalf("%d alarms opened inside maintenance window", len(state.Alarms()))
	}
}
