package feedparser

import (
	"bulletin/atom"
	"bulletin/feed"
	"bulletin/rss"
	"fmt"
)

type feedParser func(feedBody []byte) (feed.Feed, error)

var parsers []feedParser = []feedParser{
	rssParser,
	atomParser,
}

func rssParser(feedBody []byte) (feed.Feed, error) {
	rssFeed, err := rss.Parse(feedBody)
	if err != nil {
		return feed.Feed{}, err
	}
	if rssFeed == nil {
		return feed.Feed{}, fmt.Errorf("RSS: parser returned nil")
	}
	return rssFeed.AsGenericFeed(), nil
}

func atomParser(feedBody []byte) (feed.Feed, error) {
	atomFeed, err := atom.Parse(feedBody)
	if err != nil {
		return feed.Feed{}, err
	}
	return atomFeed.AsGenericFeed(), nil
}

func GetFeed(feedBody []byte) (feed.Feed, error) {
	var errs []error
	for _, p := range parsers {
		if f, err := p(feedBody); err == nil {
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
