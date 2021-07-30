package command

import (
	"feedsummary/atom"
	"feedsummary/cache"
	"feedsummary/feed"
	"feedsummary/fetcher"
	"feedsummary/log"
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
	atomFeed, err := atom.Parse(feedBody)
	if err != nil {
		return err
	}
	articles := atomFeed.GetArticles()
	html, err := feed.FormatHtml(articles)
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
	fmt.Println(html)
	return nil
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
