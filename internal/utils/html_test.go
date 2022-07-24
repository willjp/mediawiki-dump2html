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

func TestFindFirst(t *testing.T) {
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
	result := FindFirst(node, func(node *html.Node) *html.Node {
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
