package renderers

import (
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	pandoc "willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils/pandoc"
)

type Renderer interface {
	Filename(pageTitle string) string
	Setup(dump *mwdump.XMLDump, outDir string) []error
	RenderCmd() *pandoc.Cmd
	RenderExec(cmd *pandoc.Cmd, page *mwdump.Page) (rendered string, errs []error)
}
