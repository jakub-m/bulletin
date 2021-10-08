package notion

import (
	"bulletin/feed"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNotion(t *testing.T) {
	f := parseFromFile(t, "testdata/notion_blog_topic_tech.html")
	assert.Len(t, f.Articles, 3)
	// more tests to come
}

func parseFromFile(t *testing.T, path string) feed.Feed {
	parser := FeedParser{}

	file, err := os.Open(path)
	assert.NoError(t, err)
	defer file.Close()
	body, err := io.ReadAll(file)
	assert.NoError(t, err)
	fe, err := parser.ParseFeed(body)
	assert.NoError(t, err)
	return fe
}
