package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	test "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
)

// Stub out the pandoc process.
//   * Stdin is readable (we can test what it receives)
//   * Each test receives fake stdout/stderr results
//
// based on this fake pandoc output, assert method behaviour.
func TestPandocCmd(t *testing.T) {
	var stdin strings.Builder
	html := "<html><h2>foo</h2></html>"
	cmd := test.FakeCmd{
		Stdin:  test.FakeWriteCloser{Writer: stdin},
		Stderr: test.FakeReadCloser{Reader: strings.NewReader("")},
		Stdout: test.FakeReadCloser{Reader: strings.NewReader(html)},
		Args:   []string{"pandoc", "-f", "mediawiki", "-t", "html"},
	}
	pcmd := PandocCmd{Cmd: cmd}
	result, errs := pcmd.Execute(strings.NewReader("== My Header =="))
	assert.Nil(t, errs)
	assert.Equal(t, html, result)
}
