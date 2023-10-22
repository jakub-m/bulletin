package feedparser

import (
	"bulletin/feed"
	"bulletin/log"
	"bulletin/parser/atom"
	"bulletin/parser/aykevl"
	"bulletin/parser/benoitessiambre"
	"bulletin/parser/hakibenita"
	"bulletin/parser/monzo"
	"bulletin/parser/rss"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

var parsers []feed.FeedParser = []feed.FeedParser{
	atom.FeedParser,
	aykevl.FeedParser,
	benoitessiambre.FeedParser,
	hakibenita.FeedParser,
	monzo.FeedParser,
	rss.FeedParser,
}

func GetFeed(feedBody []byte, url string) (feed.Feed, error) {
	var errs []error
	for _, p := range parsers {
		log.Debugf("With %s try %s", url, p.Name())
		if f, err := p.ParseFeed(feedBody, url); err == nil {
			for i := range f.Articles {
				cleanupArticle(&f.Articles[i])
			}
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

func cleanupArticle(article *feed.Article) {
	article.Title = strings.Trim(article.Title, " \t\n\r")
	article.Id = hash(article.Id)
}

func hash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
