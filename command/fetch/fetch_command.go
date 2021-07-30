package fetch

import (
	"feedsummary/atom"
	"feedsummary/cache"
	"feedsummary/feed"
	"feedsummary/fetcher"
	"feedsummary/log"
	"flag"
	"fmt"
)

const CommandName = "fetch"

// Command fetches feed from a single source provided directly in the command line.
type Command struct {
	Cache *cache.Cache
}

func (c *Command) Execute(args []string) error {
	opts, err := getOptions(args)
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

func getOptions(args []string) (options, error) {
	var options options
	fs := flag.NewFlagSet(CommandName, flag.ContinueOnError)
	fs.StringVar(&options.url, "url", "", "the feed to fetch")
	err := fs.Parse(args)
	return options, err
}

type options struct {
	url string
}