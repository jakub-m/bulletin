package feedparser

import (
	"bulletin/atom"
	"bulletin/feed"
	"bulletin/rss"
	"fmt"
)

// GetArticles parses raw XML feed body. DEPRECATED.
func GetArticles(feedBody []byte) ([]feed.Article, error) {
	atomFeed, atomErr := atom.Parse(feedBody)
	if atomErr == nil && len(atomFeed.GetArticles()) > 0 {
		return atomFeed.GetArticles(), nil
	}
	rssFeed, rssErr := rss.Parse(feedBody)
	if rssErr == nil && len(rssFeed.GetArticles()) > 0 {
		return rssFeed.GetArticles(), nil
	}
	return nil, fmt.Errorf("could not parse. Atom error: %s. RSS error: %s", atomErr, rssErr)
}

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
