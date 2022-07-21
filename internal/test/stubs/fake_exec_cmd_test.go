package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

func TestFakeExecCmd(t *testing.T) {
	t.Run("Implements ExecCmd", func(t *testing.T) {
		var implementsIface = func(iface interfaces.ExecCmd) bool {
			return true
		}
		cmd := FakeExecCmd{}
		assert.True(t, implementsIface(cmd))
	})
}
