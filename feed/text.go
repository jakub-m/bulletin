package feed

import (
	"regexp"
	"strings"

	"bulletin/log"

	"golang.org/x/net/html"
)

var regexWhitespaces = regexp.MustCompile("[ \t\r\n]+")

const descriptionLength = 500

func GetDescriptionFromHTML(htmlBody string) string {
	d := ExtractTextFromHTML(htmlBody)
	return TrimSentences(d, descriptionLength)
}

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

// TrimSentences trims the input string to at most n charactes, preserving full sentences.
func TrimSentences(input string, n int) string {
	sentences := strings.SplitAfter(input, ".")
	t := ""
	for _, s := range sentences {
		if len(t)+len(s) <= n {
			t = t + s
		}
	}
	if t == "" {
		// if there is no "." in the input
		t = input
	}
	if len(t) > n {
		return input[:n]
	}
	return t
}
