package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

func TestFakeCmd(t *testing.T) {
	t.Run("Implements CmdExecutor", func(t *testing.T) {
		var implementsIface = func(iface interfaces.CmdExecutor) bool {
			return true
		}
		cmd := FakeCmd{}
		assert.True(t, implementsIface(&cmd))
	})
}
