package pandoc

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	test "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
)

func TestExecuteSuccess(t *testing.T) {
	var stdin strings.Builder
	stdin.WriteString("")
	cmd := test.FakeExecCmd{
		Stdin:  &test.FakeWriteCloser{Writer: &stdin},
		Stderr: &test.FakeReadCloser{Reader: strings.NewReader("")},
		Stdout: &test.FakeReadCloser{Reader: strings.NewReader("<html><h2>foo</h2></html>")},
		Args:   []string{"pandoc", "-f", "mediawiki", "-t", "html"},
	}
	pcmd := Cmd{Cmd: cmd}
	result, errs := pcmd.Execute(strings.NewReader("== My Header =="))
	assert.Nil(t, errs)
	assert.Equal(t, "<html><h2>foo</h2></html>", result)
}

func TestExecuteAppliesStdinToProcess(t *testing.T) {
	var stdin strings.Builder
	stdin.WriteString("")
	html := "<html><h2>foo</h2></html>"
	cmd := test.FakeExecCmd{
		Stdin:  &test.FakeWriteCloser{Writer: &stdin},
		Stderr: &test.FakeReadCloser{Reader: strings.NewReader("")},
		Stdout: &test.FakeReadCloser{Reader: strings.NewReader(html)},
		Args:   []string{"pandoc", "-f", "mediawiki", "-t", "html"},
	}
	pcmd := Cmd{Cmd: cmd}
	_, errs := pcmd.Execute(strings.NewReader("== My Header =="))

	assert.Nil(t, errs)
	assert.Equal(t, "== My Header ==", stdin.String())
}

func TestExecuteReturnErrors(t *testing.T) {
	var ExpectedError = errors.New("Expected Test Error")
	tcases := []struct {
		name   string
		stdin  *test.FakeWriteCloser
		stderr *test.FakeReadCloser
		stdout *test.FakeReadCloser
	}{
		// STDOUT/STDERR Close handled by exec.Cmd.Wait()
		{
			name:   "STDIN Close error returned",
			stdin:  &test.FakeWriteCloser{Writer: &strings.Builder{}, CloseError: ExpectedError},
			stdout: &test.FakeReadCloser{Reader: strings.NewReader("")},
			stderr: &test.FakeReadCloser{Reader: strings.NewReader("")},
		},
		{
			name:   "STDIN Write error returned",
			stdin:  &test.FakeWriteCloser{Writer: &strings.Builder{}, WriteError: ExpectedError},
			stdout: &test.FakeReadCloser{Reader: strings.NewReader("")},
			stderr: &test.FakeReadCloser{Reader: strings.NewReader("")},
		},
		{
			name:   "STDOUT Write error returned",
			stdin:  &test.FakeWriteCloser{Writer: &strings.Builder{}},
			stdout: &test.FakeReadCloser{Reader: strings.NewReader(""), ReadError: ExpectedError},
			stderr: &test.FakeReadCloser{Reader: strings.NewReader("")},
		},
		{
			name:   "STDERR Write error returned",
			stdin:  &test.FakeWriteCloser{Writer: &strings.Builder{}},
			stdout: &test.FakeReadCloser{Reader: strings.NewReader("")},
			stderr: &test.FakeReadCloser{Reader: strings.NewReader(""), ReadError: ExpectedError},
		},
	}
	for _, tcase := range tcases {
		t.Run(tcase.name, func(t *testing.T) {
			cmd := test.FakeExecCmd{
				Stdin:  tcase.stdin,
				Stdout: tcase.stdout,
				Stderr: tcase.stderr,
				Args:   []string{"pandoc", "-f", "mediawiki", "-t", "html"},
			}
			pcmd := Cmd{Cmd: cmd}
			_, errs := pcmd.Execute(strings.NewReader("== My Header =="))
			assert.Equal(t, 1, len(errs))
			if len(errs) != 1 {
				return
			}
			assert.Error(t, ExpectedError, errs[0])
		})
	}
}
