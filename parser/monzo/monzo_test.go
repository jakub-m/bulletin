package monzo

import (
	"bulletin/feed"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonzoParser(t *testing.T) {
	f := parseFromFile(t, &FeedParser{}, "testdata/monzo_com_blog_technology.html")
	assert.Len(t, f.Articles, 3)
}

func parseFromFile(t *testing.T, parser feed.FeedParser, path string) feed.Feed {
	file, err := os.Open(path)
	assert.NoError(t, err)
	defer file.Close()
	body, err := io.ReadAll(file)
	assert.NoError(t, err)
	fe, err := parser.ParseFeed(body)
	assert.NoError(t, err)
	return fe
}
