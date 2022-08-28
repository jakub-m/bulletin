package command

import (
	"bulletin/feed"
	"bulletin/feedparser"
	"bulletin/fetcher"
	"bulletin/log"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"time"
)

const TestCommandName = "test"

var feedSuffixes = []string{
	"all.atom.xml",
	"atom.xml",
	"feed",
	"feed.atom",
	"feed.rss",
	"feed.xml",
	"index.xml",
	"rss",
	"rss.xml",
}

type TestCommand struct {
}

func (c *TestCommand) Execute(args []string) error {
	flagSingle := false
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.BoolVar(&flagSingle, "1", false, "try single url instead of discovering")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	urls := fs.Args()
	if !flagSingle {
		extendedUrls := []string{}
		for _, base := range urls {
			extendedUrls = append(extendedUrls, base)
			for _, suffix := range feedSuffixes {
				baseUrl, err := url.Parse(base)
				log.Debugf("base url for %s is %s", base, baseUrl)
				if err != nil {
					return nil
				}
				suffixUrl, err := url.Parse(suffix)
				log.Debugf("try suffix: %s", suffixUrl)
				if err != nil {
					return err
				}
				extended := baseUrl.ResolveReference(suffixUrl).String()
				log.Debugf("extended url: %s", extended)
				extendedUrls = append(extendedUrls, extended)
			}
		}
		urls = extendedUrls
	}

	if len(urls) == 0 {
		log.Infof("pass URLs or file paths of the feeds to test as positional arguments")
		return nil
	}
	for _, url := range urls {
		log.Infof("testing %s", url)
		if articles, err := getArticles(url); err == nil {
			sortArticlesByDateAsc(articles)
			latestArticle := articles[len(articles)-1]
			hoursSinceLast := time.Since(latestArticle.Published).Hours()
			fmt.Printf("good\t%s\t%d articles, latest %.0f days ago (%s)\n", url, len(articles), hoursSinceLast/24, latestArticle.Published)
		} else {
			log.Infof("BAD\t%s\t%s\n", url, err)
		}
	}
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
