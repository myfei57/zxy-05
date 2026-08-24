package console

import (
	"errors"
	"time"
)

type notFoundError struct{}

func (*notFoundError) Error() string { return "not found" }

var errNotFound = &notFoundError{}

var errBadTime = errors.New("invalid RFC3339 time")

var errDuplicatePoint = errors.New("point name already exists")

func parseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}
