package atom

import (
	"bulletin/feed"
	btime "bulletin/time"
	"encoding/xml"
	"time"
)

func Parse(raw []byte) (*Feed, error) {
	var feed Feed
	err := xml.Unmarshal(raw, &feed)
	return &feed, err
}

// Feed represents Atom feed. See schema in https://validator.w3.org/feed/docs/atom.html
type Feed struct {
	Id       string  `xml:"id"`
	Title    string  `xml:"title"`
	Subtitle string  `xml:"subtitle"`
	Entries  []Entry `xml:"entry"`
	Links    []Link  `xml:"link"`
}

var _ feed.WithArticles = (*Feed)(nil)

type Entry struct {
	// Id identifies the entry using a universally unique and permanent URI. Two entries in a feed can have the same
	// value for id if they represent the same entry at different points in time.
	Id        string   `xml:"id"`
	Title     string   `xml:"title"`
	Published *XmlTime `xml:"published"`
	Links     []Link   `xml:"link"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr"`
}

func (f Feed) GetArticles() []feed.Article {
	var articles []feed.Article
	for _, e := range f.Entries {
		articles = append(articles, e.asArticle(f))
	}
	return articles
}

func (e Entry) asArticle(atomFeed Feed) feed.Article {
	published := e.Published.Time
	feedUrl := getBestUrl(atomFeed.Links)
	f := feed.Feed{
		Id:    atomFeed.Id,
		Title: atomFeed.Title,
		Url:   feedUrl,
	}
	articleUrl := getBestUrl(e.Links)
	return feed.Article{
		Feed:      f,
		Id:        e.Id,
		Title:     e.Title,
		Url:       articleUrl,
		Published: published,
	}
}

// getUrl returns the most appropriate link.
func getBestUrl(links []Link) string {
	if len(links) == 1 {
		return links[0].Href
	}
	for _, rel := range []string{"alternate", "self", ""} {
		for _, l := range links {
			if l.Rel == rel {
				return l.Href
			}
		}
	}
	return ""
}

type XmlTime struct {
	time.Time
}

func (x *XmlTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	return x.parse(s)
}

func (x *XmlTime) parse(value string) error {
	t, err := btime.Parse(value)
	if err != nil {
		return err
	}
	*x = XmlTime{t}
	return nil
}

func (x *XmlTime) String() string {
	return x.Format(time.RFC3339)
}
