package rss

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseSchema(t *testing.T) {
	actual := parseRssFromFile(t, "../../testdata/rss_schema.xml")
	expected := &Channel{
		Title:       "Title",
		Description: "Description",
		Items: []Item{
			{
				Title:          "Item Title",
				Description:    "Item description",
				Link:           "http://example.com/item",
				Guid:           "http://example.com/guid/0123",
				PubDate:        parseTime(t, "2021-07-26T18:00:56Z"),
				ContentEncoded: "<p>Content</p>",
			},
		},
	}
	assert.Equal(t, channelAsJson(t, expected), channelAsJson(t, actual))
}

func TestParseNetflix(t *testing.T) {
	channel := parseRssFromFile(t, "../../testdata/rss_netflix_techblog.xml")
	assert.Equal(t, "Netflix TechBlog - Medium", channel.Title)
	assert.Equal(t, 113, len(channel.Description))
	assert.Equal(t, 10, len(channel.Items))
	item := channel.Items[0]
	assert.Equal(t, "Data Movement in Netflix Studio via Data Mesh", item.Title)
	assert.Equal(t, "https://medium.com/p/3fddcceb1059", item.Guid)
	assert.Equal(t, "2021-07-26T18:00:56Z", fmt.Sprint(item.PubDate))
	assert.Equal(t, 117, len(item.Link))
	assert.Equal(t, 22294, len(item.ContentEncoded))
}

func TestParseBerthub(t *testing.T) {
	channel := parseRssFromFile(t, "../../testdata/rss_berthub.xml")
	assert.Greater(t, len(channel.Items), 0)
}

func parseRssFromFile(t *testing.T, path string) *Channel {
	f, err := os.Open(path)
	assert.NoError(t, err)
	b, err := io.ReadAll(f)
	assert.NoError(t, err)
	feed, err := Parse(b)
	assert.NoError(t, err)
	return feed
}

func parseTime(t *testing.T, value string) *RssTime {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("parseTime: %s", err)
	}
	return &RssTime{parsed}
}

// channelAsJson is used because without it time.Time does not compare well when expressed with different timezones.
func channelAsJson(t *testing.T, channel *Channel) string {
	s, err := json.MarshalIndent(channel, "", " ")
	if err != nil {
		t.Fatalf("channelAsJson: %s", err)
	}
	return string(s)
}
