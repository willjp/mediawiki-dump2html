package renderers

import "willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements"

type Renderer interface {
	Filename(page *elements.Page) string
	Setup(dump *elements.XMLDump, outDir string) error
	Render(page *elements.Page) (rendered string, err error)
}
