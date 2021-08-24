package rss

import (
	"bulletin/feed"
	"encoding/xml"
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
	Link        string  `xml:"link"`
}

func (c *Channel) GetArticles() []feed.Article {
	var articles []feed.Article
	f := feed.Feed{
		// Id is implemented as Title for RSS because we cannot extract link reliably. There are <link>
		//and <atom:link> that confuse XML parser and make it return zero element.
		Id:    c.Title,
		Title: c.Title,
		Url:   c.Link,
	}
	for _, t := range c.Items {
		a := feed.Article{
			Feed:    f,
			Id:      t.Guid,
			Title:   t.Title,
			Updated: t.PubDate.Time,
			Url:     t.Link,
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
