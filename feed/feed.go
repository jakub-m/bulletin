package feed

import (
	"encoding/json"
	"time"
)

type WithArticles interface {
	AsGenericFeed() Feed
}

// Feed is an aggregate of the articles.
type Feed struct {
	// Id uniquely identifies the Feed.
	Id    string
	Title string
	// Url from where the Feed was fetched from. Points to RSS or Atom XML document.
	Url      string
	Articles []Article
}

// Article is a generic feed agnostic to the original channel (Atom or RSS).
type Article struct {
	// Id identifies same articles. Two articles with the same Id will be included in the feed summary only once.
	Id    string
	Title string
	// Description is a short summary text (i.e. not HTML) of the Feed entry.
	Description string
	Published   time.Time
	// Url directs to the actual article.
	Url string
}

func (a Article) Marshall() ([]byte, error) {
	return json.MarshalIndent(a, "", " ")
}

func UnmarshallArticle(bytes []byte) (Article, error) {
	var a Article
	err := json.Unmarshal(bytes, &a)
	return a, err
}
