package test

import (
	"github.com/willjp/mediawiki-dump2html/internal/elements/mwdump"
	"github.com/willjp/mediawiki-dump2html/internal/interfaces"
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
