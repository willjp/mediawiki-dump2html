package writers

import (
	"sync"

	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
)

type PageProducer struct {
	XMLDump   *mwdump.XMLDump
	PageCh    chan<- mwdump.Page
	WaitGroup *sync.WaitGroup
}

func (this *PageProducer) Start() {
	defer close(this.PageCh)
	for _, page := range this.XMLDump.Pages {
		this.PageCh <- page
	}
}
