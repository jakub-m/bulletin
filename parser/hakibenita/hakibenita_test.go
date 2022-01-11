package hakibenita

import (
	"bulletin/feed"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHakibenitaParse(t *testing.T) {
	f := parseFromFile(t, &hakibenitaParser{}, "testdata/hakibenita_com.html", blogUrl)

	assert.Equal(t, "Haki Benita", f.Id)
	assert.Equal(t, "Haki Benita", f.Title)
	assert.Equal(t, blogUrl, f.Url)

	assert.Len(t, f.Articles, 11)
	firstArticle := f.Articles[0]

	assert.Equal(t, firstArticle.Id, "2021 Year in Review")
	assert.Equal(t, firstArticle.Title, "2021 Year in Review")
	assert.Equal(t, firstArticle.Description[:20], "What I've been up to")
	assert.Equal(t, "2021-12-31 00:00:00 +0000 UTC", fmt.Sprint(firstArticle.Published))
	assert.Equal(t, firstArticle.Url, "https://hakibenita.com/2021-year-in-review")
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
