package writers

import (
	"sync"

	"github.com/willjp/mediawiki-dump2html/internal/elements/mwdump"
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
