package feed

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"strings"
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

func FormatHtml(feeds []Article) (string, error) {
	buf := new(bytes.Buffer)
	err := bulletinPageTemplate.Execute(buf, feeds)
	if err != nil {
		return "", fmt.Errorf("feed: %s", err)
	}
	return strings.Trim(buf.String(), "\n"), nil
}
