package escalate

import (
	"testing"
	"time"

	"telemetryguard/internal/alarm"
	"telemetryguard/internal/store"
)

// newChannel seeds a single reachable channel so notify.Dispatch has a route.
func newChannel(t *testing.T, state *store.State) string {
	t.Helper()
	c := &store.ChannelConfig{
		ID:        "ch-1",
		Name:      "sms",
		Enabled:   true,
		Reachable: true,
		Priority:  1,
		CreatedAt: time.Date(2026, 8, 22, 0, 0, 0, 0, time.UTC),
	}
	if err := state.PutChannel(c); err != nil {
		t.Fatalf("PutChannel: %v", err)
	}
	return c.ID
}

// openAlarmWithTimer creates an open alarm and a pending escalation timer whose
// deadline is already in the past, so the next Tick is due to fire.
func openAlarmWithTimer(t *testing.T, state *store.State, openAt time.Time) *store.Alarm {
	t.Helper()
	a := &store.Alarm{
		ID:             "alarm-1",
		PointID:        "point-1",
		Status:         store.AlarmOpen,
		OriginalOpenAt: openAt,
		LastTriggerAt:  openAt,
	}
	if err := state.PutAlarm(a); err != nil {
		t.Fatalf("PutAlarm: %v", err)
	}
	if err := EnsureTimer(state, a, openAt); err != nil {
		t.Fatalf("EnsureTimer: %v", err)
	}
	return a
}

// TestTickOpenAlarmEscalates is the happy path: a still-open alarm past its SLA
// deadline escalates and dispatches exactly one notification.
func TestTickOpenAlarmEscalates(t *testing.T) {
	state := store.NewState(t.TempDir())
	newChannel(t, state)

	openAt := time.Date(2026, 8, 22, 3, 0, 0, 0, time.UTC)
	openAlarmWithTimer(t, state, openAt)
	now := openAt.Add(31 * time.Minute)

	dispatched, err := Tick(state, now)
	if err != nil {
		t.Fatalf("Tick: %v", err)
	}
	if len(dispatched) != 1 {
		t.Fatalf("open alarm past deadline should escalate once, got %d notifications", len(dispatched))
	}
	e, ok := state.Escalation("alarm-1")
	if !ok {
		t.Fatalf("escalation record missing")
	}
	if e.Level != 1 {
		t.Fatalf("level should advance to 1, got %d", e.Level)
	}
	if e.State != store.EscalationPending {
		t.Fatalf("open alarm timer should stay pending, got %q", e.State)
	}
}

// TestTickAcknowledgedAlarmDoesNotEscalate reproduces the 3 AM false escalation:
// once an operator acknowledges the alarm, later ticks must not dispatch against
// the stale open-time snapshot. The timer is cancelled instead.
func TestTickAcknowledgedAlarmDoesNotEscalate(t *testing.T) {
	state := store.NewState(t.TempDir())
	newChannel(t, state)

	openAt := time.Date(2026, 8, 22, 3, 0, 0, 0, time.UTC)
	a := openAlarmWithTimer(t, state, openAt)

	// First tick escalates the still-open alarm to level 1.
	now := openAt.Add(31 * time.Minute)
	if _, err := Tick(state, now); err != nil {
		t.Fatalf("Tick: %v", err)
	}

	// An operator acknowledges the alarm; it no longer needs escalation.
	if err := alarm.Acknowledge(state, a.ID, "oncall", "handled", now.Add(time.Minute)); err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}

	// A later tick must not escalate further: the alarm's live state is
	// "acknowledged", so the timer is cancelled rather than firing from the
	// stale open-time snapshot.
	dispatched, err := Tick(state, now.Add(2*time.Minute))
	if err != nil {
		t.Fatalf("Tick: %v", err)
	}
	if len(dispatched) != 0 {
		t.Fatalf("acknowledged alarm must not escalate, got %d notifications", len(dispatched))
	}
	e, ok := state.Escalation(a.ID)
	if !ok {
		t.Fatalf("escalation record missing")
	}
	if e.State != store.EscalationDone {
		t.Fatalf("acknowledged alarm timer should be cancelled (done), got %q", e.State)
	}
}

// TestTickResolvedAlarmDoesNotEscalate covers the recovered case: a resolved
// alarm must not keep escalating either.
func TestTickResolvedAlarmDoesNotEscalate(t *testing.T) {
	state := store.NewState(t.TempDir())
	newChannel(t, state)

	openAt := time.Date(2026, 8, 22, 3, 0, 0, 0, time.UTC)
	a := openAlarmWithTimer(t, state, openAt)

	now := openAt.Add(31 * time.Minute)
	if _, err := Tick(state, now); err != nil {
		t.Fatalf("Tick: %v", err)
	}

	if err := alarm.Resolve(state, a.ID, now.Add(time.Minute)); err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	dispatched, err := Tick(state, now.Add(2*time.Minute))
	if err != nil {
		t.Fatalf("Tick: %v", err)
	}
	if len(dispatched) != 0 {
		t.Fatalf("resolved alarm must not escalate, got %d notifications", len(dispatched))
	}
	e, ok := state.Escalation(a.ID)
	if !ok {
		t.Fatalf("escalation record missing")
	}
	if e.State != store.EscalationDone {
		t.Fatalf("resolved alarm timer should be cancelled (done), got %q", e.State)
	}
}
