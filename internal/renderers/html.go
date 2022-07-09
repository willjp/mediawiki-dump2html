package renderers

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/lithammer/dedent"
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

func (html *HTML) Filename(page *mwdump.Page) string {
	fileName := fmt.Sprint(page.Title, ".html")
	return string(utils.SanitizePath([]byte(fileName)))
}

// Hook that runs before dumping all pages. Not necessarily a pure function.
func (html *HTML) Setup(dump *mwdump.XMLDump, outDir string) error {
	return renderStylesheet(dump, outDir)
}

// Renders one page to HTML, returns as string.
func (html *HTML) Render(page *mwdump.Page) (rendered string, err error) {
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
	render := incrHeaders(renderRaw)

	// end of html
	footer := `</html>`
	return fmt.Sprint(header, title, render, footer), nil
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
