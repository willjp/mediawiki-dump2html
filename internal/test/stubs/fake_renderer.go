package test

import "willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"

type FakeRenderer struct {
	Output       string
	SetupErrors  []error
	RenderErrors []error
}

func (this *FakeRenderer) Filename(pageTitle string) string {
	return "filename"
}

func (this *FakeRenderer) Setup(dump *mwdump.XMLDump, outDir string) []error {
	return this.SetupErrors
}

func (this *FakeRenderer) Render(page *mwdump.Page) (render string, errs []error) {
	return this.Output, this.RenderErrors
}
