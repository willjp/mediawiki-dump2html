package renderers

import (
	"fmt"
	"strings"

	"github.com/lithammer/dedent"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

// Has methods for conversion, and keeps state used during conversion
type RST struct{}

func (rst *RST) Filename(page *elements.Page) string {
	fileName := fmt.Sprint(page.Title, ".rst")
	return string(utils.SanitizePath([]byte(fileName)))
}

// Converts mediawiki text to rst, with tweaks so it behaves well with sphinx-docs.
func (rst *RST) Render(page *elements.Page) (rendered string, err error) {
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

	opts := utils.PandocOptions{From: "mediawiki", To: "rst"}
	pandocRender, err := utils.PandocConvert(page, &opts)
	if err != nil {
		return "", err
	}

	// replace '<br>' with something rst understands
	render := strings.ReplaceAll(pandocRender, "<br>", ":raw-html:`<br/>`")

	return fmt.Sprint(directives, string(title), render), nil
}
