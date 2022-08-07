package renderers

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"
	"github.com/willjp/mediawiki-dump2html/internal/appfs"
	"github.com/willjp/mediawiki-dump2html/internal/interfaces"
	"github.com/willjp/mediawiki-dump2html/internal/pandoc"
	"github.com/willjp/mediawiki-dump2html/internal/utils"
)

// Struct with methods to render pandoc's sytnaxhighlighting CSS stylesheet.
//
// Pandoc's style.html templating disables syntaxhighlighting unless it is required.
// This object generates a minimal amount of stub files for pandoc to render just the
// CSS for it's syntaxhighlighting.
//
// Roughly equivalent to
//  echo '```html\n<p>placeholder</p>\n```' | pandoc --template=/tmp/dir/highlight-css.template
type HighlightCSS struct {
	pandocExecutor interfaces.PandocExecutor
	tempfileWriter interfaces.FileWriter
}

func NewHighlightCSS(pandocExecutor interfaces.PandocExecutor) HighlightCSS {
	return HighlightCSS{
		pandocExecutor: pandocExecutor,
		tempfileWriter: &utils.FileWriter{},
	}
}

// Renders pandoc's syntaxhighlighting CSS stylesheet
func (this *HighlightCSS) Render() (render string, errs []error) {
	var err error

	// dummy markdown file we'll provide to pandoc to create syntaxhighlighting css for
	htmlSeed := "```html\n<p>placeholder</p>\n```"

	// Dummy pandoc template, that sets '$highlighting-css$'.
	// Pandoc's 'style.html's preprocessor checks for this to determine if highlighting-css should be generated
	file, err := this.writeCssTemplate()
	if err != nil {
		panic(err)
	}
	defer func(file afero.File) {
		err = file.Close()
		if err != nil {
			errs = append(errs, err)
		}
		err = appfs.AppFs.Remove(file.Name())
		if err != nil {
			errs = append(errs, err)
		}
	}(file)

	// Create command to render template
	opts := pandoc.Opts{Template: file.Name()}
	cmd := opts.Command()
	stdin := strings.NewReader(htmlSeed)
	css, errs := this.pandocExecutor.Execute(&cmd, stdin)
	if errs != nil {
		return "", errs
	}

	// cat htmlSeed.md | pandoc --template={dummy-template}  # output is css
	return css, nil
}

// Writes a temporary pandoc template that will enable the generation of syntaxhighlighting-css.
// Caller is responsible for closing and removing the file.
func (this *HighlightCSS) writeCssTemplate() (file afero.File, err error) {
	Os := afero.Afero{Fs: appfs.AppFs}
	dirpath, err := Os.TempDir(os.TempDir(), "highlight-css-template-")
	if err != nil {
		return nil, err
	}

	// pandoc requires a '.template' suffix
	filepath := path.Join(dirpath, "highlight-css.template")
	file, err = Os.Create(filepath)
	if err != nil {
		return file, err
	}
	_, err = this.tempfileWriter.Write(file, []byte("$highlighting-css$\n"))
	return file, err
}
