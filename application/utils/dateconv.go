package utils

import (
	"time"
)

func ConvSecToDateString(seconds int) string {
	return time.Unix(int64(seconds), 0).Format("2006-01-02T15:04:05Z")
}

func GetCurrentTimeString() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func GetTimeNow() *time.Time {
	timeNow := time.Now().UTC()
	return &timeNow
}
