package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

func TestFakeCmd(t *testing.T) {
	t.Run("Implements CmdExecutor", func(t *testing.T) {
		var implementsIface = func(iface interfaces.ExecCmdExecutor) bool {
			return true
		}
		cmd := FakeCmd{}
		assert.True(t, implementsIface(&cmd))
	})

	t.Run("Returns FakeCmd.Render when provided", func(t *testing.T) {
		cmd := FakeCmd{Render: "abc"}
		stdin := strings.NewReader("")
		res, errs := cmd.Execute(stdin)
		assert.Nil(t, errs)
		assert.Equal(t, "abc", res)
	})

	t.Run("Returns FakeCmd.Errors when provided", func(t *testing.T) {
		var ExpectedError = errors.New("Expected")
		cmd := FakeCmd{Errors: []error{ExpectedError}}
		stdin := strings.NewReader("")
		res, errs := cmd.Execute(stdin)
		assert.Equal(t, 1, len(errs))
		assert.Error(t, ExpectedError, errs[0])
		assert.Equal(t, "", res)
	})
}
