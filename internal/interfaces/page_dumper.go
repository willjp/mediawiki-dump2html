package interfaces

import "github.com/willjp/mediawiki-dump2html/internal/elements/mwdump"

// Writes a single page to disk
type PageDumper interface {
	Dump(renderer Renderer, page *mwdump.Page, outPath string) (errs []error)
}
