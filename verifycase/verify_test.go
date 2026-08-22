package verifycase

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"telemetryguard/internal/ingest"
	"telemetryguard/internal/cycle"
	"telemetryguard/internal/store"
)

func TestEvalCursorNotAdvancedBeforeBatchDurable(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	if err := os.WriteFile(filepath.Join(dir, "batches"), []byte("blocker"), 0o644); err != nil {
		t.Fatal(err)
	}
	ingest.AppendPending(&store.Sample{PointID: "p1", PointName: "temp_main", DeviceID: "d1", Seq: 1, Value: 95, At: time.Now().UTC()})
	_, err := cycle.RunCycle(state, time.Now().UTC())
	if err == nil {
		t.Fatal("expected the batch commit to fail")
	}
	if state.Cursor() != 0 {
		t.Fatalf("cursor advanced to %d although the batch was not committed", state.Cursor())
	}
}
