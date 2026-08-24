package rule

import (
	"errors"

	"telemetryguard/internal/store"
)

// ErrUnknownVersion is returned when the target version has no snapshot.
var ErrUnknownVersion = errors.New("rule version has no snapshot")

// RollbackTo restores an earlier published version as the effective rule set.
// Only the rules of the target snapshot become effective; draft rules are
// never merged into the evaluation set.
func RollbackTo(state *store.State, target int64) error {
	entries, err := state.LoadSnapshot(target)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return ErrUnknownVersion
	}
	return state.SetEffective(target)
}
