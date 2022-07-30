package interfaces

import "willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"

// Writes a single page to disk
type PageDumper interface {
	Dump(renderer Renderer, page *mwdump.Page, outPath string) (errs []error)
}
