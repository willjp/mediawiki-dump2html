package test

import (
	"strings"
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

	t.Run("StdinPipe returns provided Stdin", func(t *testing.T) {
		writer := strings.Builder{}
		stdin := FakeWriteCloser{Writer: &writer}
		cmd := FakeExecCmd{Stdin: &stdin}
		result, err := cmd.StdinPipe()
		assert.Nil(t, err)
		assert.Equal(t, &stdin, result)
	})

	t.Run("StdoutPipe returns provided Stdout", func(t *testing.T) {
		stdout := FakeReadCloser{}
		cmd := FakeExecCmd{Stdout: &stdout}
		result, err := cmd.StdoutPipe()
		assert.Nil(t, err)
		assert.Equal(t, &stdout, result)
	})

	t.Run("StderrPipe returns provided Stderr", func(t *testing.T) {
		stderr := FakeReadCloser{}
		cmd := FakeExecCmd{Stderr: &stderr}
		result, err := cmd.StderrPipe()
		assert.Nil(t, err)
		assert.Equal(t, &stderr, result)
	})

	t.Run("Start returns nil", func(t *testing.T) {
		cmd := FakeExecCmd{}
		assert.Nil(t, cmd.Start())
	})

	t.Run("Wait returns nil", func(t *testing.T) {
		cmd := FakeExecCmd{}
		assert.Nil(t, cmd.Wait())
	})
}
