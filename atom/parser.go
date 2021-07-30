package atom

import (
	"encoding/xml"
	"time"
)

func Parse(raw []byte) (Feed, error) {
	var feed Feed
	err := xml.Unmarshal(raw, &feed)
	return feed, err
}

type Feed struct {
	Id       string  `xml:"id"`
	Title    string  `xml:"title"`
	Subtitle string  `xml:"subtitle"`
	Entries  []Entry `xml:"entry"`
}

type Entry struct {
	Id        string   `xml:"id"`
	Title     string   `xml:"title"`
	Published *XmlTime `xml:"published"`
	Updated   *XmlTime `xml:"updated"`
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
