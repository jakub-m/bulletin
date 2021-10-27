package aykevl

import (
	"bulletin/feed"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	f := parseFromFile(t, &aykevlFeedParser{}, "testdata/aykevl_nl.html", blogUrl)

	assert.Equal(t, "https://aykevl.nl/", f.Id)
	assert.Equal(t, "Ayke van Laëthem", f.Title)
	assert.Equal(t, blogUrl, f.Url)

	assert.Len(t, f.Articles, 10)
	firstArticle := f.Articles[0]

	assert.Equal(t, firstArticle.Id[:20], "What's the int type?")
	assert.Equal(t, firstArticle.Title[:20], "What's the int type?")
	assert.Equal(t, firstArticle.Description[:20], "The int type is pres")
	// assert.Equal(t, "2021-06-25 00:00:00 +0000 UTC", fmt.Sprint(firstArticle.Published))
	assert.Equal(t, firstArticle.Url, "/2021/06/what-is-int")
}

func parseFromFile(t *testing.T, parser feed.FeedParser, filePath, url string) feed.Feed {
	file, err := os.Open(filePath)
	assert.NoError(t, err)
	defer file.Close()
	body, err := io.ReadAll(file)
	assert.NoError(t, err)
	fe, err := parser.ParseFeed(body, url)
	assert.NoError(t, err)
	return fe
}
