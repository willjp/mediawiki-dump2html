package renderers

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	test "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
)

func TestCSSWriteCssFile(t *testing.T) {
	t.Run("Writes extracted CSS to file", func(t *testing.T) {
		executor := test.FakePandocExecutor{}
		renderer := NewCSS(&executor)
		pages := []mwdump.Page{
			{
				Title: "Main Page",
				Revision: []mwdump.Revision{
					{
						Text:      "== My New Header ==",
						Timestamp: time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
		}
		dump := mwdump.XMLDump{Pages: pages}
		renderer.WriteCssFile(&dump, "/var/tmp/foo.css")
	})
}

func TestCSSRender(t *testing.T) {
	t.Run("Commandline Arguments set correctly", func(t *testing.T) {
		executor := test.FakePandocExecutor{}
		renderer := NewCSS(&executor)
		pages := []mwdump.Page{
			{
				Title: "Main Page",
				Revision: []mwdump.Revision{
					{
						Text:      "== My New Header ==",
						Timestamp: time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
		}
		dump := mwdump.XMLDump{Pages: pages}
		_, err := renderer.Render(&dump)
		expects := []string{"pandoc", "-f", "mediawiki", "-t", "html", "--standalone"}
		assert.Nil(t, err)
		assert.Equal(t, expects, executor.Args())
	})

	t.Run("Renders first page from dump", func(t *testing.T) {
		executor := test.FakePandocExecutor{}
		renderer := NewCSS(&executor)
		pages := []mwdump.Page{
			{
				Title: "First Page",
				Revision: []mwdump.Revision{
					{
						Text:      "== First Page ==",
						Timestamp: time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			{
				Title: "Second Page",
				Revision: []mwdump.Revision{
					{
						Text:      "== Second Page ==",
						Timestamp: time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
		}
		dump := mwdump.XMLDump{Pages: pages}
		_, errs := renderer.Render(&dump)
		assert.Nil(t, errs)
		out, err := io.ReadAll(executor.Stdin)
		assert.Nil(t, err)
		assert.Equal(t, []byte("== First Page =="), []byte(out))
	})

	t.Run("Extracts CSS from pandoc's generated HTML file", func(t *testing.T) {
		pandocCss := htmlWhitespaceRx.ReplaceAllString(`
			html {
			  line-height: 1.5;
			  font-family: Georgia, serif;
			  font-size: 20px;
			  color: #1a1a1a;
			  background-color: #fdfdfd;
			}
			body {
			  margin: 0 auto;
			  max-width: 36em;
			  padding-left: 50px;
			  padding-right: 50px;
			  padding-top: 50px;
			  padding-bottom: 50px;
			  hyphens: auto;
			  overflow-wrap: break-word;
			  text-rendering: optimizeLegibility;
			  font-kerning: normal;
			}
		`, "")
		pandocHtml := fmt.Sprintf(htmlWhitespaceRx.ReplaceAllString(`
			<!DOCTYPE html>
			<html xmlns="http://www.w3.org/1999/xhtml" lang="" xml:lang="">
			<head>
			  <meta charset="utf-8" />
			  <meta name="generator" content="pandoc" />
			  <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=yes" />
			  <title>-</title>
			  <style>
			  %s
			  </style>
			</head>
			<body>
			Hello world
			</body>
			</html>
		`, ""), pandocCss)
		executor := test.FakePandocExecutor{Render: pandocHtml}
		renderer := NewCSS(&executor)
		pages := []mwdump.Page{
			{
				Title: "First Page",
				Revision: []mwdump.Revision{
					{
						Text:      "== First Page ==",
						Timestamp: time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
		}
		dump := mwdump.XMLDump{Pages: pages}
		render, errs := renderer.Render(&dump)
		assert.Nil(t, errs)
		assert.Equal(t, pandocCss, render)
	})
}
