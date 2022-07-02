// Dear Monzo developers. Please add RSS to your awesome tech blog.
package monzo

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
	if n := hp.FindFirstNode(doc, hp.HasTag("h1")); n != nil {
		if m := hp.FindFirstNode(n, hp.IsTextNode()); m != nil {
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
	return hp.FindFirstNode(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "ul" && hp.NodeHasAttr(n, "class", hp.StartingWith("PostListing"))
	})
}

func findAllLiPosts(doc *html.Node) []*html.Node {
	return hp.FindAllNodesRec(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "li" && hp.NodeHasAttr(n, "class", hp.StartingWith("PostListing"))
	})
}

func getArticleFromNode(doc *html.Node) feed.Article {
	article := feed.Article{}

	if n := hp.FindFirstNode(doc, hp.HasTag("h2")); n != nil {
		title := hp.FirstChildNodeText(n)
		article.Id = title
		article.Title = title
	}
	if n := hp.FindFirstNode(doc, hp.HasTag("p")); n != nil {
		article.Description = hp.FirstChildNodeText(n)
	}
	if n := hp.FindFirstNode(doc, hp.HasTag("time")); n != nil {
		//article.Description = firstChildNodeText(n)
		//time here
		timeString := strings.Trim(hp.FirstChildNodeText(n), " ")
		if t, err := time.Parse("2 January 2006", timeString); err == nil {
			article.Published = t
		}
	}
	if n := hp.FindFirstNode(doc, hp.HasTag("a")); n != nil {
		href := hp.GetAttrValue(n.Attr, "href")
		article.Url = href
	}

	return article
}
