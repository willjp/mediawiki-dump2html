package utils

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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
// Ex. []*html.Node{this, this.Parent, this.Parent.Parent, ...}
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

// Creates a lambda function that returns node, if it's parent heirarchy matches provided.
// Ex. HasParentHeirarchy(node, []atom.Atom{Atom.Head, Atom.Style})
func HasParentHeirarchy(node *html.Node, heirarchy []atom.Atom) *html.Node {
	parents := ParentHeirarchy(node)
	if parents == nil {
		return nil
	}
	if len(parents) < len(heirarchy) {
		return nil
	}
	parentAtoms := Map(
		parents[:len(heirarchy)],
		func(node *html.Node) atom.Atom {
			return node.DataAtom
		},
	)
	for i, parentAtom := range parentAtoms {
		if parentAtom != heirarchy[len(heirarchy)-i-1] {
			return nil
		}
	}
	return node
}
