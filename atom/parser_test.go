package atom

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"time"
)

func TestParserGoogleBlog(t *testing.T) {
	f, err := os.Open("../testdata/google_ai_blog.xml")
	assert.NoError(t, err)
	b, err := io.ReadAll(f)
	assert.NoError(t, err)
	feed, err := Parse(b)
	assert.NoError(t, err)

	assert.Equal(t, feed.Id, "tag:blogger.com,1999:blog-8474926331452026626")
	assert.Equal(t, feed.Title, "Google AI Blog")
	assert.Equal(t, feed.Subtitle, "The latest news from Google AI.")

	assert.Equal(t, len(feed.Entries), 25)
	entry := feed.Entries[0]

	assert.Equal(t, entry.Id, "tag:blogger.com,1999:blog-8474926331452026626.post-537064785672594983")
	assert.Equal(t, entry.Title, "Mapping Africa’s Buildings with Satellite Imagery")
}

func TestParser(t *testing.T) {
	f, err := os.Open("../testdata/schema_full.xml")
	assert.NoError(t, err)
	b, err := io.ReadAll(f)
	assert.NoError(t, err)
	feed, err := Parse(b)
	assert.NoError(t, err)

	expected := &Feed{
		Id:       "Id",
		Title:    "Title",
		Subtitle: "Subtitle",
		Entries: []*Entry{
			{
				Id:        "Id",
				Title:     "Title",
				Published: parseTime(t, "2000-01-01T00:00:00.000+01:00"),
				Updated:   parseTime(t, "2000-01-02T00:00:00.000+01:00"),
			},
		},
	}
	assert.Equal(t, expected, feed)
}

func TestParseXmlTime(t *testing.T) {
	x := &XmlTime{}
	tcs := []struct {
		time     string
		expected string
	}{
		{
			"2006-01-02T15:04:05.000-07:00",
			"2006-01-02T15:04:05-07:00",
		},
		{
			"2021-07-27T09:49:00.001-07:00",
			"2021-07-27T09:49:00-07:00",
		},
	}
	for _, tc := range tcs {
		e := x.parse(tc.time)
		assert.NoError(t, e)
		assert.Equal(t, tc.expected, x.String())
	}
}

func TestEntryUid(t *testing.T) {
	e := Entry {
		Id: "some:id",
		Updated: parseTime(t, "2000-01-02T00:00:00.000+01:00"),
	}
	assert.Equal(t, "2000-01-01T23:00:00Z-some:id", e.Uid())
}

func parseTime(t *testing.T, value string) *XmlTime {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Errorf("parseTime: %s", err)
	}
	return &XmlTime{parsed}
}
