package feed

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"math"
	"sort"
	"strings"
	"time"
)

const bulletinHeaderTimeFormat = `Monday, 02 Jan 2006`

//go:embed page_template.gohtml
var bulletinPageTemplateRaw string

type templateData struct {
	GroupedArticles [][]Article
	BulletinEndDate string
	PeriodDays      string
	Feeds           []Feed
}

func FormatFeedsAsHtml(periodDays int, periodEnd time.Time, pageTemplate *string, feeds []Feed) (string, error) {
	buf := new(bytes.Buffer)
	templateData := templateData{
		Feeds:           feeds,
		BulletinEndDate: periodEnd.Local().Format(bulletinHeaderTimeFormat),
		PeriodDays:      formatDays(periodDays),
	}
	funcMap := template.FuncMap{
		"articleDate": func(a Article) string {
			return formatArticleDate(periodEnd, a)
		},
	}
	pageTemplateBody := bulletinPageTemplateRaw
	if pageTemplate != nil {
		pageTemplateBody = *pageTemplate
	}
	bulletinPageTemplate, err := template.New("page").Funcs(funcMap).Parse(pageTemplateBody)
	if err != nil {
		return "", err
	}
	err = bulletinPageTemplate.Execute(buf, templateData)
	if err != nil {
		return "", fmt.Errorf("feed: %s", err)
	}
	return buf.String(), nil

}

// template is the gohtml rendering template. If missing, defaults to built-in template.
// DEPRECATE
func FormatHtml(periodDays int, periodEnd time.Time, pageTemplate *string, articles []Article) (string, error) {
	buf := new(bytes.Buffer)
	grouped := groupArticlesPerFeed(articles)
	templateData := templateData{
		GroupedArticles: grouped,
		BulletinEndDate: periodEnd.Local().Format(bulletinHeaderTimeFormat),
		PeriodDays:      formatDays(periodDays),
	}
	funcMap := template.FuncMap{
		"articleDate": func(a Article) string {
			return formatArticleDate(periodEnd, a)
		},
	}
	pageTemplateBody := bulletinPageTemplateRaw
	if pageTemplate != nil {
		pageTemplateBody = *pageTemplate
	}
	bulletinPageTemplate, err := template.New("page").Funcs(funcMap).Parse(pageTemplateBody)
	if err != nil {
		return "", err
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
		return !articles[i].Published.Before(articles[j].Published)
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
		return strings.ToLower(feeds[i].Title) < strings.ToLower(feeds[j].Title)
	})

	var groupedArticles [][]Article
	for _, f := range feeds {
		groupedArticles = append(groupedArticles, articlesByFeedId[f.Id])
	}
	return groupedArticles
}

func formatArticleDate(bulletinTime time.Time, article Article) string {
	dt := int(math.Round(bulletinTime.Sub(article.Published).Hours() / 24.))
	return formatDays(dt) + " old"
}

func formatDays(days int) string {
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}
