package benoitessiambre

import (
	"bulletin/testutils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	f := testutils.ParseFromFile(t, FeedParser, "testdata/benoitessiambre_com.html", blogUrl)

	assert.Equal(t, "https://benoitessiambre.com/blemish.html", f.Id)
	assert.Equal(t, "benoitessiambre.com", f.Title)
	assert.Equal(t, blogUrl, f.Url)

	assert.Len(t, f.Articles, 16)
	firstArticle := f.Articles[0]

	assert.Equal(t, firstArticle.Id, "/macro.html")
	assert.Equal(t, firstArticle.Title, "Sim Central Bank")
	assert.Equal(t, firstArticle.Description, "game")
	assert.Equal(t, "2022-06-15 00:00:00 +0000 UTC", fmt.Sprint(firstArticle.Published))
	assert.Equal(t, firstArticle.Url, "/macro.html")
}
