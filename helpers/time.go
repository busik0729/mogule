package helpers

import (
	"time"
)

func GetDateWithZeroHour(t int64) time.Time {
	a := time.Unix(t, 0)
	y, m, d := a.Date()
	l := a.Location()
	return time.Date(y, m, d, 0, 0, 0, 0, l)
}
