package hakibenita

import (
	"bulletin/feed"
	hp "bulletin/htmlparser"
	"bulletin/time"
	"bytes"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

var (
	ErrCouldNotParse = fmt.Errorf("could not parse")
	ErrBadUrl        = fmt.Errorf("bad url")
)

const blogUrl = "https://hakibenita.com/"

var FeedParser feed.FeedParser = &hakibenitaParser{}

type hakibenitaParser struct {
}

var _ feed.FeedParser = (*hakibenitaParser)(nil)

func (p *hakibenitaParser) Name() string {
	return "hakibenita"
}

var nilFeed feed.Feed = feed.Feed{}

func (p *hakibenitaParser) ParseFeed(body []byte, url string) (feed.Feed, error) {
	if !strings.HasPrefix(url, blogUrl) {
		return nilFeed, ErrBadUrl
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nilFeed, err
	}

	blogTitle := "?"
	if n := hp.FindFirstNode(doc, hp.HasAttr("class", hp.StringIs("logo"))); n != nil {
		blogTitle = hp.FirstChildNodeText(n)
	}

	postList := hp.FindFirstNode(doc, hp.All(hp.HasTag("ul"), hp.HasAttr("id", hp.StringIs("post-list"))))
	if postList == nil {
		return nilFeed, ErrCouldNotParse
	}

	articles := []feed.Article{}
	for _, n := range hp.ListChildren(postList, hp.HasTag("li")) {
		article := getArticleFromNode(n)
		articles = append(articles, article)
	}

	return feed.Feed{
		Id:       blogTitle,
		Title:    blogTitle,
		Url:      url,
		Articles: articles,
	}, nil
}

func getArticleFromNode(doc *html.Node) feed.Article {
	article := feed.Article{}

	if n := hp.FindFirstNode(doc, hp.HasTag("header")); n != nil {
		if n := hp.FindFirstNode(n, hp.HasTag("h2")); n != nil {
			title := hp.FirstChildNodeText(n)
			article.Id = title
			article.Title = title
			if n := hp.FindFirstNode(n, hp.HasTag("a")); n != nil {
				article.Url = hp.GetAttrValue(n.Attr, "href")
			}
		}
	}

	if n := hp.FindFirstNode(doc, hp.HasTag("p")); n != nil {
		article.Description = hp.FirstChildNodeText(n)
	}

	if n := hp.FindFirstNode(doc, hp.All(hp.HasTag("time"), hp.HasAttr("class", hp.StringIs("published")))); n != nil {
		dt := hp.GetAttrValue(n.Attr, "datetime")
		if t, err := time.Parse(dt); err == nil {
			article.Published = t
		} else {
			log.Printf("failed to parse date: %+v", err)
		}
	}

	// if n := hp.FindFirstNode(doc, hp.HasTag("h2")); n != nil {
	// 	title := hp.FirstChildNodeText(n)
	// 	article.Id = title
	// 	article.Title = title
	// }
	// if n := hp.FindFirstNode(doc, hp.HasTag("p")); n != nil {
	// 	article.Description = hp.FirstChildNodeText(n)
	// }
	// if n := hp.FindFirstNode(doc, hp.HasTag("time")); n != nil {
	// 	//article.Description = firstChildNodeText(n)
	// 	//time here
	// 	timeString := strings.Trim(hp.FirstChildNodeText(n), " ")
	// 	if t, err := time.Parse("2 January 2006", timeString); err == nil {
	// 		article.Published = t
	// 	}
	// }
	// if n := hp.FindFirstNode(doc, hp.HasTag("a")); n != nil {
	// 	href := hp.GetAttrValue(n.Attr, "href")
	// 	article.Url = href
	// }

	return article
}
