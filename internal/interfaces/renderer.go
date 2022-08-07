package interfaces

import (
	"github.com/willjp/mediawiki-dump2html/internal/elements/mwdump"
)

// Renders a mediawiki XML dump to various formats (ex. html, rst, ..)
type Renderer interface {
	Filename(pageTitle string) string
	Setup(dump *mwdump.XMLDump, outDir string) []error
	Render(page *mwdump.Page) (render string, errs []error)
}
