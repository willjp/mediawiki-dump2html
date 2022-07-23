package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

func TestBuildImplementsInterface(t *testing.T) {
	implementsInterface := func(interfaces.CliCommand) bool {
		return true
	}
	assert.True(t, implementsInterface(&Build{}))
}
