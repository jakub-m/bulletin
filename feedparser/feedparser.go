package feedparser

import (
	"bulletin/atom"
	"bulletin/feed"
	"bulletin/rss"
	"fmt"
)

// GetArticles parses raw XML feed body.
func GetArticles(feedBody []byte) ([]feed.Article, error) {
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
