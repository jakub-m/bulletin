package monzo

import (
	"bulletin/testutils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonzoParser(t *testing.T) {
	f := testutils.ParseFromFile(t, &monzoFeedParser{}, "testdata/monzo_com_blog_technology.html", monzoBlogUrl)

	assert.Equal(t, "Monzo - Technology", f.Id)
	assert.Equal(t, "Monzo - Technology", f.Title)
	assert.Equal(t, monzoBlogUrl, f.Url)

	assert.Len(t, f.Articles, 13)
	firstArticle := f.Articles[0]

	assert.Equal(t, firstArticle.Id[:20], "Documenting pull req")
	assert.Equal(t, firstArticle.Title[:20], "Documenting pull req")
	assert.Equal(t, firstArticle.Description[:20], "How our engineering ")
	assert.Equal(t, "2021-09-30 00:00:00 +0000 UTC", fmt.Sprint(firstArticle.Published))
	assert.Equal(t, firstArticle.Url, "/blog/2021/09/30/documenting-pull-requests-is-as-important-as-writing-good-code")
}
