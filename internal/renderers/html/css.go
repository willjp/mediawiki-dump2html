package renderers

import (
	"encoding/xml"
	"os"
	"path"
	"strings"

	"github.com/lithammer/dedent"
	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"

	htmlElement "willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/html"
)

// Writes CSS file that can be sourced in dumped HTML files.
func RenderStylesheet(dump *mwdump.XMLDump, outDir string) error {
	if len(dump.Pages) < 1 {
		return nil
	}

	extractor := newCssExtractor(&dump.Pages[0])
	css, err := extractor.Execute()
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

// Extracts pandoc's generated CSS from a page render.
//
// When pandoc is called using the `--standalone` param, it renders CSS into each page.
// This extracts that CSS, so that you could dump it to a file and reference it within each page.
type cssExtractor struct {
	pandoc *utils.Pandoc
}

func newCssExtractor(page *mwdump.Page) *cssExtractor {
	pandoc := utils.Pandoc{
		From:       "mediawiki",
		To:         "html",
		Standalone: true,
		Stdin:      strings.NewReader(page.LatestRevision().Text),
	}
	return &cssExtractor{pandoc: &pandoc}
}

func (this *cssExtractor) Execute() (string, error) {
	raw, err := this.pandoc.Execute()
	if err != nil {
		return "", err
	}
	css, err := this.extract(raw)
	if err != nil {
		return "", err
	}
	return css, nil
}

func (this *cssExtractor) execute() (string, error) {
	return this.pandoc.Execute()
}

func (this *cssExtractor) extract(raw string) (string, error) {
	var htmlNode htmlElement.Html
	xml.Unmarshal([]byte(raw), &htmlNode)
	css := dedent.Dedent(htmlNode.Head.Style)
	return css, nil
}
