package alarm

import (
	"testing"
	"time"

	"telemetryguard/internal/rule"
	"telemetryguard/internal/schedule"
	"telemetryguard/internal/store"
)

// newTestState builds a fresh file-backed state rooted in a temp directory and
// returns it along with a cleanup function.
func newTestState(t *testing.T) (*store.State, func()) {
	t.Helper()
	dir := t.TempDir()
	s := store.NewState(dir)
	return s, func() {}
}

// publishRule publishes a single threshold rule for pointName so that the
// effective snapshot contains exactly one entry.
func publishRule(t *testing.T, s *store.State, pointName, op string, threshold float64) {
	t.Helper()
	if _, err := rule.CreateRule(s, "r", pointName, op, threshold, 1); err != nil {
		t.Fatalf("CreateRule: %v", err)
	}
	if _, err := rule.PublishVersion(s); err != nil {
		t.Fatalf("PublishVersion: %v", err)
	}
}

func TestOpenOrUpdateSuppressesBreachInsideMaintenanceWindow(t *testing.T) {
	state, cleanup := newTestState(t)
	defer cleanup()

	const deviceID = "dev-1"
	const pointID = "pt-1"
	const pointName = "temp"
	publishRule(t, state, pointName, ">", 100)

	// Device is inside a maintenance window covering the evaluation time.
	at := time.Date(2026, 8, 22, 10, 0, 0, 0, time.UTC)
	if _, err := schedule.AddWindow(state, deviceID, at.Add(-time.Hour), at.Add(time.Hour)); err != nil {
		t.Fatalf("AddWindow: %v", err)
	}

	entry := mustLoadSnapshot(t, state)
	sample := &store.Sample{
		PointID:   pointID,
		PointName: pointName,
		DeviceID:  deviceID,
		Value:     250, // breaches the "> 100" rule
		At:        at,
	}

	a, err := OpenOrUpdate(state, sample, entry[0], at)
	if err != nil {
		t.Fatalf("OpenOrUpdate: %v", err)
	}
	if a != nil {
		t.Fatalf("expected breach inside maintenance window to be suppressed (nil alarm), got %v", a)
	}
	// No alarm record and no audit entry must have been written.
	if alarms := state.Alarms(); len(alarms) != 0 {
		t.Fatalf("expected no alarms persisted, got %d", len(alarms))
	}
	if entries := state.AuditEntries(); len(entries) != 0 {
		t.Fatalf("expected no audit entries persisted, got %d", len(entries))
	}
}

func TestOpenOrUpdateOpensBreachOutsideMaintenanceWindow(t *testing.T) {
	state, cleanup := newTestState(t)
	defer cleanup()

	const deviceID = "dev-1"
	const pointID = "pt-1"
	const pointName = "temp"
	publishRule(t, state, pointName, ">", 100)

	// No maintenance window at all: a breach must open an alarm.
	at := time.Date(2026, 8, 22, 10, 0, 0, 0, time.UTC)
	entry := mustLoadSnapshot(t, state)
	sample := &store.Sample{
		PointID:   pointID,
		PointName: pointName,
		DeviceID:  deviceID,
		Value:     250,
		At:        at,
	}

	a, err := OpenOrUpdate(state, sample, entry[0], at)
	if err != nil {
		t.Fatalf("OpenOrUpdate: %v", err)
	}
	if a == nil {
		t.Fatal("expected a breach outside any window to open an alarm, got nil")
	}
	if a.Status != store.AlarmOpen {
		t.Fatalf("expected open alarm, got %q", a.Status)
	}
}

func TestOpenOrUpdateSuppressesRefreshInsideMaintenanceWindow(t *testing.T) {
	state, cleanup := newTestState(t)
	defer cleanup()

	const deviceID = "dev-1"
	const pointID = "pt-1"
	const pointName = "temp"
	publishRule(t, state, pointName, ">", 100)

	// Open a real alarm first, outside any window.
	t0 := time.Date(2026, 8, 22, 9, 0, 0, 0, time.UTC)
	entry := mustLoadSnapshot(t, state)
	first := &store.Sample{PointID: pointID, PointName: pointName, DeviceID: deviceID, Value: 250, At: t0}
	if _, err := OpenOrUpdate(state, first, entry[0], t0); err != nil {
		t.Fatalf("initial OpenOrUpdate: %v", err)
	}

	// Now the device enters a maintenance window; a new breach must NOT refresh.
	t1 := t0.Add(30 * time.Minute)
	if _, err := schedule.AddWindow(state, deviceID, t1.Add(-5*time.Minute), t1.Add(time.Hour)); err != nil {
		t.Fatalf("AddWindow: %v", err)
	}
	second := &store.Sample{PointID: pointID, PointName: pointName, DeviceID: deviceID, Value: 300, At: t1}

	a, err := OpenOrUpdate(state, second, entry[0], t1)
	if err != nil {
		t.Fatalf("OpenOrUpdate: %v", err)
	}
	if a != nil {
		t.Fatalf("expected refresh inside window to be suppressed (nil), got %v", a)
	}
	// The original alarm must be untouched: its last trigger stays at t0.
	latest, ok := state.AlarmByPoint(pointID)
	if !ok {
		t.Fatal("expected the original alarm to still exist")
	}
	if !latest.LastTriggerAt.Equal(t0) {
		t.Fatalf("expected original alarm LastTriggerAt unchanged at %v, got %v", t0, latest.LastTriggerAt)
	}
}

func mustLoadSnapshot(t *testing.T, s *store.State) []store.RuleSnapshotEntry {
	t.Helper()
	entries, err := s.LoadSnapshot(s.Effective())
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one snapshot entry")
	}
	return entries
}
