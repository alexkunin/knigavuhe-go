package dom

import (
	"golang.org/x/net/html"
	"slices"
	"strings"
)

func findAttribute(node *html.Node, key string) *html.Attribute {
	for _, attribute := range node.Attr {
		if attribute.Key == key {
			return &attribute
		}
	}
	return nil
}

func IsTag(tag string) func(*html.Node) bool {
	return func(node *html.Node) bool {
		return node.Type == html.ElementNode && node.Data == tag
	}
}

func HasAttrWithValue(name string, value string) func(*html.Node) bool {
	return func(node *html.Node) bool {
		if attr := findAttribute(node, name); attr != nil {
			return attr.Val == value
		}
		return false
	}
}

func HasClass(class string) func(*html.Node) bool {
	return func(node *html.Node) bool {
		if attr := findAttribute(node, "class"); attr != nil {
			words := strings.Fields(attr.Val)
			return len(words) > 0 && slices.Contains(words, class)
		}
		return false
	}
}

func HasImmediateChild(tag string) func(*html.Node) bool {
	return func(node *html.Node) bool {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && child.Data == tag {
				return true
			}
		}
		return false
	}
}

func FindFirst(node *html.Node, predicates []func(*html.Node) bool) *html.Node {
	all := true

	for predicate := range predicates {
		if !predicates[predicate](node) {
			all = false
			break
		}
	}

	if all {
		return node
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if res := FindFirst(child, predicates); res != nil {
			return res
		}
	}

	return nil
}

func GetContent(node *html.Node) (found bool, value string) {
	if node != nil && node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		return true, strings.TrimSpace(node.FirstChild.Data)
	}
	return false, ""
}

func GetAttrValue(node *html.Node, name string) (found bool, value string) {
	if attr := findAttribute(node, name); attr != nil {
		return true, strings.TrimSpace(attr.Val)
	}
	return false, ""
}
