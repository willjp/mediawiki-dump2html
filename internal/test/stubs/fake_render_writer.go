package test

import (
	"github.com/willjp/mediawiki-dump2html/internal/elements/mwdump"
	"github.com/willjp/mediawiki-dump2html/internal/interfaces"
)

type FakeRenderWriter struct {
	DumpAllErrors []error
	DumpAllCalled bool

	// TODO: deleteme
	DumpErrors []error
	DumpCalled bool
}

func (this *FakeRenderWriter) DumpAll(renderer interfaces.Renderer, dump *mwdump.XMLDump, outDir string) (errs []error) {
	this.DumpAllCalled = true
	return this.DumpAllErrors
}

// TODO: deleteme
func (this *FakeRenderWriter) Dump(renderer interfaces.Renderer, page *mwdump.Page, outPath string) []error {
	this.DumpCalled = true
	return this.DumpErrors
}
