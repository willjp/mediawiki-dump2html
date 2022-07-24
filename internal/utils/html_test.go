package utils

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var htmlWhitespaceRx = regexp.MustCompile(`(?m)(^\s+|\n)`)

func TestFindFirstChild(t *testing.T) {
	rawHtml := htmlWhitespaceRx.ReplaceAllString(`
	<html>
	  <head>
	    <title>Main Page</title>
	    <link rel="stylesheet" href="style.css"/>
	  </head>
	  <body>
	    <h1 id="main_page">Main Page</h1>
	    <blockquote>
	      <h2 id="bar">Bar</h2>
	    </blockquote>
	    <h2 id="foo">Foo</h2>
	    <blockquote>
	      <h2 id="baz">Baz</h2>
	    </blockquote>
	  </body>
	</html>
	`, "")
	node, _ := html.Parse(strings.NewReader(rawHtml))
	result := FindFirstChild(node, func(node *html.Node) *html.Node {
		// CSS ex. body + h2
		if node.DataAtom != atom.H2 {
			return nil
		}
		if node.Parent == nil {
			return nil
		}
		if node.Parent.DataAtom != atom.Body {
			return nil
		}
		return node
	})
	expects := html.Attribute{Key: "id", Val: "foo"}
	assert.Equal(t, expects, result.Attr[0])
}

func TestParentHeirarchy(t *testing.T) {
	rawHtml := htmlWhitespaceRx.ReplaceAllString(`
	<html>
	  <head>
	    <title>Main Page</title>
	    <link rel="stylesheet" href="style.css"/>
	  </head>
	</html>
	`, "")
	node, _ := html.Parse(strings.NewReader(rawHtml))
	target := node.FirstChild.FirstChild
	heirarchy := ParentHeirarchy(target)
	expects := []*html.Node{
		target,
		node.FirstChild,
		node,
	}
	assert.Equal(t, expects, heirarchy)
}

func TestHasParentHeirarchy(t *testing.T) {
	rawHtml := htmlWhitespaceRx.ReplaceAllString(`
	<html>
	  <head>
	    <style id="style">
	      html {
		line-height: 1.5;
	      }
	    </style>
	  </head>
	</html>
	`, "")

	t.Run("Returns node when matches", func(t *testing.T) {
		node, _ := html.Parse(strings.NewReader(rawHtml))
		target := node.FirstChild.FirstChild.FirstChild // document.html.head.style
		wants := []atom.Atom{atom.Head, atom.Style}     // head.style
		result := HasParentHeirarchy(target, wants)
		expects := html.Attribute{Key: "id", Val: "style"}
		assert.Equal(t, expects, result.Attr[0])
	})

	t.Run("Returns nil when node does not match", func(t *testing.T) {
		node, _ := html.Parse(strings.NewReader(rawHtml))
		target := node.FirstChild.FirstChild        // document.html.head
		wants := []atom.Atom{atom.Head, atom.Style} // head.style
		result := HasParentHeirarchy(target, wants)
		assert.Nil(t, result)
	})
}
