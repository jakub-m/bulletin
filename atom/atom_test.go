package atom

import (
	"io"
	"os"
	"sort"
	"testing"

	"bulletin/testutils"

	"github.com/stretchr/testify/assert"
)

func TestParserGoogleBlog(t *testing.T) {
	feed := parseAtomFromFile(t, "../testdata/atom_google_ai_blog.xml")

	assert.Equal(t, "tag:blogger.com,1999:blog-8474926331452026626", feed.Id)
	assert.Equal(t, "Google AI Blog", feed.Title)
	assert.Equal(t, "The latest news from Google AI.", feed.Subtitle)
	assert.Equal(t, len(feed.Entries), 25)

	var rels []string
	for _, l := range feed.Links {
		rels = append(rels, l.Rel)
	}
	sort.Strings(rels)
	assert.Equal(t, []string{"alternate", "hub", "next", "self"}, rels)

	entry := feed.Entries[0]
	assert.Equal(t, "tag:blogger.com,1999:blog-8474926331452026626.post-537064785672594983", entry.Id)
	assert.Equal(t, "Mapping Africa’s Buildings with Satellite Imagery", entry.Title)
	assert.Equal(t, "http://ai.googleblog.com/feeds/537064785672594983/comments/default", entry.Links[0].Href)
}

func TestParseAtomSchema(t *testing.T) {
	feed := parseAtomFromFile(t, "../testdata/atom_schema.xml")

	expected := &Feed{
		Id:       "Id",
		Title:    "Title",
		Subtitle: "Subtitle",
		Entries: []Entry{
			{
				Id:        "Id",
				Title:     "Title",
				Published: parseTime(t, "2000-01-01T00:00:00.000+01:00"),
				Updated:   parseTime(t, "2000-01-02T00:00:00.000+01:00"),
				Links:     nil,
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

func parseAtomFromFile(t *testing.T, path string) *Feed {
	f, err := os.Open(path)
	assert.NoError(t, err)
	b, err := io.ReadAll(f)
	assert.NoError(t, err)
	feed, err := Parse(b)
	assert.NoError(t, err)
	return feed
}

func parseTime(t *testing.T, value string) *XmlTime {
	parsed := testutils.ParseTime(t, value)
	return &XmlTime{parsed}
}
