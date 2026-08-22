package rule

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"telemetryguard/internal/store"
)

// ErrBadOperator is returned for an unsupported comparison operator.
var ErrBadOperator = errors.New("unsupported rule operator")

// CreateRule creates a draft rule.
func CreateRule(state *store.State, name, pointName, op string, threshold float64, level int) (*store.Rule, error) {
	name = strings.TrimSpace(name)
	pointName = strings.TrimSpace(pointName)
	if name == "" || pointName == "" {
		return nil, errors.New("rule name and point name are required")
	}
	if op != ">" && op != ">=" && op != "<" && op != "<=" {
		return nil, ErrBadOperator
	}
	if level < 1 || level > 3 {
		return nil, errors.New("rule level must be 1..3")
	}
	r := &store.Rule{
		ID:        uuid.NewString(),
		Name:      name,
		PointName: pointName,
		Op:        op,
		Threshold: threshold,
		Level:     level,
		State:     store.RuleDraft,
		CreatedAt: time.Now().UTC(),
	}
	if err := state.PutRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

// ListRules returns all rules ordered by name.
func ListRules(state *store.State) []*store.Rule {
	return state.Rules()
}
