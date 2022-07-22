package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

func TestFakeFile(t *testing.T) {
	t.Run("Implements interfaces.osFile", func(t *testing.T) {
		var implementsIface = func(iface interfaces.OsFile) bool {
			return true
		}
		assert.True(t, implementsIface(&os.File{}))
	})
}
