package hakibenita

import (
	"bulletin/testutils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHakibenitaParse(t *testing.T) {
	f := testutils.ParseFromFile(t, &hakibenitaParser{}, "testdata/hakibenita_com.html", blogUrl)

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
