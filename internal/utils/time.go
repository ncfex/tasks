package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mergestat/timediff"
)

func FormatTimeToHuman(t time.Time) string {
	return timediff.TimeDiff(t)
}

func ParseHumanToTime(humanTime string) (time.Time, error) {
	humanTime = strings.ToLower(strings.TrimSpace(humanTime))
	now := time.Now()

	if humanTime == "tomorrow" {
		return now.AddDate(0, 0, 1), nil
	}

	re := regexp.MustCompile(`^(\d+)\s+(second|minute|hour|day|week|month|year)s?\s+ago$`)
	matches := re.FindStringSubmatch(humanTime)

	if len(matches) == 3 {
		quantity, err := strconv.Atoi(matches[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid quantity: %s", matches[1])
		}

		unit := matches[2]
		switch unit {
		case "second":
			return now.Add(time.Duration(-quantity) * time.Second), nil
		case "minute":
			return now.Add(time.Duration(-quantity) * time.Minute), nil
		case "hour":
			return now.Add(time.Duration(-quantity) * time.Hour), nil
		case "day":
			return now.AddDate(0, 0, -quantity), nil
		case "week":
			return now.AddDate(0, 0, -quantity*7), nil
		case "month":
			return now.AddDate(0, -quantity, 0), nil
		case "year":
			return now.AddDate(-quantity, 0, 0), nil
		}
	}

	re = regexp.MustCompile(`^in\s+(\d+)\s+(second|minute|hour|day|week|month|year)s?$`)
	matches = re.FindStringSubmatch(humanTime)

	if len(matches) == 3 {
		quantity, err := strconv.Atoi(matches[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid quantity: %s", matches[1])
		}

		unit := matches[2]
		switch unit {
		case "second":
			return now.Add(time.Duration(quantity) * time.Second), nil
		case "minute":
			return now.Add(time.Duration(quantity) * time.Minute), nil
		case "hour":
			return now.Add(time.Duration(quantity) * time.Hour), nil
		case "day":
			return now.AddDate(0, 0, quantity), nil
		case "week":
			return now.AddDate(0, 0, quantity*7), nil
		case "month":
			return now.AddDate(0, quantity, 0), nil
		case "year":
			return now.AddDate(quantity, 0, 0), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time string: %s", humanTime)
}
