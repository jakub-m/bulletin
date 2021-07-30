package feed

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"
)

var htmlTemplateArticles *template.Template

func init() {
	t, err := template.New("articles").Parse(htmlTemplateArticlesRaw)
	if err != nil {
		log.Fatal("feed: cannot parse html template")
	}
	htmlTemplateArticles = t
}

func FormatHtml(feeds []Article) (string, error) {
	buf := new(bytes.Buffer)
	err := htmlTemplateArticles.Execute(buf, feeds)
	if err != nil {
		return "", fmt.Errorf("feed: %s", err)
	}
	return strings.Trim(buf.String(), "\n"), nil
}

const htmlTemplateArticlesRaw = `
{{range .}}<a href="{{.Url}}">{{.Title}}</a></br>
{{end}}
`
