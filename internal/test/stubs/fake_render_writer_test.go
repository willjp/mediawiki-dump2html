package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/willjp/mediawiki-dump2html/internal/elements/mwdump"
	"github.com/willjp/mediawiki-dump2html/internal/interfaces"
)

func samplePages() []mwdump.Page {
	page := mwdump.Page{
		Title: "file",
		Revisions: []mwdump.Revision{
			{Text: "== My New Header ==", Timestamp: time.Now()},
		},
	}
	return []mwdump.Page{page}
}

func TestFakeRenderWriterImplementsInterface(t *testing.T) {
	implementsInterface := func(interfaces.RenderWriter) bool {
		return true
	}
	assert.True(t, implementsInterface(&FakeRenderWriter{}))
}

func TestFakeRenderWriterDumpAll(t *testing.T) {
	t.Run("Sets DumpAllCalled", func(t *testing.T) {
		writer := FakeRenderWriter{}
		dump := mwdump.XMLDump{Pages: samplePages()}
		renderer := FakeRenderer{}

		errs := writer.DumpAll(&renderer, &dump, "/var/tmp")
		assert.Nil(t, errs)
		assert.True(t, writer.DumpAllCalled)
	})

	t.Run("Returns DumpAllErrors when provided", func(t *testing.T) {
		writer := FakeRenderWriter{DumpAllErrors: []error{ExpectedError}}
		dump := mwdump.XMLDump{Pages: samplePages()}
		renderer := FakeRenderer{}

		errs := writer.DumpAll(&renderer, &dump, "/var/tmp")
		assert.Equal(t, 1, len(errs))
		assert.Error(t, ExpectedError, errs[0])
	})
}

func TestFakeRenderWriterDump(t *testing.T) {
	t.Run("Sets DumpCalled", func(t *testing.T) {
		writer := FakeRenderWriter{}
		page := samplePages()[0]
		renderer := FakeRenderer{}

		errs := writer.Dump(&renderer, &page, "/var/tmp")
		assert.Nil(t, errs)
		assert.True(t, writer.DumpCalled)
	})

	t.Run("Returns DumpErrors", func(t *testing.T) {
		writer := FakeRenderWriter{DumpErrors: []error{ExpectedError}}
		page := samplePages()[0]
		renderer := FakeRenderer{}

		errs := writer.Dump(&renderer, &page, "/var/tmp")
		assert.Equal(t, 1, len(errs))
		assert.Error(t, ExpectedError, errs[0])
	})
}
