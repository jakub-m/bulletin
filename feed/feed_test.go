package feed_test

import (
	"bulletin/atom"
	"bulletin/feed"
	"bulletin/testutils"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestGoogleBlogArticles(t *testing.T) {
	articles := parseArticlesFromXml(t, "../testdata/atom_google_ai_blog.xml")
	expectedArticles := []feed.Article{
		{
			Feed: feed.Feed{
				Id:    "tag:blogger.com,1999:blog-8474926331452026626",
				Title: "Google AI Blog",
				Url:   "http://ai.googleblog.com/",
			},
			Id:      "tag:blogger.com,1999:blog-8474926331452026626.post-537064785672594983",
			Title:   "Mapping Africa’s Buildings with Satellite Imagery",
			Url:     "http://ai.googleblog.com/2021/07/mapping-africas-buildings-with.html",
			Updated: testutils.ParseTime(t, "2021-07-29T13:05:10.956-07:00"),
		},
	}
EXPECTED:
	for _, expected := range expectedArticles {
		for _, actual := range articles {
			if fmt.Sprintf("%+v", actual) == fmt.Sprintf("%+v", expected) {
				continue EXPECTED
			}
		}
		t.Errorf("Article not found: %+v", expected)
	}
	if t.Failed() {
		t.Log("Actual articles:")
		for i, a := range articles {
			t.Logf("\t%d\t%+v", i, a)
		}
	}
}

func parseArticlesFromXml(t *testing.T, path string) []feed.Article {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	parsed, err := atom.Parse(b)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	return parsed.GetArticles()
}
