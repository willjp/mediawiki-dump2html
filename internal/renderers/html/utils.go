package renderers

import (
	"regexp"
	"strings"
)

var idInvalidRx *regexp.Regexp

func init() {
	idInvalidRx = regexp.MustCompile(`[^a-z0-9\-]+`)
}

// Downcases, and sanitizes characters in a HTML header to assign to html ID.
//   (ex. 'My  Page' --> 'my_page')
func toHtmlId(value string) string {
	downcased := strings.ToLower(value)
	return idInvalidRx.ReplaceAllString(downcased, "_")
}

