package renderers

import (
	"encoding/xml"
	"strings"

	"github.com/lithammer/dedent"
	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/pandoc"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"

	htmlElement "willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/html"
)

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
	logger.Infof("Writing: %s\n", filepath)
	return utils.FileReplace(render, filepath)
}

// Renders CSS stylesheet for page.
func (this *CSS) Render(dump *mwdump.XMLDump) (render string, errs []error) {
	if len(dump.Pages) < 1 {
		return "", nil
	}

	cmd := this.pandocCommand()
	stdin := strings.NewReader(dump.Pages[0].LatestRevision().Text)
	html, errs := this.pandocExecutor.Execute(&cmd, stdin)
	if errs != nil {
		return "", errs
	}

	var htmlNode htmlElement.Html
	xml.Unmarshal([]byte(html), &htmlNode)
	css := dedent.Dedent(htmlNode.Head.Style)
	if errs != nil {
		return "", errs
	}

	return css, nil
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
