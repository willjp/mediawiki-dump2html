package test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/pandoc"
)

func TestFakePandocExecutorExecute(t *testing.T) {
	t.Run("Records provided stdin", func(t *testing.T) {
		stdin := strings.NewReader("abc")
		cmd := pandoc.Cmd{}
		executor := FakePandocExecutor{}
		executor.Execute(&cmd, stdin)

		res, err := io.ReadAll(executor.Stdin)
		assert.Nil(t, err)
		assert.Equal(t, []byte("abc"), res)
	})

	t.Run("Returns provided errors", func(t *testing.T) {
		var ExpectedError = errors.New("Expected")
		stdin := strings.NewReader("")
		cmd := pandoc.Cmd{}
		executor := FakePandocExecutor{Errors: []error{ExpectedError}}
		out, errs := executor.Execute(&cmd, stdin)

		assert.Equal(t, "", out)
		assert.Equal(t, []error{ExpectedError}, errs)
	})

	t.Run("Returns provided result", func(t *testing.T) {
		stdin := strings.NewReader("")
		cmd := pandoc.Cmd{}
		executor := FakePandocExecutor{Render: "Hole is good music"}
		out, errs := executor.Execute(&cmd, stdin)

		assert.Nil(t, errs)
		assert.Equal(t, "Hole is good music", out)
	})
}

func TestFakePandocExecutorArgs(t *testing.T) {
	t.Run("Returns Args when provided", func(t *testing.T) {
		args := []string{"abc", "def"}
		executor := FakePandocExecutor{args: args}
		assert.Equal(t, args, executor.Args())
	})
}
