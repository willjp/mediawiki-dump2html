package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
)

func TestFakeRendererFilename(t *testing.T) {
	renderer := FakeRenderer{}
	assert.Equal(t, "filename", renderer.Filename(""))
}

func TestFakeRendererSetup(t *testing.T) {
	t.Run("Returns SetupErrors if provided", func(t *testing.T) {
		var ExpectedError = errors.New("Expected")
		renderer := FakeRenderer{SetupErrors: []error{ExpectedError}}
		dump := mwdump.XMLDump{}
		errs := renderer.Setup(&dump, "/var/tmp")
		assert.Equal(t, 1, len(errs))
		assert.Error(t, ExpectedError, errs[0])
	})
}

func TestFakeRendererRender(t *testing.T) {
	t.Run("Returns Output if provided", func(t *testing.T) {
		renderer := FakeRenderer{Output: "abc"}
		page := mwdump.Page{}
		out, errs := renderer.Render(&page)
		assert.Nil(t, errs)
		assert.Equal(t, "abc", out)
	})

	t.Run("Returns RenderErrors if provided", func(t *testing.T) {
		var ExpectedError = errors.New("Expected")
		renderer := FakeRenderer{RenderErrors: []error{ExpectedError}}
		page := mwdump.Page{}
		_, errs := renderer.Render(&page)
		assert.Equal(t, 1, len(errs))
		assert.Error(t, ExpectedError, errs[0])
	})
}
