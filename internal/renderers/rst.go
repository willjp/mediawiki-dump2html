package renderers

import (
	"fmt"
	"strings"

	"github.com/lithammer/dedent"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

// Has methods for conversion, and keeps state used during conversion
type RST struct{}

func (rst *RST) Filename(pageTitle string) string {
	fileName := fmt.Sprint(pageTitle, ".rst")
	return string(utils.SanitizeFilename([]byte(fileName)))
}

// Hook that runs before dumping all pages. Not necessarily a pure function.
func (html *RST) Setup(dump *mwdump.XMLDump, outDir string) error {
	return nil
}

// Converts mediawiki text to rst, with tweaks so it behaves well with sphinx-docs.
func (rst *RST) Render(page *mwdump.Page) (rendered string, err error) {
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

	pandoc := utils.Pandoc{
		From:  "mediawiki",
		To:    "rst",
		Stdin: strings.NewReader(page.LatestRevision().Text),
	}
	pandocRender, err := pandoc.Execute()
	if err != nil {
		return "", err
	}

	// replace '<br>' with something rst understands
	render := strings.ReplaceAll(pandocRender, "<br>", ":raw-html:`<br/>`")

	return fmt.Sprint(directives, string(title), render), nil
}
