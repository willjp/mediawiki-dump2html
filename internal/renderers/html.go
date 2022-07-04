package renderers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

type HTML struct{}

func (html *HTML) Filename(page *elements.Page) string {
	fileName := fmt.Sprint(page.Title, ".html")
	return string(utils.SanitizePath([]byte(fileName)))
}

func (html *HTML) Render(page *elements.Page) (rendered string, err error) {
	title := fmt.Sprintf("<h1 id=\"%s\">%s</h1>\n",
		toHtmlId(page.Title),
		page.Title)

	renderRaw, err := utils.PandocConvert(page, "mediawiki", "html")
	if err != nil {
		return "", err
	}
	render := incrHeaders(renderRaw)

	return fmt.Sprint(title, render), nil
}

var headerRx *regexp.Regexp
var idInvalidRx *regexp.Regexp

func init() {
	headerRx = regexp.MustCompile(fmt.Sprint(
		`(?P<head><[/]?[ \t]*h)`, // '<h'  '</h'
		`(?P<lv>[1-6])`,          // '1'
		`(?P<tail>[^>]*>)`,       // '>'
	))
	idInvalidRx = regexp.MustCompile(`[^a-z0-9\-]+`)
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
