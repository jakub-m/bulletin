package rss

import "encoding/xml"

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
	Items       []*Item `xml:"item"`
}

type Item struct {
	Title          string `xml:"title"`
	Link           string `xml:"link"`
	Guid           string `xml:"guid"`
	ContentEncoded string `xml:"encoded"` // content:encoded
	//PubDate RssTime `xml:"pubDate"`
}
