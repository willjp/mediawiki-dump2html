package interfaces

import "github.com/willjp/mediawiki-dump2html/internal/elements/mwdump"

// Writes output from a Renderer to disk
type RenderWriter interface {
	DumpAll(renderer Renderer, dump *mwdump.XMLDump, outDir string) (errs []error)
}
