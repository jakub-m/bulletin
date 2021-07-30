package feed

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestFormatHtml(t *testing.T) {
	feeds := []Article{
		{
			Id:      "id-1",
			Title:   "title-1",
			Url:     "http://example.com/1",
			Updated: time.Time{},
		},
		{
			Id:      "id-2",
			Title:   "title-2",
			Url:     "http://example.com/1",
			Updated: time.Time{}.Add(1 * time.Minute),
		},
	}
	body, err := FormatHtml(feeds)
	assert.NoError(t, err)
	expected := `
<a href="http://example.com/1">title-1</a>
<a href="http://example.com/1">title-2</a>
`

	assert.Equal(t, strings.Trim(expected, "\n"), body)
}
