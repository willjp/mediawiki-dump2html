package test

import (
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

type FakePageDumper struct {
	DumpErrors []error
	DumpCalled bool
	DumpArgs   []any
}

func (this *FakePageDumper) Dump(renderer interfaces.Renderer, page *mwdump.Page, outPath string) []error {
	this.DumpCalled = true
	this.DumpArgs = []any{renderer, page, outPath}
	return this.DumpErrors
}
