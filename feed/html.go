package feed

import (
	"regexp"
	"strings"

	"bulletin/log"

	"golang.org/x/net/html"
)

var regexWhitespaces = regexp.MustCompile("[ \t\r\n]+")

// ExtractTextFromHTML recursively extracts all the text from the HTML input.
func ExtractTextFromHTML(htmlBody string) string {
	node, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		log.Debugf("Error while parsing HTML: %s", err)
		return ""
	}
	var f func(n *html.Node) string
	f = func(n *html.Node) string {
		nodeText := ""
		if n.Type == html.TextNode {
			nodeText = n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			nodeText = nodeText + " " + f(c)
		}
		return nodeText
	}
	t := f(node)
	t = regexWhitespaces.ReplaceAllString(t, " ")
	t = strings.Trim(t, " ")
	return t
}
