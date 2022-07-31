package writers

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	test "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
)

func TestAsyncRenderWriterGetNumWorkers(t *testing.T) {
	t.Run("Defaults to num CPUs", func(t *testing.T) {
		writer := AsyncRenderWriter{}
		num := writer.GetNumWorkers()
		assert.Equal(t, runtime.NumCPU(), num)
	})

	t.Run("Respects NumWorkers field", func(t *testing.T) {
		writer := AsyncRenderWriter{NumWorkers: 2}
		num := writer.GetNumWorkers()
		assert.Equal(t, 2, num)
	})
}

func TestAsyncRenderWriterDumpAll(t *testing.T) {
	t.Run("Renders a Page", func(t *testing.T) {
		renderer := test.FakeRenderer{}
		dumper := test.FakePageDumper{}
		dump := mwdump.XMLDump{
			Pages: []mwdump.Page{
				{Title: "file", Revisions: []mwdump.Revision{{Text: "== My Header =="}}},
			},
		}
		writer := AsyncRenderWriter{
			NumWorkers: 1,
			pageDumper: &dumper,
		}
		errs := writer.DumpAll(&renderer, &dump, "/var/tmp")
		assert.Nil(t, errs)
		assert.True(t, dumper.DumpCalled)
	})

	t.Run("Render Errors are returned", func(t *testing.T) {
		renderer := test.FakeRenderer{}
		dumper := test.FakePageDumper{DumpErrors: []error{ExpectedError}}
		dump := mwdump.XMLDump{
			Pages: []mwdump.Page{
				{Title: "file", Revisions: []mwdump.Revision{{Text: "== My Header =="}}},
			},
		}
		writer := AsyncRenderWriter{
			NumWorkers: 1,
			pageDumper: &dumper,
		}
		errs := writer.DumpAll(&renderer, &dump, "/var/tmp")
		assert.Equal(t, 1, len(errs))
		if len(errs) == 1 {
			assert.Error(t, ExpectedError, errs[0])
		}
	})
}
