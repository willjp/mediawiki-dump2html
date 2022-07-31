package writers

import (
	"path"
	"sync"

	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

type PageConsumer struct {
	OutDir     string
	Renderer   interfaces.Renderer
	PageDumper interfaces.PageDumper
	PageCh     <-chan mwdump.Page
	ErrorCh    chan<- error
	WaitGroup  *sync.WaitGroup
}

func (this *PageConsumer) Start() {
	var errs []error
	for page := range this.PageCh {
		outPath := path.Join(this.OutDir, this.Renderer.Filename(page.Title))
		errs = this.PageDumper.Dump(this.Renderer, &page, outPath)
		for _, err := range errs {
			this.ErrorCh <- err
		}
	}
	this.WaitGroup.Done()
}
