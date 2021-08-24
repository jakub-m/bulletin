package atom

import (
	"bulletin/feed"
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
	Id       string   `xml:"id"`
	Title    string   `xml:"title"`
	Subtitle string   `xml:"subtitle"`
	Entries  []*Entry `xml:"entry"`
	Link     []*Link  `xml:"link"`
}

var _ feed.WithArticles = (*Feed)(nil)

type Entry struct {
	// Id identifies the entry using a universally unique and permanent URI. Two entries in a feed can have the same
	// value for id if they represent the same entry at different points in time.
	Id        string   `xml:"id"`
	Title     string   `xml:"title"`
	Published *XmlTime `xml:"published"`
	Updated   *XmlTime `xml:"updated"`
	OrigLink  string   `xml:"origLink"` // feedburner:origLink
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

// getUrl returns the most appropriate link.
func (f Feed) getUrl() string {
	if len(f.Link) == 1 {
		return f.Link[0].Href
	}
	for _, s := range []string{"alternate", "self", ""} {
		for _, l := range f.Link {
			if l.Rel == s {
				return l.Href
			}
		}

	}
	return ""
}

func (e Entry) asArticle(f Feed) feed.Article {
	updated := e.Published.Time
	if e.Updated != nil {
		updated = e.Updated.Time
	}
	ff := feed.Feed{
		Id:    f.Id,
		Title: f.Title,
		Url:   f.getUrl(),
	}
	return feed.Article{
		Feed:    ff,
		Id:      e.Id,
		Title:   e.Title,
		Url:     e.OrigLink,
		Updated: updated,
	}
}

type XmlTime struct {
	time.Time
}

const timeFormat = `2006-01-02T15:04:05.000-07:00`

func (x *XmlTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	return x.parse(s)
}

func (x *XmlTime) parse(value string) error {
	t, err := time.Parse(timeFormat, value)
	if err != nil {
		return err
	}
	*x = XmlTime{t}
	return nil
}

func (x *XmlTime) String() string {
	return x.Format(time.RFC3339)
}
