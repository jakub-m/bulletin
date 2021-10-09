package monzo

import (
	"bulletin/feed"
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type FeedParser struct {
}

var _ feed.FeedParser = (*FeedParser)(nil)

func (p *FeedParser) Name() string {
	return "monzo"
}

func (p *FeedParser) ParseFeed(body []byte) (feed.Feed, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return feed.Feed{}, err
	}

	nodePostListing := findPostListing(doc)

	fmt.Printf("xxx %v\n", nodePostListing)
	return feed.Feed{}, nil
}

func findPostListing(doc *html.Node) *html.Node {
	var f func(*html.Node) *html.Node
	f = func(n *html.Node) *html.Node {
		fmt.Printf("xxx data = %s, %v\n", n.Data, n.Type)
		if n.Type == html.ElementNode && n.Data == "ul" && hasClass(n.Attr, startingWith("PostListing")) {
			return n
		}
		fmt.Printf("xxx firstchild %v\n", n.FirstChild)

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			return f(c)
		}
		return nil
	}
	return f(doc)
}

type stringMatcher func(string) bool

func hasClass(attrs []html.Attribute, matcher stringMatcher) bool {
	for _, attr := range attrs {
		fmt.Printf("xxx key %s\n", attr.Key)
		if attr.Key == "class" {
			for _, part := range strings.Split(attr.Val, " ") {
				if matcher(part) {
					return true
				}
			}
		}
	}
	return false
}

func startingWith(prefix string) stringMatcher {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}
