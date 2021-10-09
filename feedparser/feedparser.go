package feedparser

import (
	"bulletin/feed"
	"bulletin/parser/atom"
	"bulletin/parser/monzo"
	"bulletin/parser/rss"
	"fmt"
)

var parsers []feed.FeedParser = []feed.FeedParser{
	atom.FeedParser,
	rss.FeedParser,
	monzo.FeedParser,
}

func GetFeed(feedBody []byte, url string) (feed.Feed, error) {
	var errs []error
	for _, p := range parsers {
		if f, err := p.ParseFeed(feedBody, url); err == nil {
			if len(f.Articles) > 0 {
				return f, nil
			} else {
				errs = append(errs, fmt.Errorf("parsed but no articles"))
			}
		} else {
			errs = append(errs, err)
		}
	}
	return feed.Feed{}, fmt.Errorf("could not parse: %v", errs)
}
