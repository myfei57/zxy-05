package eval

import (
	"testing"
	"time"

	"telemetryguard/internal/rule"
	"telemetryguard/internal/schedule"
	"telemetryguard/internal/store"
)

// TestEvaluateSamplesSuppressesBreachesInsideMaintenanceWindow is the
// end-to-end guard for the reported bug: when a device is inside a
// maintenance window, threshold-breaching samples must produce no alarm,
// no audit entry, and the caller must receive an empty opened slice so no
// SLA escalation timer is started.
func TestEvaluateSamplesSuppressesBreachesInsideMaintenanceWindow(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)

	const deviceID = "dev-1"
	const pointID = "pt-1"
	const pointName = "temp"
	if _, err := rule.CreateRule(state, "r", pointName, ">", 100, 1); err != nil {
		t.Fatalf("CreateRule: %v", err)
	}
	if _, err := rule.PublishVersion(state); err != nil {
		t.Fatalf("PublishVersion: %v", err)
	}

	at := time.Date(2026, 8, 22, 10, 0, 0, 0, time.UTC)
	if _, err := schedule.AddWindow(state, deviceID, at.Add(-time.Hour), at.Add(time.Hour)); err != nil {
		t.Fatalf("AddWindow: %v", err)
	}

	samples := []*store.Sample{{
		PointID:   pointID,
		PointName: pointName,
		DeviceID:  deviceID,
		Value:     250, // breaches "> 100"
		At:        at,
	}}

	opened, err := EvaluateSamples(state, samples, at)
	if err != nil {
		t.Fatalf("EvaluateSamples: %v", err)
	}
	if len(opened) != 0 {
		t.Fatalf("expected no alarms opened inside maintenance window, got %d: %+v", len(opened), opened)
	}
	if alarms := state.Alarms(); len(alarms) != 0 {
		t.Fatalf("expected no alarm records persisted, got %d", len(alarms))
	}
	if entries := state.AuditEntries(); len(entries) != 0 {
		t.Fatalf("expected no audit entries persisted, got %d", len(entries))
	}
	if es := state.Escalations(); len(es) != 0 {
		t.Fatalf("expected no escalation timers created, got %d", len(es))
	}
}

// TestEvaluateSamplesOpensBreachesOutsideMaintenanceWindow confirms the
// happy path still works once the device leaves its window.
func TestEvaluateSamplesOpensBreachesOutsideMaintenanceWindow(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)

	const deviceID = "dev-1"
	const pointID = "pt-1"
	const pointName = "temp"
	if _, err := rule.CreateRule(state, "r", pointName, ">", 100, 1); err != nil {
		t.Fatalf("CreateRule: %v", err)
	}
	if _, err := rule.PublishVersion(state); err != nil {
		t.Fatalf("PublishVersion: %v", err)
	}

	// A window that already ended before evaluation time.
	at := time.Date(2026, 8, 22, 10, 0, 0, 0, time.UTC)
	if _, err := schedule.AddWindow(state, deviceID, at.Add(-2*time.Hour), at.Add(-time.Hour)); err != nil {
		t.Fatalf("AddWindow: %v", err)
	}

	samples := []*store.Sample{{
		PointID:   pointID,
		PointName: pointName,
		DeviceID:  deviceID,
		Value:     250,
		At:        at,
	}}

	opened, err := EvaluateSamples(state, samples, at)
	if err != nil {
		t.Fatalf("EvaluateSamples: %v", err)
	}
	if len(opened) != 1 {
		t.Fatalf("expected one alarm opened, got %d", len(opened))
	}
	if opened[0].Status != store.AlarmOpen {
		t.Fatalf("expected open alarm, got %q", opened[0].Status)
	}
}
