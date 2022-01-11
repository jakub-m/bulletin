package time

import (
	"testing"
	gotime "time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	tcs := []string{
		"Thu, 26 Aug 2021 12:22:53 +0200",
		"Fri, 25 Jun 2021 11:03:04 GMT",
		"2021-06-25T11:03:03.265Z",
		"2021-08-23T03:19:41.679-04:00",
		"2021-01-23T00:00",
		"2021-02-01",
	}
	for _, tc := range tcs {
		tim, err := Parse(tc)
		assert.NoError(t, err)
		assert.True(t, tim.After(gotime.Time{}), "wrong time for `%s`", tc)
	}
}
