package audit

import (
	"path/filepath"

	"telemetryguard/internal/store"
)

// CountActive returns the persisted active alarm count, falling back to the
// live registry when no count has been recorded yet.
func CountActive(state *store.State) int {
	var count int
	if err := store.LoadJSON(filepath.Join(state.Root(), "count.json"), &count); err == nil && count > 0 {
		return count
	}
	return len(state.ActiveAlarms())
}
