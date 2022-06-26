package elements

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

// Represents a mediawiki <page> element
type Page struct {
	Title    string     `xml:"title"`
	Revision []Revision `xml:"revision"`
}

func (page *Page) LatestRevision() Revision {
	return page.Revision[len(page.Revision)-1]
}

func (page *Page) WriteRst(sphinxRoot string) {
	var fileModified time.Time
	fileName := fmt.Sprint(page.Title, ".rst")
	rstPath := path.Join(sphinxRoot, fileName)
	rstStat, err := os.Stat(rstPath)
	switch {
	case err == nil:
		fileModified = rstStat.ModTime()
	case errors.Is(err, fs.ErrNotExist):
		fileModified = time.Unix(0, 0)
	default:
		panic(err)
	}

	revision := page.LatestRevision()
	if revision.Timestamp.After(fileModified) {
		file, err := os.Create(rstPath)
		panicOn(err)

		_, err = file.WriteString(page.renderRst())
		panicOn(err)
	}
}

func (page *Page) renderRst() string {
	directives := `
	.. role:: raw-html(raw)
	  :format: html
	`

	// page title between '='s
	titleLen := len([]rune(page.Title))
	title := fmt.Sprint(
		strings.Repeat("=", titleLen),
		page.Title,
		strings.Repeat("=", titleLen),
	)

	// raw=$(cat $PAGE | pandoc -f mediawiki -t rst)
	cmd := exec.Command("pandoc", "-f", "mediawiki", "-t", "rst")
	writer, err := cmd.StdinPipe()
	panicOn(err)
	_, err = writer.Write([]byte(page.LatestRevision().Text))
	panicOn(err)
	raw, err := cmd.Output()
	panicOn(err)

	// replace '<br>' with something rst understands
	rendered := strings.ReplaceAll(string(raw), "<br>", ":raw-html:`<br/>`")

	return fmt.Sprintf(directives, title, rendered)
}
