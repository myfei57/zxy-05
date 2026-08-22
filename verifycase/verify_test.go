package verifycase

import (
	"testing"
	"time"

	"telemetryguard/internal/alarm"
	"telemetryguard/internal/channel"
	"telemetryguard/internal/escalate"
	"telemetryguard/internal/store"
)

func TestEscalationReadsLiveAlarmState(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	if _, err := channel.RegisterChannel(state, "sms", 1); err != nil {
		t.Fatal(err)
	}
	t0 := time.Now().UTC().Add(-2 * time.Hour)
	al := &store.Alarm{ID: "a1", PointID: "p1", PointName: "temp_main", DeviceID: "d1", Status: store.AlarmOpen, OriginalOpenAt: t0, LastTriggerAt: t0}
	if err := state.PutAlarm(al); err != nil {
		t.Fatal(err)
	}
	if err := escalate.EnsureTimer(state, al, t0); err != nil {
		t.Fatal(err)
	}
	if err := alarm.Acknowledge(state, "a1", "zhao", "handled", t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	notifs, err := escalate.Tick(state, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if len(notifs) != 0 {
		t.Fatalf("escalation fired %d notifications for an acknowledged alarm", len(notifs))
	}
}
