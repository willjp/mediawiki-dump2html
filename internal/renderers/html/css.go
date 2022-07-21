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
	pandoc "willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils/pandoc"

	htmlElement "willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/html"
)

// Writes CSS file that can be sourced in dumped HTML files.
func RenderStylesheet(dump *mwdump.XMLDump, outDir string) []error {
	if len(dump.Pages) < 1 {
		return nil
	}

	cmd := PandocCommand()
	css, errs := ExtractCss(&cmd, &dump.Pages[0])
	if errs != nil {
		return errs
	}

	cssPath := path.Join(outDir, stylesheetName)
	file, err := os.Create(cssPath)
	utils.PanicOn(err)
	defer func() {
		err := file.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}()

	logger.Infof("Writing: %s\n", cssPath)
	_, err = file.WriteString(css)
	if err != nil {
		utils.RmFileOn(file, err)
		errs = append(errs, err)
		return errs
	}
	return nil
}

// Builds pandoc command to render HTML with CSS.
func PandocCommand() pandoc.Cmd {
	opts := pandoc.Opts{
		From:       "mediawiki",
		To:         "html",
		Standalone: true,
	}
	cmd := opts.Command()
	return cmd
}

// Executes pandoc command, and extracts CSS
func ExtractCss(cmd *pandoc.Cmd, src *mwdump.Page) (string, []error) {
	html, errs := cmd.Execute(strings.NewReader(src.LatestRevision().Text))
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
