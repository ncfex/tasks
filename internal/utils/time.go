package utils

import (
	"time"

	"github.com/mergestat/timediff"
)

func FormatTimeToHuman(t time.Time) string {
	return timediff.TimeDiff(t)
}
