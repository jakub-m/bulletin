package command

import (
	"bulletin/atom"
	"bulletin/cache"
	"bulletin/feed"
	"bulletin/fetcher"
	"bulletin/log"
	"bulletin/rss"
	"flag"
	"fmt"
)

const FetchCommandName = "fetch"

// FetchCommand fetches feed from a single source provided directly in the command line.
type FetchCommand struct {
	Cache *cache.Cache
}

func (c *FetchCommand) Execute(args []string) error {
	opts, err := getFetchOptions(args)
	if err != nil {
		return err
	}
	log.Infof("Fetch feed from %s", opts.url)
	feedBody, err := fetcher.Get(opts.url)
	if err != nil {
		return err
	}
	articles, err := parseArticles(feedBody)
	if err != nil {
		return err
	}
	log.Infof("Caching %d articles", len(articles))
	for _, a := range articles {
		err := c.Cache.StoreArticle(a)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseArticles(feedBody []byte) ([]feed.Article, error) {
	atomFeed, atomErr := atom.Parse(feedBody)
	if atomErr == nil && len(atomFeed.GetArticles()) > 0 {
		return atomFeed.GetArticles(), nil
	}
	rssFeed, rssErr := rss.Parse(feedBody)
	if rssErr == nil && len(rssFeed.GetArticles()) > 0 {
		return rssFeed.GetArticles(), nil
	}
	return nil, fmt.Errorf("could not parse. Atom error: %s. Rss error: %s", atomErr, rssErr)
}

func getFetchOptions(args []string) (fetchOptions, error) {
	var options fetchOptions
	fs := flag.NewFlagSet(FetchCommandName, flag.ContinueOnError)
	fs.StringVar(&options.url, "url", "", "the feed to fetch")
	err := fs.Parse(args)
	return options, err
}

type fetchOptions struct {
	url string
}
