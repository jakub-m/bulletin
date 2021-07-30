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

// FetchCommand fetches feed from a single source provided directly in the command line.
type FetchCommand struct {
}

func (c *FetchCommand) Execute(commonOpts Options, args []string) error {
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
	feedCache, err := cache.NewCache(commonOpts.CacheDir)
	if err != nil {
		return err
	}
	for _, a := range articles {
		err := feedCache.StoreArticle(a)
		if err != nil {
			return err
		}
	}
	fmt.Println(html)
	return nil
}

func getOptions(args []string) (options, error) {
	var options options
	fs := flag.NewFlagSet("fetch", flag.ContinueOnError)
	fs.StringVar(&options.url, "url", "", "the feed to fetch")
	err := fs.Parse(args)
	return options, err
}

type options struct {
	url string
}
