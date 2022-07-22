package renderers

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/pandoc"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

var validSchemeRx = regexp.MustCompile(`^(http|https|ftp|file|fax|mailto|tel)$`)

type HTML struct {
	pandocExecutor interfaces.PandocExecutor
}

func New() HTML {
	executor := pandoc.Executor{}
	return HTML{pandocExecutor: &executor}
}

// alt constructor for tests
func newHTML(pandocExecutor interfaces.PandocExecutor) HTML {
	return HTML{pandocExecutor: pandocExecutor}
}

func (html *HTML) Filename(pageTitle string) string {
	// downcase everything - mediawiki has some links that are not case sensitive
	fileName := strings.ToLower(fmt.Sprint(pageTitle, ".html"))
	return string(utils.SanitizeFilename([]byte(fileName)))
}

// Hook that runs before dumping all pages. Not necessarily a pure function.
func (this *HTML) Setup(dump *mwdump.XMLDump, outDir string) (errs []error) {
	cssFile := path.Join(outDir, stylesheetName)
	css := NewCSS(this.pandocExecutor)
	return css.WriteCssFile(dump, cssFile)
}

func (this *HTML) Render(page *mwdump.Page) (render string, errs []error) {
	cmd := this.renderCmd()
	stdin := strings.NewReader(page.LatestRevision().Text)
	raw, errs := this.pandocExecutor.Execute(&cmd, stdin)
	if errs != nil {
		return "", errs
	}

	// parses/modifies/re-renders HTML (correcting links, setting header-levels, ...)
	node, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		utils.LogWarnOn(err)
		errs = append(errs, err)
		return "", errs
	}
	node, _ = this.adjust(node, page)
	var finalRender strings.Builder
	html.Render(&finalRender, node)

	return finalRender.String(), nil
}

// Prepares Command
func (this *HTML) renderCmd() pandoc.Cmd {
	opts := pandoc.Opts{
		From: "mediawiki",
		To:   "html",
	}
	return opts.Command()
}

// Rebuilds HTML Tree, adjusted for serving over static html
func (this *HTML) adjust(node *html.Node, page *mwdump.Page) (*html.Node, error) {
	var err error

	// process
	node = this.adjustHeadNode(node, page)
	node = this.adjustBodyNode(node, page)
	node = this.adjustHNode(node, page)
	err = this.adjustAnchorNode(node)
	if err != nil {
		return nil, err
	}

	// recurse through and modify children
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		this.adjust(child, page)
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

func (this *HTML) adjustHNode(node *html.Node, page *mwdump.Page) *html.Node {
	if node.Type != html.ElementNode {
		return node
	}
	switch node.DataAtom {
	case atom.H1:
		// Don't modify the page's <h1>
		for _, attr := range node.Attr {
			if attr.Key != "id" {
				continue
			}
			if attr.Val == toHtmlId(page.Title) {
				return node
			}
		}
		node.DataAtom = atom.H2
		node.Data = "h2"
	case atom.H2:
		node.DataAtom = atom.H3
		node.Data = "h3"
	case atom.H3:
		node.DataAtom = atom.H4
		node.Data = "h4"
	case atom.H4:
		node.DataAtom = atom.H5
		node.Data = "h5"
	case atom.H5:
		node.DataAtom = atom.H6
		node.Data = "h6"
	}
	return node
}

// Modifes `<a href="">` elements, so they point to files we have written to disk.
//
//    - files on disk use POSIX compatible characters in filename
//    - since serving statically without webserver, appends '.html' to filename
func (this *HTML) adjustAnchorNode(node *html.Node) error {
	if node.Type != html.ElementNode {
		return nil
	}
	if node.DataAtom != atom.A {
		return nil
	}
	var attrs []html.Attribute
	for _, attr := range node.Attr {
		if attr.Key != "href" {
			attrs = append(attrs, attr)
		}

		if !attrSchemeValid(attr.Val) {
			attr.Val = this.Filename(attr.Val)
		}

		attrs = append(attrs, attr)
	}
	node.Attr = attrs

	return nil
}

// Returns true if provided URL has a scheme that is allowlisted
func attrSchemeValid(uri string) bool {
	target, err := url.Parse(uri)
	return err == nil && target.IsAbs() && validSchemeRx.MatchString(target.Scheme)
}
