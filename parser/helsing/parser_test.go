package helsing

import (
	"bulletin/testutils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelsingParser(t *testing.T) {
	f := testutils.ParseFromFile(t, &helsingFeedParser{}, "testdata/blog_helsing_ai.html", helsingBlogUrl)

	assert.Equal(t, "Helsing - Helsing Blog", f.Id)
	assert.Equal(t, "Helsing - Helsing Blog", f.Title)
	assert.Equal(t, helsingBlogUrl, f.Url)

	assert.Len(t, f.Articles, 26)

	firstArticle := f.Articles[0]
	assert.Equal(t, "AI-assisted vetting of software packages", firstArticle.Id)
	assert.Equal(t, "AI-assisted vetting of software packages", firstArticle.Title)
	assert.Equal(t, firstArticle.Description[:20], "Describes Helsing's ")
	assert.Equal(t, "2025-09-10 00:00:00 +0000 UTC", fmt.Sprint(firstArticle.Published))
	assert.Equal(t, "https://blog.helsing.ai/posts/ai-assisted-vetting-of-software-packages/", firstArticle.Url)
}
