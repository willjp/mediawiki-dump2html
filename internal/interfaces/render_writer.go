package interfaces

import "willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"

type RenderWriter interface {
	DumpAll(renderer Renderer, dump *mwdump.XMLDump, outDir string) (errs []error)
	Dump(renderer Renderer, page *mwdump.Page, outPath string) []error
}
