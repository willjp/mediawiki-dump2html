package renderers

import (
	"fmt"
	"strings"

	"github.com/lithammer/dedent"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/pandoc"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

// Has methods for conversion, and keeps state used during conversion
type RST struct {
	pandocExecutor interfaces.PandocExecutor
}

func New() RST {
	executor := pandoc.Executor{}
	return RST{pandocExecutor: &executor}
}

func (this *RST) Filename(pageTitle string) string {
	fileName := fmt.Sprint(pageTitle, ".rst")
	return string(utils.SanitizeFilename([]byte(fileName)))
}

// Hook that runs before dumping all pages. Not necessarily a pure function.
func (this *RST) Setup(dump *mwdump.XMLDump, outDir string) []error {
	return nil
}

func (this *RST) Render(page *mwdump.Page) (render string, errs []error) {
	cmd := this.renderCmd()
	stdin := strings.NewReader(page.LatestRevision().Text)
	raw, errs := this.pandocExecutor.Execute(&cmd, stdin)
	if errs != nil {
		return "", errs
	}

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

	// replace '<br>' with something rst understands
	renderRaw := strings.ReplaceAll(raw, "<br>", ":raw-html:`<br/>`")

	return fmt.Sprint(directives, string(title), renderRaw), nil
}

func (this *RST) renderCmd() pandoc.Cmd {
	opts := pandoc.Opts{
		From: "mediawiki",
		To:   "rst",
	}
	return opts.Command()
}
