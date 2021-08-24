package testutils

import (
	"testing"
	"time"
)

func ParseTime(t *testing.T, value string) time.Time {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("parseTime: %s", err)
	}
	return parsed
}
