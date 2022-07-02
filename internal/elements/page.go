package elements

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/lithammer/dedent"
	"willpittman.net/x/logger"
)

var invalidPathCh *regexp.Regexp

func init() {
	var err error
	invalidPathCh, err = regexp.Compile("[/]")
	panicOn(err)
}

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

func (page *Page) RelPath() string {
	fileName := fmt.Sprint(page.Title, ".rst")
	sanitized := invalidPathCh.ReplaceAll([]byte(fileName), []byte("-"))
	return string(sanitized)
}

func (page *Page) WriteRst(sphinxRoot string) error {
	rmFileOn := func(file *os.File, err error) {
		if err != nil {
			logger.Errorf("Error encountered, removing: %s", file.Name())
			os.Remove(file.Name())
		}
	}

	var fileModified time.Time
	rstPath := path.Join(sphinxRoot, page.RelPath())
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

		logger.Infof("Writing: %s\n", rstPath)
		rendered, err := page.renderRst()
		if err != nil {
			rmFileOn(file, err)
			return err
		}
		_, err = file.WriteString(rendered)
		if err != nil {
			rmFileOn(file, err)
			return err
		}
	}
	return nil
}

// Converts mediawiki text to rst, with tweaks so it behaves well with sphinx-docs.
func (page *Page) renderRst() (rendered string, err error) {
	directives := dedent.Dedent(`
	.. role:: raw-html(raw)
	  :format: html

	`)

	// page title between '='s
	titleLen := len([]rune(page.Title))
	title := fmt.Sprint(
		strings.Repeat("=", titleLen), "\n",
		page.Title, "\n",
		strings.Repeat("=", titleLen), "\n\n",
	)

	// cat $PAGE | pandoc -f mediawiki -t rst
	pandocRender, err := page.pandocWikiToRst()
	if err != nil {
		return "", err
	}

	// replace '<br>' with something rst understands
	render := strings.ReplaceAll(pandocRender, "<br>", ":raw-html:`<br/>`")

	return fmt.Sprint(directives, string(title), render), nil
}

// Uses pandoc to convert mediawiki to rst (without additional modifications)
func (page *Page) pandocWikiToRst() (rendered string, err error) {
	// raw=$(cat $PAGE | pandoc -f mediawiki -t rst)
	// TODO: instead of chan, mv-on-write?
	cmd := exec.Command("pandoc", "-f", "mediawiki", "-t", "rst")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer stdout.Close()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	defer stderr.Close()
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	ch := make(chan error, 1)
	go func(ch chan<- error) {
		defer stdin.Close()
		_, err := stdin.Write([]byte(page.LatestRevision().Text))
		ch <- err
	}(ch)

	cmd.Start()
	if err != nil {
		return "", err
	}
	err = <-ch
	if err != nil {
		return "", err
	}
	outAll, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	errAll, err := io.ReadAll(stderr)
	if err != nil {
		return "", err
	}
	if err = cmd.Wait(); err != nil {
		logger.Debugf("\n------\nSTDIN\n------\n%s\n------\nSTDERR\n------\n%s", page.LatestRevision().Text, errAll)
		return "", err
	}
	return string(outAll), nil
}
