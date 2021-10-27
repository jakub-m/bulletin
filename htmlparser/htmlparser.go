package htmlparser

import (
	"strings"

	"golang.org/x/net/html"
)

type NodeMatcher func(*html.Node) bool

func FindFirstNode(n *html.Node, fn NodeMatcher) *html.Node {
	var f func(*html.Node) *html.Node
	f = func(n *html.Node) *html.Node {
		if fn(n) {
			return n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if r := f(c); r != nil {
				return r
			}
		}
		return nil
	}
	return f(n)
}

func FirstChildNodeText(doc *html.Node) string {
	if t := FindFirstNode(doc, func(n *html.Node) bool {
		return n.Type == html.TextNode
	}); t != nil {
		return t.Data
	}
	return ""
}

func FindAllNodesRec(n *html.Node, fn NodeMatcher) []*html.Node {
	var f func(*html.Node) []*html.Node
	f = func(n *html.Node) []*html.Node {
		nodes := []*html.Node{}
		if fn(n) {
			nodes = append(nodes, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if r := f(c); r != nil {
				nodes = append(nodes, r...)
			}
		}
		return nodes
	}
	return f(n)
}

func HasTag(tag string) NodeMatcher {
	return func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == tag
	}
}

func HasAttr(key string, matcher StringMatcher) NodeMatcher {
	return func(n *html.Node) bool {
		for _, attr := range n.Attr {
			if attr.Key == key {
				for _, part := range strings.Split(attr.Val, " ") {
					if matcher(part) {
						return true
					}
				}
			}
		}
		return false
	}
}

func IsTextNode() NodeMatcher {
	return func(n *html.Node) bool {
		return n.Type == html.TextNode
	}
}

type StringMatcher func(string) bool

// NodeHasAttr is deprecated, use HasAttr instead.
func NodeHasAttr(n *html.Node, key string, matcher StringMatcher) bool {
	for _, attr := range n.Attr {
		if attr.Key == key {
			for _, part := range strings.Split(attr.Val, " ") {
				if matcher(part) {
					return true
				}
			}
		}
	}
	return false
}

func StartingWith(prefix string) StringMatcher {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}

func StringIs(expected string) StringMatcher {
	return func(s string) bool {
		return s == expected
	}
}

func GetAttrValue(attributes []html.Attribute, key string) string {
	for _, a := range attributes {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
