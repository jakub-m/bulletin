package atom

import (
	"bulletin/feed"
	"bulletin/log"
	btime "bulletin/time"
	"encoding/xml"
	"fmt"
	"time"
)

var FeedParser feed.FeedParser = &atomFeedParser{}

type atomFeedParser struct {
}

func (p *atomFeedParser) Name() string {
	return "Atom"
}

func (p *atomFeedParser) ParseFeed(body []byte, url string) (feed.Feed, error) {
	ch, err := Parse(body)
	if err == nil && ch == nil {
		err = fmt.Errorf("atom parser returned nil")
	}
	if err != nil {
		return feed.Feed{}, fmt.Errorf("atomFeedParser: %v", err)
	}
	return ch.AsGenericFeed(), nil
}

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

type Entry struct {
	// Id identifies the entry using a universally unique and permanent URI. Two entries in a feed can have the same
	// value for id if they represent the same entry at different points in time.
	Id        string   `xml:"id"`
	Title     string   `xml:"title"`
	Published *XmlTime `xml:"published"`
	Updated   *XmlTime `xml:"updated"`
	Links     []Link   `xml:"link"`
	Content   string   `xml:"content"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"` // TODO fix relative links
	Type string `xml:"type,attr"`
}

func (atomFeed Feed) AsGenericFeed() feed.Feed {
	articles := []feed.Article{}
	for _, e := range atomFeed.Entries {
		genericArt, err := e.asGenericArticle()
		if err != nil {
			log.Debugf("AsGenericFeed: feed %s :%v", atomFeed.Id, err)
			continue
		}
		articles = append(articles, genericArt)
	}
	feedUrl := getBestUrl(atomFeed.Links)
	gf := feed.Feed{
		Id:       atomFeed.Id,
		Title:    atomFeed.Title,
		Url:      feedUrl,
		Articles: articles,
	}
	return gf
}

func (e Entry) asGenericArticle() (feed.Article, error) {
	if e.Published == nil && e.Updated == nil {
		return feed.Article{}, fmt.Errorf("asGenericArticle: cannot determine date for the article")
	}
	var published time.Time
	if e.Published == nil {
		// It happens that Published is missing, so we do our best here.
		published = e.Updated.Time
	} else {
		published = e.Published.Time
	}
	articleUrl := getBestUrl(e.Links)
	description := feed.GetDescriptionFromHTML(e.Content)
	art := feed.Article{
		Id:          e.Id,
		Title:       e.Title,
		Url:         articleUrl,
		Published:   published,
		Description: description,
	}
	return art, nil
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
