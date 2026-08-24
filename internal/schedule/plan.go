package schedule

import (
	"time"

	"github.com/google/uuid"

	"telemetryguard/internal/store"
)

// NewBatch creates a pending evaluation batch for the next cursor position.
func NewBatch(state *store.State, pointCount int) (*store.EvalBatch, error) {
	b := &store.EvalBatch{
		ID:         uuid.NewString(),
		Cursor:     state.Cursor() + 1,
		State:      store.BatchPending,
		PointCount: pointCount,
		StartedAt:  time.Now().UTC(),
	}
	return b, nil
}
