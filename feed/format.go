package feed

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"sort"
)

//go:embed page_template.gohtml
var bulletinPageTemplateRaw string

var bulletinPageTemplate *template.Template

func init() {
	var err error
	bulletinPageTemplate, err = template.New("page").Parse(bulletinPageTemplateRaw)
	if err != nil {
		log.Fatal("feed: cannot parse html template")
	}
}

func FormatHtml(articles []Article) (string, error) {
	buf := new(bytes.Buffer)
	grouped := groupArticlesPerFeed(articles)
	err := bulletinPageTemplate.Execute(buf, grouped)
	if err != nil {
		return "", fmt.Errorf("feed: %s", err)
	}
	return buf.String(), nil
}

func groupArticlesPerFeed(articles []Article) [][]Article {
	feedIdMap := make(map[string]Feed)
	articlesByFeedId := make(map[string][]Article)
	for _, article := range articles {
		feedId := article.Feed.Id
		feedIdMap[feedId] = article.Feed
		articlesByFeedId[feedId] = append(articlesByFeedId[feedId], article)
	}
	var feeds []Feed
	for _, f := range feedIdMap {
		feeds = append(feeds, f)
	}
	sort.Slice(feeds, func(i, j int) bool {
		return feeds[i].Title < feeds[j].Title
	})

	var groupedArticles [][]Article
	for _, f := range feeds {
		groupedArticles = append(groupedArticles, articlesByFeedId[f.Id])
	}
	return groupedArticles
}
