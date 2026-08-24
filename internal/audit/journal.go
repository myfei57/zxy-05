package audit

import (
	"path/filepath"

	"telemetryguard/internal/store"
)

// Record durably appends one audit entry and maintains the persisted active
// alarm count for open and close transitions.
func Record(state *store.State, entry *store.AuditEntry) error {
	if err := state.AppendAudit(entry); err != nil {
		return err
	}
	count := persistedCount(state)
	switch entry.Action {
	case "alarm_open":
		count++
	case "alarm_close":
		if count > 0 {
			count--
		}
	}
	return store.SaveJSON(filepath.Join(state.Root(), "count.json"), count)
}

func persistedCount(state *store.State) int {
	var count int
	if err := store.LoadJSON(filepath.Join(state.Root(), "count.json"), &count); err != nil {
		return 0
	}
	return count
}
