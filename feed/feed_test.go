package feed_test

import (
	"bulletin/atom"
	"bulletin/feed"
	"bulletin/rss"
	"bulletin/testutils"
	"encoding/json"
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
	assertArticlesContain(t, expectedArticles, articles)
}

func TestNetflixArticles(t *testing.T) {
	articles := parseArticlesFromXml(t, "../testdata/rss_netflix_techblog.xml")
	expectedArticles := []feed.Article{
		{
			Feed: feed.Feed{
				Id:    "Netflix TechBlog - Medium",
				Title: "Netflix TechBlog - Medium",
				Url:   "",
			},
			Id:      "https://medium.com/p/3fddcceb1059",
			Title:   "Data Movement in Netflix Studio via Data Mesh",
			Url:     "https://netflixtechblog.com/data-movement-in-netflix-studio-via-data-mesh-3fddcceb1059?source=rss----2615bd06b42e---4",
			Updated: testutils.ParseTime(t, "2021-07-26T18:00:56+00:00"),
		},
	}
	assertArticlesContain(t, expectedArticles, articles)
}

func assertArticlesContain(t *testing.T, expected []feed.Article, actual []feed.Article) {
EXPECTED:
	for _, e := range expected {
		expectedJson, err := json.Marshal(e)
		expectedJsonString := string(expectedJson)
		if err != nil {
			t.Fatal(err)
		}
		for _, a := range actual {
			actualJson, err := json.Marshal(a)
			actualJsonString := string(actualJson)
			if err != nil {
				t.Fatal(err)
			}
			// if fmt.Sprintf("%+v", a) == fmt.Sprintf("%+v", e) {
			// 	continue EXPECTED
			// }
			if actualJsonString == expectedJsonString {
				continue EXPECTED
			}
		}
		t.Errorf("Article not found: %s", expectedJsonString)
	}
	if t.Failed() {
		t.Log("Actual articles:")
		for i, a := range actual {
			actualJson, _ := json.Marshal(a)
			t.Logf("\t%d\t%s", i, string(actualJson))
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
	atomFeed, atomErr := atom.Parse(b)
	if atomErr == nil && len(atomFeed.Entries) > 0 {
		return atomFeed.GetArticles()
	}
	rssFeed, rssErr := rss.Parse(b)
	if rssErr == nil && len(rssFeed.Items) > 0 {
		return rssFeed.GetArticles()
	}
	t.Fatalf("No articles. Atom err: %s, RSS err: %s", atomErr, rssErr)
	return []feed.Article{}
}
