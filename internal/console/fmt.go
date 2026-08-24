package console

import (
	"strconv"
)

func formatInt(v int64) string {
	return strconv.FormatInt(v, 10)
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 3, 64)
}
