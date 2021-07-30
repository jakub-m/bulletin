package feed

import "time"

type WithArticles interface {
	GetArticles() []Article
}

// Article is a generic feed agnostic to the original channel (Atom or RSS).
type Article struct {
	// Id identifies same articles. Two articles with the same Id will be included in the feed summary only once.
	Id      string
	Title   string
	Updated time.Time
	Url     string
}
