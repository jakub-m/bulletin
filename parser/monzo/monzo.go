// Dear Monzo developers. Please add RSS to your awesome tech blog.
package monzo

import (
	"bulletin/feed"
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

const monzoBlogUrl = "https://monzo.com/blog/"

var FeedParser feed.FeedParser = &monzoFeedParser{}

type monzoFeedParser struct {
}

var _ feed.FeedParser = (*monzoFeedParser)(nil)

func (p *monzoFeedParser) Name() string {
	return "monzo"
}

var nilFeed feed.Feed = feed.Feed{}

func (p *monzoFeedParser) ParseFeed(body []byte, url string) (feed.Feed, error) {
	if !strings.HasPrefix(url, monzoBlogUrl) {
		return nilFeed, ErrBadUrl
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nilFeed, err
	}

	blogTitle := "?"
	if n := findFirstNode(doc, hasTag("h1")); n != nil {
		if m := findFirstNode(n, isTextNode()); m != nil {
			blogTitle = m.Data
		}
	}

	nodePostListing := findUlPostListing(doc)
	if nodePostListing == nil {
		return nilFeed, ErrCouldNotParse
	}

	articles := []feed.Article{}
	for _, li := range findAllLiPosts(nodePostListing) {
		article := getArticleFromNode(li)
		articles = append(articles, article)
	}

	return feed.Feed{
		Id:       "Monzo - " + blogTitle,
		Title:    "Monzo - " + blogTitle,
		Url:      url,
		Articles: articles,
	}, nil
}

func findUlPostListing(doc *html.Node) *html.Node {
	return findFirstNode(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "ul" && hasAttr(n, "class", startingWith("PostListing"))
	})
}

func findAllLiPosts(doc *html.Node) []*html.Node {
	return findAllNodesRec(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "li" && hasAttr(n, "class", startingWith("PostListing"))
	})
}

func getArticleFromNode(doc *html.Node) feed.Article {
	article := feed.Article{}

	if n := findFirstNode(doc, hasTag("h2")); n != nil {
		title := firstChildNodeText(n)
		article.Id = title
		article.Title = title
	}
	if n := findFirstNode(doc, hasTag("p")); n != nil {
		article.Description = firstChildNodeText(n)
	}
	if n := findFirstNode(doc, hasTag("time")); n != nil {
		//article.Description = firstChildNodeText(n)
		//time here
		timeString := strings.Trim(firstChildNodeText(n), " ")
		if t, err := time.Parse("2 January 2006", timeString); err == nil {
			article.Published = t
		}
	}
	if n := findFirstNode(doc, hasTag("a")); n != nil {
		href := getAttrValue(n.Attr, "href")
		article.Url = href
	}

	return article
}

type nodeMatcher func(*html.Node) bool

func hasTag(tag string) nodeMatcher {
	return func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == tag
	}
}

func isTextNode() nodeMatcher {
	return func(n *html.Node) bool {
		return n.Type == html.TextNode
	}
}

func firstChildNodeText(doc *html.Node) string {
	if t := findFirstNode(doc, func(n *html.Node) bool {
		return n.Type == html.TextNode
	}); t != nil {
		return t.Data
	}
	return ""
}

func findFirstNode(n *html.Node, fn nodeMatcher) *html.Node {
	var f func(*html.Node) *html.Node
	f = func(n *html.Node) *html.Node {
		if fn(n) {
			return n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if r := f(c); r != nil {
				return r
			}
		}
		return nil
	}
	return f(n)
}

func findAllNodesRec(n *html.Node, fn nodeMatcher) []*html.Node {
	var f func(*html.Node) []*html.Node
	f = func(n *html.Node) []*html.Node {
		nodes := []*html.Node{}
		if fn(n) {
			nodes = append(nodes, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if r := f(c); r != nil {
				nodes = append(nodes, r...)
			}
		}
		return nodes
	}
	return f(n)
}

type stringMatcher func(string) bool

func startingWith(prefix string) stringMatcher {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}

func hasAttr(n *html.Node, key string, matcher stringMatcher) bool {
	for _, attr := range n.Attr {
		if attr.Key == key {
			for _, part := range strings.Split(attr.Val, " ") {
				if matcher(part) {
					return true
				}
			}
		}
	}
	return false
}

func getAttrValue(attributes []html.Attribute, key string) string {
	for _, a := range attributes {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
