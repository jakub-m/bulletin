package time

import (
	"fmt"
	gotime "time"
)

var timeFormats = []string{
	gotime.RFC1123,
	gotime.RFC1123Z,
	gotime.RFC3339,
	`2006-01-02T15:04:05.000-07:00`,
}

func Parse(value string) (gotime.Time, error) {
	for _, f := range timeFormats {
		t, err := gotime.Parse(f, value)
		if err == nil {
			return t, nil
		}
	}
	return gotime.Time{}, fmt.Errorf("cannot parse time: `%s`", value)
}
