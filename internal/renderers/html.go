package renderers

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/lithammer/dedent"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

var headerRx *regexp.Regexp
var idInvalidRx *regexp.Regexp
var stylesheetName string

func init() {
	headerRx = regexp.MustCompile(fmt.Sprint(
		`(?P<head><[/]?[ \t]*h)`, // '<h'  '</h'
		`(?P<lv>[1-6])`,          // '1'
		`(?P<tail>[^>]*>)`,       // '>'
	))
	idInvalidRx = regexp.MustCompile(`[^a-z0-9\-]+`)
	stylesheetName = "style.css"
}

type HTML struct{}

func (html *HTML) Filename(pageTitle string) string {
	fileName := fmt.Sprint(pageTitle, ".html")
	return string(utils.SanitizePath([]byte(fileName)))
}

// Hook that runs before dumping all pages. Not necessarily a pure function.
func (html *HTML) Setup(dump *mwdump.XMLDump, outDir string) error {
	return renderStylesheet(dump, outDir)
}

// Renders one page to HTML, returns as string.
func (this *HTML) Render(page *mwdump.Page) (rendered string, err error) {
	// html header
	header := dedent.Dedent(fmt.Sprintf(`
		<html>
		<head>
		  <title>%s</title>
		  <link rel="stylesheet" href="%s" />
		</head>
		`, page.Title, stylesheetName,
	))

	// h1
	title := fmt.Sprintf(
		"<h1 id=\"%s\">%s</h1>\n",
		toHtmlId(page.Title),
		page.Title,
	)

	// rendered wiki
	opts := utils.PandocOptions{From: "mediawiki", To: "html"}
	renderRaw, err := utils.PandocConvert(page, &opts)
	if err != nil {
		return "", err
	}

	// parses/modifies HTML (corrects links, increments header-levels, ...)
	node, err := html.Parse(strings.NewReader(renderRaw))
	if err != nil {
		utils.LogWarnOn(err)
		return "", err
	}
	finalBody, _ := this.adjustHtmlNode(node)

	var render strings.Builder
	html.Render(&render, finalBody)

	// end of html
	footer := `</html>`
	return fmt.Sprint(header, title, render.String(), footer), nil
	return "", nil
}

func (this *HTML) adjustHtmlNode(node *html.Node) (finalNode *html.Node, err error) {
	// process
	newNode, err := this.adjustAnchorLinks(node)
	if err != nil {
		return nil, err
	}

	// mutate children
	var children []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		newChild, err := this.adjustHtmlNode(child)
		if err != nil {
			return nil, err
		}
		children = append(children, newChild)
	}

	// update this node so points to updated children
	if len(children) > 0 {
		newNode.FirstChild = children[0]
		newNode.LastChild = children[len(children)-1]
	}

	// now update children so they point to correct sibling
	for index, child := range children {
		if 0 < index && index < len(children)-1 {
			child.PrevSibling = children[index-1]
			child.NextSibling = children[index+1]
		}
	}

	return newNode, nil
}

// Increments the header-level of every HTML header in 'render'.
// (ex. <h1>foo</h1> --> <h2>foo</h2>)
func (this *HTML) incrementHeaderLvl(node *html.Node) (finalNode *html.Node, err error) {
	return node, nil
}

// Modifes `<a href="">` elements, so they point to files we have written to disk.
//  - files on disk use POSIX compatible characters in filename
//  - since serving statically without webserver, appends '.html' to filename
func (this *HTML) adjustAnchorLinks(node *html.Node) (finalNode *html.Node, err error) {
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

// Increments the header-level of every HTML header in 'render'.
// (ex. <h1>foo</h1> --> <h2>foo</h2>)
func incrHeaders(render string) string {
	return headerRx.ReplaceAllStringFunc(render, func(match string) string {
		submatches := headerRx.FindStringSubmatch(match)
		lv, err := strconv.Atoi(submatches[2])
		utils.PanicOn(err)
		return fmt.Sprint(submatches[1], lv+1, submatches[3])
	})
}

// Downcases, and sanitizes characters in a HTML header to assign to html ID.
// (ex. 'My  Page' --> 'my_page')
func toHtmlId(value string) string {
	downcased := strings.ToLower(value)
	return idInvalidRx.ReplaceAllString(downcased, "_")
}

// Writes CSS file that can be sourced in dumped HTML files.
func renderStylesheet(dump *mwdump.XMLDump, outDir string) error {
	if len(dump.Pages) < 1 {
		return nil
	}

	css, err := utils.PandocExtractCss(&dump.Pages[0])
	if err != nil {
		return err
	}
	cssPath := path.Join(outDir, stylesheetName)
	file, err := os.Create(cssPath)
	defer file.Close()
	utils.PanicOn(err)

	logger.Infof("Writing: %s\n", cssPath)
	_, err = file.WriteString(css)
	if err != nil {
		utils.RmFileOn(file, err)
		return err
	}
	return nil
}
