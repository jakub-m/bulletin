package feed

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.Contains(t, body, `<a href="http://example.com/1">title-1</a>`)
	assert.Contains(t, body, `<a href="http://example.com/1">title-2</a>`)
}
