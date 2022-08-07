package writers

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/willjp/mediawiki-dump2html/internal/elements/mwdump"
	test "github.com/willjp/mediawiki-dump2html/internal/test/stubs"
)

func TestPageConsumer(t *testing.T) {
	t.Run("Dumps Received Pages with PageDumper", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		pageCh := make(chan mwdump.Page, 1)
		errorCh := make(chan error, 1)
		consumer := PageConsumer{
			OutDir:     "/var/tmp",
			Renderer:   &test.FakeRenderer{},
			PageDumper: &test.FakePageDumper{},
			PageCh:     pageCh,
			ErrorCh:    errorCh,
			WaitGroup:  &wg,
		}
		page := mwdump.Page{
			Title:     "file",
			Revisions: []mwdump.Revision{{Text: "== My New Header ==", Timestamp: time.Now()}},
		}
		pageCh <- page
		close(pageCh)
		consumer.Start()

		castDumper := consumer.PageDumper.(*test.FakePageDumper)
		assert.True(t, castDumper.DumpCalled)
	})

	t.Run("Emits PageDumper errors to ErrorCh", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		pageCh := make(chan mwdump.Page, 1)
		errorCh := make(chan error, 1)
		consumer := PageConsumer{
			OutDir:     "/var/tmp",
			Renderer:   &test.FakeRenderer{},
			PageDumper: &test.FakePageDumper{DumpErrors: []error{ExpectedError}},
			PageCh:     pageCh,
			ErrorCh:    errorCh,
			WaitGroup:  &wg,
		}
		page := mwdump.Page{
			Title:     "file",
			Revisions: []mwdump.Revision{{Text: "== My New Header ==", Timestamp: time.Now()}},
		}
		pageCh <- page
		close(pageCh)
		consumer.Start()

		assert.Equal(t, 1, len(errorCh))
		if len(errorCh) == 1 {
			err := <-errorCh
			assert.Error(t, ExpectedError, err)
		}
	})

	t.Run("Writes to expected file", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		pageCh := make(chan mwdump.Page, 1)
		errorCh := make(chan error, 1)
		consumer := PageConsumer{
			OutDir:     "/var/tmp",
			Renderer:   &test.FakeRenderer{},
			PageDumper: &test.FakePageDumper{},
			PageCh:     pageCh,
			ErrorCh:    errorCh,
			WaitGroup:  &wg,
		}
		page := mwdump.Page{
			Title:     "file",
			Revisions: []mwdump.Revision{{Text: "== My New Header ==", Timestamp: time.Now()}},
		}
		pageCh <- page
		close(pageCh)
		consumer.Start()

		castDumper := consumer.PageDumper.(*test.FakePageDumper)
		assert.True(t, castDumper.DumpCalled)
		assert.Equal(t, "/var/tmp/file", castDumper.DumpArgs[2])
	})
}
