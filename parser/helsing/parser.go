// Dear Helsing developers. Please add RSS to your tech blog.
package helsing

import (
	"bulletin/feed"
	"bytes"
	"fmt"
	"strings"
	"time"

	hp "github.com/jakub-m/htmlquery"

	"golang.org/x/net/html"
)

var (
	ErrCouldNotParse = fmt.Errorf("could not parse")
	ErrBadUrl        = fmt.Errorf("bad url")
)

const helsingBlogUrl = "https://blog.helsing.ai/"

var FeedParser feed.FeedParser = &helsingFeedParser{}

type helsingFeedParser struct {
}

var _ feed.FeedParser = (*helsingFeedParser)(nil)

func (p *helsingFeedParser) Name() string {
	return "helsing"
}

var nilFeed feed.Feed = feed.Feed{}

func (p *helsingFeedParser) ParseFeed(body []byte, url string) (feed.Feed, error) {
	if !strings.HasPrefix(url, helsingBlogUrl) {
		return nilFeed, ErrBadUrl
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nilFeed, err
	}

	blogTitle := "?"
	if n := hp.FindFirstNode(doc, hp.All(hp.HasTag("h1"), hp.HasAttr("class", hp.StringIs("site-title")))); n != nil {
		blogTitle = hp.FirstChildNodeText(n)
	}

	articleNodes := hp.FindAllNodesRec(doc, hp.All(hp.HasTag("article"), hp.HasAttr("class", hp.StringIs("post-item"))))
	if len(articleNodes) == 0 {
		return nilFeed, ErrCouldNotParse
	}

	articles := []feed.Article{}
	for _, an := range articleNodes {
		articles = append(articles, getArticleFromNode(an))
	}

	return feed.Feed{
		Id:       "Helsing - " + blogTitle,
		Title:    "Helsing - " + blogTitle,
		Url:      url,
		Articles: articles,
	}, nil
}

func getArticleFromNode(doc *html.Node) feed.Article {
	article := feed.Article{}

	if n := hp.FindFirstNode(doc, hp.All(hp.HasTag("a"), hp.HasAttr("class", hp.StringIs("post-title")))); n != nil {
		title := hp.FirstChildNodeText(n)
		article.Id = title
		article.Title = title
		article.Url = hp.GetAttrValue(n.Attr, "href")
	}
	if n := hp.FindFirstNode(doc, hp.HasTag("time")); n != nil {
		datetime := strings.TrimSpace(hp.GetAttrValue(n.Attr, "datetime"))
		if t, err := time.Parse("2006-01-02", datetime); err == nil {
			article.Published = t
		}
	}
	if n := hp.FindFirstNode(doc, hp.All(hp.HasTag("p"), hp.HasAttr("class", hp.StringIs("summary")))); n != nil {
		article.Description = hp.FirstChildNodeText(n)
	}

	return article
}
