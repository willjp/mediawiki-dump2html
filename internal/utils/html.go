package utils

import "golang.org/x/net/html"

// Returns the first html node returned by the provided visitor.
func FindFirst(node *html.Node, visitor func(node *html.Node) *html.Node) *html.Node {
	match := visitor(node)
	if match != nil {
		return match
	}

	// recurse through and modify children
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		match = FindFirst(child, visitor)
		if match != nil {
			return match
		}
	}
	return nil
}
