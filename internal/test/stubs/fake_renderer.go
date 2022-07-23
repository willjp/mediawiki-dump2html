package test

import "willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"

type FakeRenderer struct {
	Output       string
	SetupErrors  []error
	RenderErrors []error
	SetupCalled  bool
}

func (this *FakeRenderer) Filename(pageTitle string) string {
	return pageTitle
}

func (this *FakeRenderer) Setup(dump *mwdump.XMLDump, outDir string) []error {
	this.SetupCalled = true
	return this.SetupErrors
}

func (this *FakeRenderer) Render(page *mwdump.Page) (render string, errs []error) {
	return this.Output, this.RenderErrors
}
