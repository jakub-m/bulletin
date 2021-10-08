package notion

import (
	"bulletin/feed"
	"bulletin/log"
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type FeedParser struct {
}

var _ feed.FeedParser = (*FeedParser)(nil)

func (p *FeedParser) Name() string {
	return "notion"
}

func (p *FeedParser) ParseFeed(body []byte) (feed.Feed, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return feed.Feed{}, err
	}
	var parsePostGrid func(*html.Node)

	parsePostGrid = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" && hasClass(n.Attr, "posts-grid") {
			log.Infof("xxx here %s", n.Data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parsePostGrid(c)
		}
	}
	parsePostGrid(doc)

	return feed.Feed{}, nil
}

func hasClass(attrs []html.Attribute, className string) bool {
	for _, attr := range attrs {
		fmt.Printf("xxx %+v\n", attr)
		for _, val := range strings.Split(attr.Val, " ") {
			if val == className {
				return true
			}
		}
	}
	return false
}
