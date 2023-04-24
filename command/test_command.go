package command

import (
	"bulletin/feed"
	"bulletin/feedparser"
	"bulletin/fetcher"
	"bulletin/log"
	"bytes"
	"flag"
	"fmt"
	"io"
	gourl "net/url"
	"os"
	"sort"
	"time"

	hq "github.com/jakub-m/htmlquery"

	"golang.org/x/net/html"
)

const TestCommandName = "test"

type TestCommand struct {
}

func (c *TestCommand) Execute(args []string) error {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("expected exactly one url")
	}
	url := fs.Args()[0]

	body, err := fetcher.Get(url)
	if err != nil {
		return fmt.Errorf("error while fetching %s: %s", url, err)
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error while parsing %s: %s", url, err)
	}
	/*
	   <link rel="alternate" type="application/rss+xml" href="/blog/rss.xml" title="QuestDB Blog RSS Feed">
	   <link rel="alternate" type="application/atom+xml" href="/blog/atom.xml" title="QuestDB Blog Atom Feed">
	*/
	feedNodes := hq.FindAllNodesRec(doc,
		hq.All(hq.HasTag("link"),
			hq.Any(
				hq.HasAttr("type", hq.StringIs("application/rss+xml")),
				hq.HasAttr("type", hq.StringIs("application/atom+xml")),
			)))
	log.Debugf("Got %d feed nodes: %s", len(feedNodes), feedNodes)
	var feedUrl string
	if len(feedNodes) == 0 {
		feedUrl = url
	} else {
		feedUrl = getAttr(feedNodes[0], "href")
	}
	if feedUrl == "" {
		return fmt.Errorf("feed url missing for %s", url)
	}

	feedUrl, err = joinPaths(url, feedUrl)
	if err != nil {
		return fmt.Errorf("error while testing %s and %s: %s", url, feedUrl, err)
	}

	log.Infof("Testing %s", feedUrl)
	if articles, err := getArticles(feedUrl); err == nil {
		sortArticlesByDateAsc(articles)
		latestArticle := articles[len(articles)-1]
		hoursSinceLast := time.Since(latestArticle.Published).Hours()
		log.Infof("good\t%s\t%d articles, latest %.0f days ago (%s)\n", feedUrl, len(articles), hoursSinceLast/24, latestArticle.Published)
	} else {
		log.Infof("BAD\t%s\t%s\n", feedUrl, err)
	}
	fmt.Println(feedUrl)
	return nil
}

func getArticles(url string) ([]feed.Article, error) {
	body, err := fetchOrRead(url)
	log.Debugf("Got %d KB", len(body)/(1<<10))
	if err != nil {
		return nil, err
	}
	f, err := feedparser.GetFeed(body, url)
	if err != nil {
		return nil, err
	}
	return f.Articles, nil
}

func sortArticlesByDateAsc(articles []feed.Article) {
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Published.Before(articles[j].Published)
	})
}

func fetchOrRead(url string) ([]byte, error) {
	if _, err := os.Stat(url); err == nil {
		f, err := os.Open(url)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return io.ReadAll(f)
	} else {
		return fetcher.Get(url)
	}
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func joinPaths(url1, url2 string) (string, error) {
	u1, err := gourl.Parse(url1)
	if err != nil {
		return "", err
	}
	u2, err := gourl.Parse(url2)
	if err != nil {
		return "", err
	}
	if u2.Host != "" {
		return u2.String(), nil
	}
	return gourl.JoinPath(fmt.Sprintf("%s://%s", u1.Scheme, u1.Host), u2.Path)
}
