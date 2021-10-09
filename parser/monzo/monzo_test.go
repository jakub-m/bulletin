package monzo

import (
	"bulletin/feed"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonzoParser(t *testing.T) {
	f := parseFromFile(t, &monzoFeedParser{}, "testdata/monzo_com_blog_technology.html")
	assert.Len(t, f.Articles, 13)
	firstArticle := f.Articles[0]

	assert.Equal(t, firstArticle.Id[:20], "Documenting pull req")
	assert.Equal(t, firstArticle.Title[:20], "Documenting pull req")
	assert.Equal(t, firstArticle.Description[:20], "How our engineering ")
	assert.Equal(t, "2021-09-30 00:00:00 +0000 UTC", fmt.Sprint(firstArticle.Published))
	assert.Equal(t, firstArticle.Url, "/blog/2021/09/30/documenting-pull-requests-is-as-important-as-writing-good-code")
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
