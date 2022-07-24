package utils

import (
	"golang.org/x/net/html"
)

// Returns the first html node returned by the provided visitor.
func FindFirstChild(node *html.Node, visitor func(node *html.Node) *html.Node) *html.Node {
	match := visitor(node)
	if match != nil {
		return match
	}

	// recurse through and modify children
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		match = FindFirstChild(child, visitor)
		if match != nil {
			return match
		}
	}
	return nil
}

// Returns of provided node's parents (including provided node).
func ParentHeirarchy(node *html.Node) (nodes []*html.Node) {
	if node != nil {
		nodes = append(nodes, node)
	}
	if node.Parent != nil {
		parent_nodes := ParentHeirarchy(node.Parent)
		nodes = append(nodes, parent_nodes...)
	}
	return nodes
}
