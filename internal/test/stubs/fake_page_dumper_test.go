package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakePageDumperDump(t *testing.T) {
	t.Run("Sets DumpCalled", func(t *testing.T) {
		writer := FakePageDumper{}
		page := samplePages()[0]
		renderer := FakeRenderer{}

		errs := writer.Dump(&renderer, &page, "/var/tmp")
		assert.Nil(t, errs)
		assert.True(t, writer.DumpCalled)
	})

	t.Run("Returns DumpErrors", func(t *testing.T) {
		writer := FakePageDumper{DumpErrors: []error{ExpectedError}}
		page := samplePages()[0]
		renderer := FakeRenderer{}

		errs := writer.Dump(&renderer, &page, "/var/tmp")
		assert.Equal(t, 1, len(errs))
		assert.Error(t, ExpectedError, errs[0])
	})

	t.Run("Sets DumpArgs", func(t *testing.T) {
		writer := FakePageDumper{}
		page := samplePages()[0]
		renderer := FakeRenderer{}

		errs := writer.Dump(&renderer, &page, "/var/tmp")
		assert.Nil(t, errs)
		assert.Equal(t, []any{&renderer, &page, "/var/tmp"}, writer.DumpArgs)
	})
}
