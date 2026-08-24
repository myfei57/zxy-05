package verifycase

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"telemetryguard/internal/alarm"
	"telemetryguard/internal/store"
)

func TestAckWaitsForOwnerDurable(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	t0 := time.Now().UTC()
	al := &store.Alarm{ID: "alarm-1", PointID: "p1", PointName: "temp_main", DeviceID: "d1", Status: store.AlarmOpen, OriginalOpenAt: t0, LastTriggerAt: t0}
	if err := state.PutAlarm(al); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "owners"), []byte("blocker"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := alarm.Acknowledge(state, "alarm-1", "zhao", "check", t0.Add(time.Minute)); err == nil {
		t.Fatal("expected the owner write to fail")
	}
	got, _ := state.Alarm("alarm-1")
	if got.Status != store.AlarmOpen {
		t.Fatalf("alarm status = %s, want %s", got.Status, store.AlarmOpen)
	}
}
