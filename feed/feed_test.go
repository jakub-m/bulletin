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

	"github.com/stretchr/testify/assert"
)

func TestGoogleBlogArticles(t *testing.T) {
	articles := parseArticlesFromXml(t, "../testdata/atom_google_ai_blog.xml")
	expected := feed.Article{
		Feed: feed.Feed{
			Id:    "tag:blogger.com,1999:blog-8474926331452026626",
			Title: "Google AI Blog",
			Url:   "http://ai.googleblog.com/",
		},
		Id:      "tag:blogger.com,1999:blog-8474926331452026626.post-537064785672594983",
		Title:   "Mapping Africa’s Buildings with Satellite Imagery",
		Url:     "http://ai.googleblog.com/2021/07/mapping-africas-buildings-with.html",
		Updated: testutils.ParseTime(t, "2021-07-29T13:05:10.956-07:00"),
	}
	assert.Equal(t, asJson(t, expected), asJson(t, articles[0]))
}

func TestNetflixArticles(t *testing.T) {
	articles := parseArticlesFromXml(t, "../testdata/rss_netflix_techblog.xml")
	expected := feed.Article{
		Feed: feed.Feed{
			Id:    "https://netflixtechblog.com/feed",
			Title: "Netflix TechBlog - Medium",
			Url:   "https://netflixtechblog.com/feed",
		},
		Id:      "https://medium.com/p/3fddcceb1059",
		Title:   "Data Movement in Netflix Studio via Data Mesh",
		Url:     "https://netflixtechblog.com/data-movement-in-netflix-studio-via-data-mesh-3fddcceb1059?source=rss----2615bd06b42e---4",
		Updated: testutils.ParseTime(t, "2021-07-26T18:00:56+00:00"),
	}
	assert.Equal(t, asJson(t, expected), asJson(t, articles[0]))
}

func asJson(t *testing.T, a feed.Article) string {
	t.Helper()
	j, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	return string(j)
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
