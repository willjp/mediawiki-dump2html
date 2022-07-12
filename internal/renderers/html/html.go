package renderers

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

type HTML struct{}

func (html *HTML) Filename(pageTitle string) string {
	// downcase everything - mediawiki has some links that are not case sensitive
	fileName := strings.ToLower(fmt.Sprint(pageTitle, ".html"))
	return string(utils.SanitizeFilename([]byte(fileName)))
}

// Hook that runs before dumping all pages. Not necessarily a pure function.
func (html *HTML) Setup(dump *mwdump.XMLDump, outDir string) error {
	return RenderStylesheet(dump, outDir)
}

// Renders one page to HTML, returns as string.
func (this *HTML) Render(page *mwdump.Page) (rendered string, err error) {
	// rendered wiki
	opts := utils.PandocOpts{
		From: "mediawiki",
		To:   "html",
	}
	cmd := opts.Command()
	renderRaw, err := cmd.Execute(strings.NewReader(page.LatestRevision().Text))
	if err != nil {
		return "", err
	}

	// parses/modifies/re-renders HTML (correcting links, setting header-levels, ...)
	node, err := html.Parse(strings.NewReader(renderRaw))
	if err != nil {
		utils.LogWarnOn(err)
		return "", err
	}
	node, _ = this.adjust(node, page)
	var render strings.Builder
	html.Render(&render, node)

	return render.String(), nil
}

// Rebuilds HTML Tree, adjusted for serving over static html
func (this *HTML) adjust(node *html.Node, page *mwdump.Page) (*html.Node, error) {
	var err error

	// process
	node = this.adjustHeadNode(node, page)
	node = this.adjustBodyNode(node, page)
	node, err = this.adjustAnchorNode(node)
	if err != nil {
		return nil, err
	}

	// recurse through children
	var children []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		child, err = this.adjust(child, page)
		if err != nil {
			return child, err
		}
		children = append(children, child)
	}

	// point Child/Sibling info in structs to the new children
	if len(children) > 0 {
		node.FirstChild = children[0]
		node.LastChild = children[len(children)-1]
	}
	for index, child := range children {
		if 0 < index && index < len(children)-1 {
			child.PrevSibling = children[index-1]
			child.NextSibling = children[index+1]
		}
	}

	return node, nil
}

// Adjusts HTML <head>
//
//    - Links stylesheet
//    - Sets Page Title.
func (this *HTML) adjustHeadNode(node *html.Node, page *mwdump.Page) *html.Node {
	if node.Type != html.ElementNode {
		return node
	}
	if node.DataAtom != atom.Head {
		return node
	}

	titleVal := html.Node{
		Type:      html.TextNode,
		Data:      page.Title,
		Namespace: node.Namespace,
	}

	title := html.Node{
		Type:       html.ElementNode,
		DataAtom:   atom.Title,
		Data:       "title",
		Namespace:  node.Namespace,
		FirstChild: &titleVal,
		LastChild:  &titleVal,
	}

	linkStyle := html.Node{
		Type:      html.ElementNode,
		DataAtom:  atom.Link,
		Data:      "link",
		Namespace: node.Namespace,
		Attr: []html.Attribute{
			{Namespace: node.Namespace, Key: "rel", Val: "stylesheet"},
			{Namespace: node.Namespace, Key: "href", Val: stylesheetName},
		},
	}

	node.AppendChild(&title)
	node.AppendChild(&linkStyle)
	return node
}

// Adjusts HTML Body
//
//    - Prepends a <h1> with the page title to the body.
func (this *HTML) adjustBodyNode(node *html.Node, page *mwdump.Page) *html.Node {
	if node.Type != html.ElementNode {
		return node
	}
	if node.DataAtom != atom.Body {
		return node
	}

	headerVal := html.Node{
		Type:      html.TextNode,
		Data:      page.Title,
		Namespace: node.Namespace,
	}

	header := html.Node{
		Type:       html.ElementNode,
		DataAtom:   atom.H1,
		Data:       "h1",
		Namespace:  node.Namespace,
		FirstChild: &headerVal,
		LastChild:  &headerVal,
		Attr: []html.Attribute{
			{Namespace: node.Namespace, Key: "id", Val: toHtmlId(page.Title)},
		},
	}

	node.InsertBefore(&header, node.FirstChild)
	return node
}

// Modifes `<a href="">` elements, so they point to files we have written to disk.
//
//    - files on disk use POSIX compatible characters in filename
//    - since serving statically without webserver, appends '.html' to filename
func (this *HTML) adjustAnchorNode(node *html.Node) (finalNode *html.Node, err error) {
	if node.Type != html.ElementNode {
		return node, nil
	}
	if node.DataAtom != atom.A {
		return node, nil
	}

	var attrs []html.Attribute
	for _, attr := range node.Attr {
		if attr.Key != "href" {
			attrs = append(attrs, attr)
			continue
		}

		// ignore error, since we are correcting invalid urls
		target, err := url.Parse(attr.Val)
		if err != nil || !target.IsAbs() {
			newAttr := html.Attribute{
				Namespace: attr.Namespace,
				Key:       attr.Key,
				Val:       this.Filename(attr.Val),
			}
			attrs = append(attrs, newAttr)
		} else {
			attrs = append(attrs, attr)
		}
	}

	return &html.Node{
		Parent:      node.Parent,
		FirstChild:  node.FirstChild,
		LastChild:   node.LastChild,
		PrevSibling: node.PrevSibling,
		NextSibling: node.NextSibling,
		Type:        node.Type,
		DataAtom:    node.DataAtom,
		Data:        node.Data,
		Namespace:   node.Namespace,
		Attr:        attrs,
	}, nil
}
