package utils

import "time"

func HumanReadableTime(time time.Time) string {
	return time.GoString()
}
