package renderers

import (
	"io"
	"path"
	"regexp"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/appfs"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
	test "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
)

var htmlWhitespaceRx = regexp.MustCompile(`(?m)(^\s+|\n)`)

func TestHTMLNew(t *testing.T) {
	html := New()
	assert.NotNil(t, html.pandocExecutor)
}

func TestHTMLFilename(t *testing.T) {
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

func TestHTMLSetup(t *testing.T) {
	appfs.AppFs = afero.NewMemMapFs()
	Os := afero.Afero{Fs: appfs.AppFs}
	t.Run("Writes Expected File", func(t *testing.T) {
		filepath := path.Join("/var/tmp", stylesheetName)
		pages := []mwdump.Page{
			{
				Title: "main_page",
				Revision: []mwdump.Revision{
					{
						Text:      "== My New Header ==",
						Timestamp: time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
		}
		dump := mwdump.XMLDump{Pages: pages}
		executor := test.FakePandocExecutor{}

		renderer := newHTML(&executor)
		errs := renderer.Setup(&dump, "/var/tmp")
		assert.Nil(t, errs)
		exists, err := Os.Exists(filepath)
		assert.Nil(t, err)
		assert.True(t, exists)
	})
}

func TestHTMLRender(t *testing.T) {
	t.Run("Commandline Arguments set correctly", func(t *testing.T) {
		executor := test.FakePandocExecutor{}
		renderer := newHTML(&executor)
		page := mwdump.Page{
			Title: "Main Page",
			Revision: []mwdump.Revision{
				{
					Text:      "== My New Header ==",
					Timestamp: time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
			},
		}

		_, err := renderer.Render(&page)
		expects := []string{"pandoc", "-f", "mediawiki", "-t", "html"}
		assert.Nil(t, err)
		assert.Equal(t, expects, executor.Args())
	})

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
		executor := test.FakePandocExecutor{}
		renderer := newHTML(&executor)
		renderer.Render(&page)
		out, err := io.ReadAll(executor.Stdin)
		assert.Nil(t, err)
		assert.Equal(t, []byte("== My New Header =="), []byte(out))
	})

	tcases := []struct {
		name    string
		render  string
		expects string
	}{
		{
			name:   "Inserts Page Info",
			render: "",
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
			name: "Increments all headers in page",
			render: htmlWhitespaceRx.ReplaceAllString(`
				<h1>H1</h1>
				<h2>H2</h2>
				<h3>H3</h3>
				<h4>H4</h4>
				<h5>H5</h5>
				<h6>H6</h6>
			`, ""),
			expects: htmlWhitespaceRx.ReplaceAllString(`
				<html>
				  <head>
				    <title>Main Page</title>
				    <link rel="stylesheet" href="style.css"/>
				  </head>
				  <body>
				    <h1 id="main_page">Main Page</h1>
				    <h2>H1</h2>
				    <h3>H2</h3>
				    <h4>H3</h4>
				    <h5>H4</h5>
				    <h6>H5</h6>
				    <h6>H6</h6>
				  </body>
				</html>
			`, ""),
		},
		{
			name:   "Adds .html suffix to link URLs",
			render: `<a href="another_page">Another Page</a>`,
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
			name:   "Sanitizes link urls to match files written to filesystem",
			render: `<a href="Programming: Concepts">Programming: Concepts</a>`,
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
			executor := test.FakePandocExecutor{Render: tcase.render}
			renderer := newHTML(&executor)
			render, errs := renderer.Render(&page)

			assert.Nil(t, errs)
			assert.Equal(t, tcase.expects, render)
		})
	}
}
