package rss

import (
	"encoding/xml"
	"bulletin/feed"
	"time"
)

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
	Title       string  `xml:"title"`
	Description string  `xml:"description"`
	Items       []*Item `xml:"item"`
}

func (c *Channel) GetArticles() []feed.Article {
	var articles []feed.Article
	for _, t := range c.Items {
		a := feed.Article{
			Id:    t.Guid,
			Title: t.Title,
			//Updated: t.
			Url: t.Link,
		}
		articles = append(articles, a)
	}
	return articles
}

type Item struct {
	Title          string   `xml:"title"`
	Link           string   `xml:"link"`
	Guid           string   `xml:"guid"`
	ContentEncoded string   `xml:"encoded"` // content:encoded
	PubDate        *RssTime `xml:"pubDate"`
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
	t, err := time.Parse(time.RFC1123, value)
	if err != nil {
		return err
	}
	*x = RssTime{t}
	return nil
}

func (x *RssTime) String() string {
	return x.Format(time.RFC3339)
}
