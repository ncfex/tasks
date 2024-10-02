package utils

import (
	"time"

	"github.com/mergestat/timediff"
)

func HumanReadableTime(time time.Time) string {
	return timediff.TimeDiff(time)
}
