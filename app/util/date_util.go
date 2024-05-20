package util

import (
	"fmt"
	"time"
)

// ParseTime parses a string into a time.Time object
func ParseTime(timeStr string) (time.Time, error) {

	timeFormats := []string{time.DateTime, time.RFC3339Nano, time.RFC3339}

	for _, format := range timeFormats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time format %s", timeStr)
}
