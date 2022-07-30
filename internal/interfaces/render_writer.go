package interfaces

import "willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"

// Writes output from a Renderer to disk
type RenderWriter interface {
	DumpAll(renderer Renderer, dump *mwdump.XMLDump, outDir string) (errs []error)
}
