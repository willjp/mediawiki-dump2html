package renderers

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	test "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
	pandoc "willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils/pandoc"
)

var htmlWhitespaceRx = regexp.MustCompile(`(?m)(^\s+|\n)`)

func TestFilename(t *testing.T) {
	tcases := []struct {
		name      string
		pageTitle string
		expects   string
	}{
		{
			name:      "URL Valid PageName",
			pageTitle: "mainpage",
			expects:   "mainpage.html",
		},
		{
			name:      "URL InValid PageName",
			pageTitle: "Programming: Main Page",
			expects:   "programming__main_page.html",
		},
	}
	for _, tcase := range tcases {
		t.Run(tcase.name, func(t *testing.T) {
			renderer := HTML{}
			res := renderer.Filename(tcase.pageTitle)
			assert.Equal(t, tcase.expects, res)
		})
	}
}

func TestRenderCmd(t *testing.T) {
	renderer := HTML{}
	cmd := renderer.RenderCmd()
	expects := []string{"pandoc", "-f", "mediawiki", "-t", "html"}
	assert.Equal(t, expects, cmd.Args)
}

func TestRenderExec(t *testing.T) {
	t.Run("Renders Latest Page Revision", func(t *testing.T) {
		page := mwdump.Page{
			Title: "Main Page",
			Revision: []mwdump.Revision{
				{
					Text:      "== My New Header ==",
					Timestamp: time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
			},
		}
		stdin := strings.Builder{}
		stdin.Write([]byte(""))
		cmd := test.FakeCmd{
			Stdin:  &test.FakeWriteCloser{Writer: &stdin},
			Stderr: &test.FakeReadCloser{Reader: strings.NewReader("")},
			Stdout: &test.FakeReadCloser{Reader: strings.NewReader("<html><h2>foo</h2></html>")},
		}
		pcmd := pandoc.Cmd{Cmd: cmd}
		renderer := HTML{}
		renderer.RenderExec(&pcmd, &page)
		assert.Equal(t, "== My New Header ==", stdin.String())
	})

	tcases := []struct {
		name         string
		pandocRender string
		expects      string
	}{
		{
			name:         "Inserts Page Info",
			pandocRender: "",
			expects: htmlWhitespaceRx.ReplaceAllString(`
				<html>
				  <head>
				    <title>Main Page</title>
				    <link rel="stylesheet" href="style.css"/>
				  </head>
				  <body>
				    <h1 id="main_page">Main Page</h1>
				  </body>
				</html>
			`, ""),
		},
		{
			name:         "Increments all headers in page",
			pandocRender: "<h1>Documentation</h1>",
			expects: htmlWhitespaceRx.ReplaceAllString(`
				<html>
				  <head>
				    <title>Main Page</title>
				    <link rel="stylesheet" href="style.css"/>
				  </head>
				  <body>
				    <h1 id="main_page">Main Page</h1>
				    <h2>Documentation</h2>
				  </body>
				</html>
			`, ""),
		},
		{
			name:         "Adds .html suffix to link URLs",
			pandocRender: `<a href="another_page">Another Page</a>`,
			expects: htmlWhitespaceRx.ReplaceAllString(`
				<html>
				  <head>
				    <title>Main Page</title>
				    <link rel="stylesheet" href="style.css"/>
				  </head>
				  <body>
				    <h1 id="main_page">Main Page</h1>
				    <a href="another_page.html">Another Page</a>
				  </body>
				</html>
			`, ""),
		},
		{
			name:         "Sanitizes link urls to match files written to filesystem",
			pandocRender: `<a href="Programming: Concepts">Programming: Concepts</a>`,
			expects: htmlWhitespaceRx.ReplaceAllString(`
				<html>
				  <head>
				    <title>Main Page</title>
				    <link rel="stylesheet" href="style.css"/>
				  </head>
				  <body>
				    <h1 id="main_page">Main Page</h1>
				    <a href="programming__concepts.html">Programming: Concepts</a>
				  </body>
				</html>
			`, ""),
		},
	}
	for _, tcase := range tcases {
		t.Run(tcase.name, func(t *testing.T) {
			page := mwdump.Page{
				Title:    "Main Page",
				Revision: []mwdump.Revision{{}},
			}
			stdin := strings.Builder{}
			stdin.Write([]byte(""))
			cmd := test.FakeCmd{
				Stdin:  &test.FakeWriteCloser{Writer: &stdin},
				Stderr: &test.FakeReadCloser{Reader: strings.NewReader("")},
				Stdout: &test.FakeReadCloser{Reader: strings.NewReader(tcase.pandocRender)},
			}
			pcmd := pandoc.Cmd{Cmd: cmd}
			renderer := HTML{}
			render, errs := renderer.RenderExec(&pcmd, &page)

			assert.Nil(t, errs)
			assert.Equal(t, tcase.expects, render)
		})
	}
}
