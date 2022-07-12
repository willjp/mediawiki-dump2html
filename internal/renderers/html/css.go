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

	cmd := PandocCommand()
	css, err := ExtractCss(cmd, &dump.Pages[0])
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

// Builds pandoc command to render HTML with CSS.
func PandocCommand() *utils.PandocCmd {
	opts := utils.PandocOpts{
		From:       "mediawiki",
		To:         "html",
		Standalone: true,
	}
	cmd := opts.Command()
	return cmd
}

// Executes pandoc command, and extracts CSS
func ExtractCss(cmd *utils.PandocCmd, src *mwdump.Page) (string, error) {
	html, err := cmd.Execute(strings.NewReader(src.LatestRevision().Text))
	if err != nil {
		return "", err
	}

	var htmlNode htmlElement.Html
	xml.Unmarshal([]byte(html), &htmlNode)
	css := dedent.Dedent(htmlNode.Head.Style)
	if err != nil {
		return "", err
	}
	return css, nil
}
