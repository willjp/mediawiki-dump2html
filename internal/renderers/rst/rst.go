package renderers

import (
	"fmt"
	"strings"

	"github.com/lithammer/dedent"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
	pandoc "willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils/pandoc"
)

// Has methods for conversion, and keeps state used during conversion
type RST struct{}

func (rst *RST) Filename(pageTitle string) string {
	fileName := fmt.Sprint(pageTitle, ".rst")
	return string(utils.SanitizeFilename([]byte(fileName)))
}

// Hook that runs before dumping all pages. Not necessarily a pure function.
func (html *RST) Setup(dump *mwdump.XMLDump, outDir string) []error {
	return nil
}

func (this *RST) RenderCmd() *pandoc.Cmd {
	opts := pandoc.Opts{
		From: "mediawiki",
		To:   "rst",
	}
	return opts.Command()
}

// Converts mediawiki text to rst, with tweaks so it behaves well with sphinx-docs.
func (rst *RST) RenderExec(cmd *pandoc.Cmd, page *mwdump.Page) (rendered string, errs []error) {
	directives := dedent.Dedent(`
	.. role:: raw-html(raw)
	  :format: html

	`)

	// page title between '='s
	titleLen := len([]rune(page.Title))
	title := fmt.Sprint(
		strings.Repeat("=", titleLen), "\n",
		page.Title, "\n",
		strings.Repeat("=", titleLen), "\n\n",
	)

	pandocRender, errs := cmd.Execute(strings.NewReader(page.LatestRevision().Text))
	if errs != nil {
		return "", errs
	}

	// replace '<br>' with something rst understands
	render := strings.ReplaceAll(pandocRender, "<br>", ":raw-html:`<br/>`")

	return fmt.Sprint(directives, string(title), render), nil
}
