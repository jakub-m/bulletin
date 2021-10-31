package command

import (
	"bulletin/feed"
	"bulletin/feedparser"
	"bulletin/fetcher"
	"bulletin/log"
	"fmt"
	"io"
	"os"
)

const TestCommandName = "test"

type TestCommand struct {
}

func (c *TestCommand) Execute(args []string) error {
	urls := args
	if len(urls) == 0 {
		log.Infof("pass URLs or file paths of the feeds to test as positional arguments")
		return nil
	}
	for _, url := range urls {
		log.Infof("testing %s", url)
		if articles, err := getArticles(url); err == nil {
			fmt.Printf("good\t%s\t%d articles\n", url, len(articles))
		} else {
			fmt.Printf("BAD\t%s\t%s\n", url, err)
		}
	}
	return nil
}

func getArticles(url string) ([]feed.Article, error) {
	body, err := fetchOrRead(url)
	log.Debugf("Got %d KB", len(body)/1<<10)
	if err != nil {
		return nil, err
	}
	f, err := feedparser.GetFeed(body, url)
	if err != nil {
		return nil, err
	}
	return f.Articles, nil
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
