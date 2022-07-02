// Dear Ayke. Please add RSS to your awesome tech blog.
package aykevl

import (
	"bulletin/feed"
	hp "bulletin/htmlparser"
	"bytes"
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var (
	ErrCouldNotParse = fmt.Errorf("could not parse")
	ErrBadUrl        = fmt.Errorf("bad url")
)

const blogUrl = "https://aykevl.nl/"

var FeedParser feed.FeedParser = &aykevlFeedParser{}

type aykevlFeedParser struct {
}

func (p *aykevlFeedParser) Name() string {
	return "aykevl.nl"
}

var nilFeed feed.Feed = feed.Feed{}

func (p *aykevlFeedParser) ParseFeed(body []byte, url string) (feed.Feed, error) {
	if !strings.HasPrefix(url, blogUrl) {
		return nilFeed, ErrBadUrl
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nilFeed, err
	}

	blogTitle := "?"
	if n := hp.FindFirstNode(doc, hp.HasTag("header")); n != nil {
		if m := hp.FindFirstNode(n, hp.HasAttr("class", hp.StringIs("text"))); m != nil {
			blogTitle = hp.FirstChildNodeText(m)
		}
	}

	articles := []feed.Article{}

	for _, art := range hp.FindAllNodesRec(doc, hp.HasTag("article")) {
		article := feed.Article{}
		if a := hp.FindFirstNode(art, hp.HasTag("a")); a != nil {
			title := hp.FirstChildNodeText(a)
			article.Id = title
			article.Title = title
			article.Url = hp.GetAttrValue(a.Attr, "href")
		}

		if t := hp.FindFirstNode(art, hp.HasTag("time")); t != nil {
			datetime := hp.GetAttrValue(t.Attr, "datetime")
			tt, err := time.Parse(time.RFC3339, datetime)
			if err != nil {
				return nilFeed, fmt.Errorf("failed to parse time from feed: %v", err)
			}
			article.Published = tt
		}

		for _, p := range hp.FindAllNodesRec(art, hp.HasTag("p")) {
			if text := hp.FirstChildNodeText(p); len(text) > 0 {
				article.Description = text
			}
		}

		articles = append(articles, article)
	}

	return feed.Feed{
		Id:       url,
		Title:    blogTitle,
		Url:      url,
		Articles: articles,
	}, nil
}
