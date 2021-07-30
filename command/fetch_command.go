package command

import (
	"feedsummary/atom"
	"feedsummary/feed"
	"feedsummary/fetcher"
	"feedsummary/log"
	"flag"
	"fmt"
)

type Command interface {
	Execute(args []string) error
}

// FetchCommand fetches feed from a single source provided directly in the command line.
type FetchCommand struct {
}

func (c *FetchCommand) Execute(args []string) error {
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
	html, err := feed.FormatHtml(atomFeed.GetArticles())
	if err != nil {
		return err
	}
	fmt.Println(html)
	return nil
}

func getOptions(args []string) (options, error) {
	var options options
	fs := flag.NewFlagSet("FetchCommand", flag.ContinueOnError)
	fs.StringVar(&options.url, "url", "", "the feed to fetch")
	err := fs.Parse(args)
	return options, err
}

type options struct {
	url string
}
