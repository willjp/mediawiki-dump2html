package test

import (
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

type FakeRenderWriter struct {
	DumpAllErrors []error
	DumpErrors    []error
	DumpAllCalled bool
	DumpCalled    bool
}

func (this *FakeRenderWriter) DumpAll(renderer interfaces.Renderer, dump *mwdump.XMLDump, outDir string) (errs []error) {
	this.DumpAllCalled = true
	return this.DumpAllErrors
}

func (this *FakeRenderWriter) Dump(renderer interfaces.Renderer, page *mwdump.Page, outPath string) []error {
	this.DumpCalled = true
	return this.DumpErrors
}
