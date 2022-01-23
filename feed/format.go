package feed

import (
	"bytes"
	"crypto/md5"
	_ "embed"
	"fmt"
	"html/template"
	"math"
	"net/url"
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
		"prettyUrl": func(u string) string {
			return formatPrettyUrl(u)
		},
		"hash": func(s string) string {
			return fmt.Sprintf("%x", md5.Sum([]byte(s)))
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

func formatPrettyUrl(rawUrl string) string {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}
	u.Fragment = ""
	u.RawQuery = ""
	return u.String()
}
