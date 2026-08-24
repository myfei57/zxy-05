package verifycase

import (
	"testing"
	"time"

	"telemetryguard/internal/eval"
	"telemetryguard/internal/rule"
	"telemetryguard/internal/store"
)

func TestRollbackExcludesDraftRules(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	if _, err := rule.CreateRule(state, "rule-a", "temp_main", ">", 80, 1); err != nil {
		t.Fatal(err)
	}
	v1, err := rule.PublishVersion(state)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := rule.CreateRule(state, "rule-b", "vib_main", ">", 5, 1); err != nil {
		t.Fatal(err)
	}
	if err := rule.RollbackTo(state, v1); err != nil {
		t.Fatal(err)
	}
	sample := &store.Sample{PointID: "p2", PointName: "vib_main", DeviceID: "d1", Seq: 1, Value: 9, At: time.Now().UTC()}
	opened, err := eval.EvaluateSamples(state, []*store.Sample{sample}, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if len(opened) != 0 {
		t.Fatalf("%d alarms opened from a draft rule after rollback", len(opened))
	}
}
