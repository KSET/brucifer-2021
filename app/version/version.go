package version

import (
	"fmt"
	"time"
)

const TimeFormat = "2006-01-02T15:04:05-0700"

var (
	CommitHash                = "dev"
	buildTimestamp            = TimeFormat
	initialTime               = time.Now()
	parsedTime     *time.Time = nil
)

func BuildTime() *time.Time {
	if parsedTime != nil {
		return parsedTime
	}

	if buildTimestamp == TimeFormat {
		parsedTime = &initialTime
		return parsedTime
	}

	t, err := time.Parse(TimeFormat, buildTimestamp)

	if err == nil {
		localTime := t.Local()
		parsedTime = &localTime
	} else {
		parsedTime = &initialTime
	}

	return parsedTime
}

func BuildVersion() string {
	return fmt.Sprintf("%s (%s)", CommitHash, BuildTime().Format(TimeFormat))
}
