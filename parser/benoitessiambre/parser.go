// Dear Benoit, please add RSS feed to your blog. Thanks!

package benoitessiambre

import (
	"bulletin/feed"
	"bulletin/time"
	"bytes"
	"fmt"
	"strings"
	gotime "time"

	hp "github.com/jakub-m/htmlquery"

	"golang.org/x/net/html"
)

var (
	ErrCouldNotParse = fmt.Errorf("could not parse")
	ErrBadUrl        = fmt.Errorf("bad url")
)

const blogUrl = "https://benoitessiambre.com/blemish.html"

var FeedParser feed.FeedParser = &feedParser{}

type feedParser struct {
}

func (p *feedParser) Name() string {
	return "benoitessiambre.com"
}

var nilFeed feed.Feed = feed.Feed{}

func (p *feedParser) ParseFeed(body []byte, url string) (feed.Feed, error) {
	if !strings.HasPrefix(url, blogUrl) {
		return nilFeed, ErrBadUrl
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nilFeed, err
	}
	_ = doc

	blogTitle := p.Name()

	articles := []feed.Article{}

	if dl := hp.FindFirstNode(doc, hp.HasTag("dl")); dl != nil {
		times := []gotime.Time{}
		for _, timeNode := range hp.ListChildren(dl, hp.HasAttr("class", hp.StringIs("time"))) {
			t, err := time.Parse(hp.FirstChildNodeText(timeNode))
			if err != nil {
				// just carry on. I am sorry, I did my best.
				continue
			}
			times = append(times, t)
		}
		for i, dd := range hp.ListChildren(dl, hp.HasTag("dd")) {
			article := feed.Article{}
			if a := hp.FindFirstNode(dd, hp.HasTag("a")); a != nil {
				article.Url = hp.GetAttrValue(a.Attr, "href")
				article.Id = article.Url
				article.Title = hp.FirstChildNodeText(a)
			}
			if tag := hp.FindFirstNode(dd, hp.HasAttr("class", hp.StringIs("tag"))); tag != nil {
				article.Description = hp.FirstChildNodeText(tag)
			}
			if i >= len(times) {
				continue
			}
			article.Published = times[i]
			articles = append(articles, article)
		}
	}

	return feed.Feed{
		Id:       url,
		Title:    blogTitle,
		Url:      url,
		Articles: articles,
	}, nil
}
