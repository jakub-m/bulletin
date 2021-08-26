package feed

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"math"
	"sort"
	"time"
)

const bulletinHeaderTimeFormat = `Monday, 02 Jan 2006`

//go:embed page_template.gohtml
var bulletinPageTemplateRaw string

type templateData struct {
	GroupedArticles [][]Article
	BulletinDate    string
	PeriodDays      string
}

func FormatHtml(periodDays int, now time.Time, articles []Article) (string, error) {
	buf := new(bytes.Buffer)
	grouped := groupArticlesPerFeed(articles)
	templateData := templateData{
		GroupedArticles: grouped,
		BulletinDate:    now.Local().Format(bulletinHeaderTimeFormat),
		PeriodDays:      formatDays(periodDays),
	}
	funcMap := template.FuncMap{
		"articleDate": func(a Article) string {
			return formatArticleDate(now, a)
		},
	}
	bulletinPageTemplate, err := template.New("page").Funcs(funcMap).Parse(bulletinPageTemplateRaw)
	if err != nil {
		log.Fatal("feed: cannot parse html template")
	}
	err = bulletinPageTemplate.Execute(buf, templateData)
	if err != nil {
		return "", fmt.Errorf("feed: %s", err)
	}
	return buf.String(), nil
}

func groupArticlesPerFeed(articles []Article) [][]Article {
	// sort articles per descending date.
	sort.Slice(articles, func(i, j int) bool {
		return !articles[i].Updated.Before(articles[j].Updated)
	})

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

func formatArticleDate(bulletinTime time.Time, article Article) string {
	dt := int(math.Round(bulletinTime.Sub(article.Updated).Hours() / 24.))
	return formatDays(dt) + " old"
}

func formatDays(days int) string {
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}
