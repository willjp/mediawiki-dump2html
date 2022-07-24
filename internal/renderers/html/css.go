package renderers

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/log"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/pandoc"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

var UnableToFindCssError = errors.New("Unable to locate stylesheet within html")

type CSS struct {
	pandocExecutor interfaces.PandocExecutor
}

func NewCSS(pandocExecutor interfaces.PandocExecutor) CSS {
	return CSS{pandocExecutor: pandocExecutor}
}

// Convenience method that Renders, and writes file, collecting errors.
func (this *CSS) WriteCssFile(dump *mwdump.XMLDump, filepath string) []error {
	render, errs := this.Render(dump)
	if errs != nil {
		return errs
	}
	log.Log.Infof("Writing: %s\n", filepath)
	return utils.FileReplace(render, filepath)
}

// Renders CSS stylesheet for page.
func (this *CSS) Render(dump *mwdump.XMLDump) (render string, errs []error) {
	if len(dump.Pages) < 1 {
		return "", nil
	}

	cmd := this.pandocCommand()
	stdin := strings.NewReader(dump.Pages[0].LatestRevision().Text)
	rawHtml, errs := this.pandocExecutor.Execute(&cmd, stdin)
	if errs != nil {
		return "", errs
	}

	// extract css from parsed html
	css, err := this.extractCssFromHtml(rawHtml)
	if err != nil {
		return "", []error{err}
	}
	return css, nil
}

// extract CSS from first html.head.style element
func (this *CSS) extractCssFromHtml(rawHtml string) (render string, err error) {
	node, err := html.Parse(strings.NewReader(rawHtml))
	if err != nil {
		return "", err
	}
	cssNode := utils.FindFirstChild(node, func(node *html.Node) *html.Node {
		return utils.HasParentHeirarchy(node, []atom.Atom{atom.Head, atom.Style})
	})
	if cssNode == nil {
		return "", UnableToFindCssError
	}
	for child := cssNode.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			return child.Data, nil
		}
	}
	return "", UnableToFindCssError
}

// Builds pandoc command to render HTML with CSS.
func (this *CSS) pandocCommand() pandoc.Cmd {
	opts := pandoc.Opts{
		From:       "mediawiki",
		To:         "html",
		Standalone: true,
	}
	cmd := opts.Command()
	return cmd
}
