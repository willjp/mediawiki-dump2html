package renderers

import "willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"

type Renderer interface {
	Filename(page *mwdump.Page) string
	Setup(dump *mwdump.XMLDump, outDir string) error
	Render(page *mwdump.Page) (rendered string, err error)
}
