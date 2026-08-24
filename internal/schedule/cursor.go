package schedule

import (
	"errors"

	"telemetryguard/internal/store"
)

// ErrBatchNotCommitted is returned when advancing past an uncommitted batch.
var ErrBatchNotCommitted = errors.New("batch is not committed")

// AdvanceCursor moves the evaluation cursor forward only after the batch was
// durably committed, so a failed batch never gets skipped silently.
func AdvanceCursor(state *store.State, batch *store.EvalBatch) error {
	return state.SetCursor(batch.Cursor)
}
