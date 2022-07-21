package pandoc

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	stubs "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
)

// tests
func TestExecute(t *testing.T) {
	t.Run("Forwards STDIN to cmd.Execute()", func(t *testing.T) {
		stdin := strings.NewReader("abc")
		cmd := stubs.FakeCmd{}
		executor := Executor{}
		_, errs := executor.Execute(&cmd, stdin)
		assert.Nil(t, errs)

		in, err := io.ReadAll(cmd.Stdin)
		assert.Nil(t, err)
		assert.Equal(t, []byte("abc"), in)
	})

	t.Run("Performs/Wraps cmd.Execute()", func(t *testing.T) {
		stdin := strings.NewReader("abc")
		cmd := stubs.FakeCmd{Render: "abc"}
		executor := Executor{}
		out, errs := executor.Execute(&cmd, stdin)
		assert.Nil(t, errs)
		assert.Equal(t, "abc", out)
	})
}

func TestArgs(t *testing.T) {
	t.Run("Returns contained cmd.Args()", func(t *testing.T) {
		stdin := strings.NewReader("")
		args := []string{"foo", "-a", "--verbose"}
		cmd := stubs.FakeCmd{CliArgs: args}
		executor := Executor{}
		_, errs := executor.Execute(&cmd, stdin)
		assert.Nil(t, errs)
		assert.Equal(t, args, executor.Args())
	})
}
