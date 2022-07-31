package writers

import (
	"runtime"
	"sync"

	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

type AsyncRenderWriter struct {
	NumWorkers int
	pageDumper interfaces.PageDumper
}

func NewAsyncRenderWriter(numWorkers int) *AsyncRenderWriter {
	return &AsyncRenderWriter{
		NumWorkers: numWorkers,
		pageDumper: &RenderWriter{},
	}
}

func (this *AsyncRenderWriter) GetNumWorkers() int {
	if this.NumWorkers <= 0 {
		return runtime.NumCPU()
	} else {
		return this.NumWorkers
	}
}

func (this *AsyncRenderWriter) DumpAll(renderer interfaces.Renderer, dump *mwdump.XMLDump, outDir string) (errs []error) {
	errs = renderer.Setup(dump, outDir)

	numWorkers := this.GetNumWorkers()
	pageCh := make(chan mwdump.Page, 100)
	errorCh := make(chan error, 50*numWorkers)
	consumers := make([]*PageConsumer, 0, numWorkers)
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	// build producer/consumers
	producer := PageProducer{
		XMLDump:   dump,
		PageCh:    pageCh,
		WaitGroup: &wg,
	}
	for i := 0; i < numWorkers; i++ {
		consumer := PageConsumer{
			PageCh:     pageCh,
			ErrorCh:    errorCh,
			Renderer:   renderer,
			PageDumper: this.pageDumper,
			OutDir:     outDir,
			WaitGroup:  &wg,
		}
		consumers = append(consumers, &consumer)
	}

	// start producer/consumers
	go producer.Start()
	for _, consumer := range consumers {
		go consumer.Start()
	}

	// wait for workers to complete.
	// nothing else will write to 'errorCh'
	wg.Wait()
	close(errorCh)
	if len(errorCh) == 0 {
		return nil
	}
	for err := range errorCh {
		errs = append(errs, err)
	}
	return errs
}
