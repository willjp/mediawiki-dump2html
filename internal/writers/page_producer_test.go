package writers

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
)

func TestPageProducer(t *testing.T) {
	samplePage := func(title string) mwdump.Page {
		return mwdump.Page{
			Title: title,
			Revision: []mwdump.Revision{
				{Text: "== My New Header ==", Timestamp: time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)},
			},
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	pageCh := make(chan mwdump.Page, 2)
	dump := mwdump.XMLDump{Pages: []mwdump.Page{samplePage("one"), samplePage("two")}}
	producer := PageProducer{
		XMLDump:   &dump,
		PageCh:    pageCh,
		WaitGroup: &wg,
	}
	producer.Start()

	assert.Equal(t, 2, len(pageCh))
	p1 := <-pageCh
	p2 := <-pageCh
	assert.Equal(t, "one", p1.Title)
	assert.Equal(t, "two", p2.Title)
}
