package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

func TestFakeRenderWriterImplementsInterface(t *testing.T) {
	implementsInterface := func(interfaces.RenderWriter) bool {
		return true
	}
	assert.True(t, implementsInterface(&FakeRenderWriter{}))
}
