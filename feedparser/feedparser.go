package feedparser

import (
	"bulletin/atom"
	"bulletin/feed"
	"bulletin/rss"
	"fmt"
)

var parsers []feed.FeedParser = []feed.FeedParser{
	atom.FeedParser,
	rss.FeedParser,
}

func GetFeed(feedBody []byte) (feed.Feed, error) {
	var errs []error
	for _, p := range parsers {
		if f, err := p.ParseFeed(feedBody); err == nil {
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
