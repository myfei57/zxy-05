package verifycase

import (
	"testing"
	"time"

	"telemetryguard/internal/ingest"
	"telemetryguard/internal/point"
	"telemetryguard/internal/store"
)

func TestDedupWindowRejectsLateSequence(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	d, err := point.RegisterDevice(state, "furnace-1")
	if err != nil {
		t.Fatal(err)
	}
	p, err := point.RegisterPoint(state, d.ID, "temp_main", "c")
	if err != nil {
		t.Fatal(err)
	}
	t0 := time.Now().UTC()
	if _, err := ingest.Receive(state, p.ID, 5, 80, t0); err != nil {
		t.Fatal(err)
	}
	if _, err := ingest.Receive(state, p.ID, 5, 80, t0.Add(61*time.Second)); err == nil {
		t.Fatal("late duplicate sequence accepted after window rollover")
	}
}
