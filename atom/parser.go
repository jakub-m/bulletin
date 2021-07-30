package atom

import (
	"encoding/xml"
	"fmt"
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
	Entries  []*Entry `xml:"entry"`
}

type Entry struct {
	// Id identifies the entry using a universally unique and permanent URI. Two entries in a feed can have the same
	// value for id if they represent the same entry at different points in time.
	Id        string   `xml:"id"`
	Title     string   `xml:"title"`
	Published *XmlTime `xml:"published"`
	Updated   *XmlTime `xml:"updated"`
}

// Uid return unique Id of the Entry. If the same content was published or updated at different points in time,
// it will have a different Uid.
func (e *Entry) Uid() string {
	return fmt.Sprintf("%s-%s", e.Updated.UTC().Format(time.RFC3339), e.Id)
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