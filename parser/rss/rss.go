package rss

import (
	"bulletin/feed"
	"bulletin/log"
	btime "bulletin/time"
	"encoding/xml"
	"fmt"
	"net/url"
	"time"
)

type rssFeedParser struct {
}

var FeedParser feed.FeedParser = &rssFeedParser{}

func (p *rssFeedParser) Name() string {
	return "RSS"
}

func (p *rssFeedParser) ParseFeed(body []byte) (feed.Feed, error) {
	ch, err := Parse(body)
	if err == nil && ch == nil {
		err = fmt.Errorf("rss parser returned nil")
	}
	if err != nil {
		return feed.Feed{}, fmt.Errorf("rssFeedParser: %v", err)
	}
	return ch.AsGenericFeed(), nil
}

func Parse(raw []byte) (*Channel, error) {
	var r rssFeed
	err := xml.Unmarshal(raw, &r)
	return r.Channel, err
}

type rssFeed struct {
	Channel *Channel `xml:"channel"`
}

// Reed is RSS feed.
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
	Links       []Link `xml:"link"`
}

type Item struct {
	Title          string   `xml:"title"`
	Description    string   `xml:"description"`
	Link           string   `xml:"link"`
	Guid           string   `xml:"guid"`
	ContentEncoded string   `xml:"encoded"` // content:encoded
	PubDate        *RssTime `xml:"pubDate"`
}

type Link struct {
	Value string `xml:",chardata"`
	Href  string `xml:"href,attr"`
	Rel   string `xml:"rel,attr"`
	Type  string `xml:"type,attr"`
}

func (c *Channel) AsGenericFeed() feed.Feed {
	var articles []feed.Article
	feedLink := getBestLink(c.Links)
	feedTitle := c.Title
	if feedTitle == "" {
		if u, err := url.Parse(feedLink); err == nil {
			feedTitle = u.Host
		} else {
			log.Debugf("could not parse url %s: %s", feedLink, err)
		}
	}
	gf := feed.Feed{
		Id:    feedLink,
		Title: feedTitle,
		Url:   feedLink,
	}
	for _, t := range c.Items {
		a := feed.Article{
			Id:          t.Guid,
			Title:       t.Title,
			Description: getDescription(t),
			Published:   t.PubDate.Time,
			Url:         t.Link,
		}
		articles = append(articles, a)
	}
	gf.Articles = articles
	return gf
}

func getBestLink(links []Link) string {
	basicLinks := filterLinks(links, func(l Link) bool {
		return l.Value != ""
	})
	if len(basicLinks) > 0 {
		return firstLinkValue(basicLinks)
	}
	selfLinks := filterLinks(links, func(l Link) bool {
		return l.Rel == "self"
	})
	if len(selfLinks) == 1 {
		return firstLinkValue(selfLinks)
	}
	return firstLinkValue(links)
}

func firstLinkValue(links []Link) string {
	if len(links) == 0 {
		return ""
	}
	first := links[0]
	if first.Href != "" {
		return first.Href
	}
	return first.Value
}

func filterLinks(links []Link, fn func(l Link) bool) []Link {
	var filtered []Link
	for _, l := range links {
		if fn(l) {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

func getDescription(i Item) string {
	d := feed.GetDescriptionFromHTML(i.Description)
	if d == "" {
		return feed.GetDescriptionFromHTML(i.ContentEncoded)
	}
	return d
}

type RssTime struct {
	time.Time
}

func (x *RssTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	return x.parse(s)
}

func (x *RssTime) parse(value string) error {
	t, err := btime.Parse(value)
	if err != nil {
		return err
	}
	*x = RssTime{t}
	return nil
}

func (x *RssTime) String() string {
	return x.Format(time.RFC3339)
}
