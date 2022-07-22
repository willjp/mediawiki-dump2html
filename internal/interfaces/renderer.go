package interfaces

import (
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
)

type Renderer interface {
	Filename(pageTitle string) string
	Setup(dump *mwdump.XMLDump, outDir string) []error
	Render(page *mwdump.Page) (render string, errs []error)
}