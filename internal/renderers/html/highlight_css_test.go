package renderers

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/willjp/mediawiki-dump2html/internal/appfs"
	test "github.com/willjp/mediawiki-dump2html/internal/test/stubs"
)

var ExpectedError = errors.New("Expected")

func TestNewHighlightCSS(t *testing.T) {
	executor := test.FakePandocExecutor{}
	css := NewHighlightCSS(&executor)
	assert.NotNil(t, css.pandocExecutor)
}

func TestHighlightCSSRender(t *testing.T) {
	setup := func(t *testing.T) {
		appfs.AppFs = afero.NewMemMapFs()
		Os := afero.Afero{Fs: appfs.AppFs}
		err := Os.MkdirAll(os.TempDir(), 0755)
		assert.Nil(t, err)
	}

	t.Run("htmlSeed is written to pandoc STDIN", func(t *testing.T) {
		setup(t)
		executor := test.FakePandocExecutor{}
		css := NewHighlightCSS(&executor)

		_, errs := css.Render()
		assert.Nil(t, errs)

		expects := "```html\n<p>placeholder</p>\n```"
		out := make([]byte, len(expects))
		_, err := executor.Stdin.Read(out)
		assert.Nil(t, err)
		assert.Equal(t, []byte(expects), out)
	})

	t.Run("$highlighting-css$ is written to pandoc template tempfile", func(t *testing.T) {
		setup(t)
		written := bytes.NewBuffer(nil)
		executor := test.FakePandocExecutor{}
		css := HighlightCSS{
			pandocExecutor: &executor,
			tempfileWriter: &test.FakeFileWriter{Written: written},
		}

		_, errs := css.Render()
		assert.Nil(t, errs)

		out, err := io.ReadAll(written)
		assert.Nil(t, err)
		assert.Equal(t, []byte("$highlighting-css$\n"), out)
	})

	t.Run("pandoc template tempfile is passed as cli argument to pandoc", func(t *testing.T) {
		setup(t)
		written := bytes.NewBuffer(nil)
		executor := test.FakePandocExecutor{}
		css := HighlightCSS{
			pandocExecutor: &executor,
			tempfileWriter: &test.FakeFileWriter{Written: written},
		}

		_, errs := css.Render()
		assert.Nil(t, errs)

		args := executor.Args()
		paramMatch := regexp.MustCompile(`^--template=([a-zA-Z]:|/)`)
		assert.True(t, paramMatch.Match([]byte(args[1])))
	})

	t.Run("returns css on render success", func(t *testing.T) {
		setup(t)
		returnCss := `
			code span.vs { color: #4070a0; } /* VerbatimString */
			code span.wa { color: #60a0b0; font-weight: bold; font-style: italic; } /* Warning */
			`
		executor := test.FakePandocExecutor{Render: returnCss}
		css := NewHighlightCSS(&executor)

		out, errs := css.Render()
		assert.Nil(t, errs)
		assert.Equal(t, returnCss, out)
	})

	t.Run("returns errors on render failure", func(t *testing.T) {
		setup(t)
		executor := test.FakePandocExecutor{Errors: []error{ExpectedError}}
		css := NewHighlightCSS(&executor)

		_, errs := css.Render()
		assert.Equal(t, 1, len(errs))
		assert.Error(t, ExpectedError, errs[0])
	})
}
