package feed_test

import (
	"bulletin/feed"
	"bulletin/parser/atom"
	"bulletin/parser/rss"
	"bulletin/testutils"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoogleBlogArticles(t *testing.T) {
	actual := getFeed(t, "../parser/atom/testdata/atom_google_ai_blog.xml")
	expected := feed.Feed{
		Id:    "tag:blogger.com,1999:blog-8474926331452026626",
		Title: "Google AI Blog",
		Url:   "http://ai.googleblog.com/",
		Articles: []feed.Article{
			{
				Id:          "tag:blogger.com,1999:blog-8474926331452026626.post-537064785672594983",
				Title:       "Mapping Africa’s Buildings with Satellite Imagery",
				Url:         "http://feedproxy.google.com/~r/blogspot/gJZg/~3/bcEqeVSMnBQ/mapping-africas-buildings-with.html",
				Published:   testutils.ParseTime(t, "2021-07-28T08:27:00.004-07:00"),
				Description: "Posted by ",
			},
		},
	}
	assert.Equal(t, asJson(t, expected), asJson(t, actual))
}

func TestNetflixArticles(t *testing.T) {
	actual := getFeed(t, "../parser/rss/testdata/rss_netflix_techblog.xml")
	expected := feed.Feed{
		Id:    "https://netflixtechblog.com?source=rss----2615bd06b42e---4",
		Title: "Netflix TechBlog - Medium",
		Url:   "https://netflixtechblog.com?source=rss----2615bd06b42e---4",
		Articles: []feed.Article{
			{
				Id:          "https://medium.com/p/3fddcceb1059",
				Title:       "Data Movement in Netflix Studio via Data Mesh",
				Url:         "https://netflixtechblog.com/data-movement-in-netflix-studio-via-data-mesh-3fddcceb1059?source=rss----2615bd06b42e---4",
				Published:   testutils.ParseTime(t, "2021-07-26T18:00:56+00:00"),
				Description: "By Andrew ",
			},
		},
	}
	assert.Equal(t, asJson(t, expected), asJson(t, actual))
}

func TestMuratArticles(t *testing.T) {
	actual := getFeed(t, "../parser/atom/testdata/atom_murat.xml")
	expected := feed.Feed{
		Id:    "tag:blogger.com,1999:blog-8436330762136344379",
		Title: "Metadata",
		Url:   "http://muratbuffalo.blogspot.com/",
		Articles: []feed.Article{
			{
				Id:          "tag:blogger.com,1999:blog-8436330762136344379.post-8449165989112346419",
				Title:       "There is plenty of room at the bottom",
				Url:         "http://muratbuffalo.blogspot.com/2021/08/there-is-plenty-of-room-at-bottom.html",
				Published:   testutils.ParseTime(t, "2021-08-17T09:35:00.008-04:00"),
				Description: "This is a ",
			},
		},
	}
	assert.Equal(t, asJson(t, expected), asJson(t, actual))
}

func TestDropboxArticles(t *testing.T) {
	actual := getFeed(t, "../parser/rss/testdata/rss_dropbox.xml")
	expected := feed.Feed{
		Id:    "https://dropbox.tech/feed",
		Title: "dropbox.tech",
		Url:   "https://dropbox.tech/feed",
		Articles: []feed.Article{
			{
				Id:          "https://dropbox.tech/infrastructure/making-dropbox-data-centers-carbon-neutral",
				Title:       "How we’re making Dropbox data centers 100% carbon neutral",
				Url:         "https://dropbox.tech/infrastructure/making-dropbox-data-centers-carbon-neutral",
				Published:   testutils.ParseTime(t, "2021-08-03T06:00:00-07:00"),
				Description: "As you may",
			},
		},
	}
	assert.Equal(t, asJson(t, expected), asJson(t, actual))
}

func TestFixRelativeUrls(t *testing.T) {
	actual := &feed.Feed{
		Url: "http://example.com/abc",
		Articles: []feed.Article{
			{
				Url: "",
			},
			{
				Url: "foo",
			},
			{
				Url: "/bar",
			},
			{
				Url: "http://example.com/xoxoxo",
			},
			{
				Url: "http://localhost",
			},
		},
	}
	expected := &feed.Feed{
		Url: "http://example.com/abc",
		Articles: []feed.Article{
			{
				Url: "http://example.com",
			},
			{
				Url: "http://example.com/foo",
			},
			{
				Url: "http://example.com/bar",
			},
			{
				Url: "http://example.com/xoxoxo",
			},
			{
				Url: "http://localhost",
			},
		},
	}

	feed.FixRelativeUrls(actual)
	assert.Equal(t, expected, actual)
}

func asJson(t *testing.T, a interface{}) string {
	t.Helper()
	j, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	return string(j)
}

func getFeed(t *testing.T, path string) feed.Feed {
	f := parseFeedFromXml(t, path)
	processArticles(f.Articles)
	f.Articles = f.Articles[:1]
	return f
}

func parseFeedFromXml(t *testing.T, path string) feed.Feed {
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
		t.Log("xxx here?")
		return atomFeed.AsGenericFeed()
	}
	rssFeed, rssErr := rss.Parse(b)
	if rssErr == nil && len(rssFeed.Items) > 0 {
		return rssFeed.AsGenericFeed("")
	}
	t.Fatalf("No articles. Atom err: %s, RSS err: %s", atomErr, rssErr)
	return feed.Feed{}
}

func processArticles(articles []feed.Article) {
	for i := range articles {
		articles[i].Description = fmt.Sprintf("%.10s", articles[i].Description)
	}
}
