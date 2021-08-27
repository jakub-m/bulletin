package command

import (
	"bulletin/feed"
	"bulletin/feedparser"
	"bulletin/fetcher"
	"bulletin/log"
	"fmt"
)

const TestCommandName = "test"

type TestCommand struct {
}

func (c *TestCommand) Execute(args []string) error {
	urls := args
	if len(urls) == 0 {
		log.Infof("pass urls of the feeds to test as positional arguments")
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
	body, err := fetcher.Get(url)
	log.Debugf("Got %d KB", len(body)/1<<10)
	if err != nil {
		return nil, err
	}
	return feedparser.GetArticles(body)
}
