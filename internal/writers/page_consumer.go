package writers

import (
	"path"
	"sync"

	"github.com/willjp/mediawiki-dump2html/internal/elements/mwdump"
	"github.com/willjp/mediawiki-dump2html/internal/interfaces"
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
