package elements

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"willpittman.net/x/logger"
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
	panicAndRmOn := func(file *os.File, err error) {
		if err == nil {
			return
		}
		logger.Errorf("Error encountered, removing: %s", file.Name())
		os.Remove(file.Name())
		panic(err)
	}

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

		logger.Infof("Writing: %s\n", rstPath)
		rendered, err := page.renderRst()
		panicAndRmOn(file, err)
		_, err = file.WriteString(rendered)
		panicAndRmOn(file, err)
	}
}

func (page *Page) renderRst() (rendered string, err error) {
	directives := `
	.. role:: raw-html(raw)
	  :format: html
	`

	// page title between '='s
	titleLen := len([]rune(page.Title))
	title := fmt.Sprint(
		strings.Repeat("=", titleLen), "\n",
		page.Title, "\n",
		strings.Repeat("=", titleLen), "\n",
	)

	// raw=$(cat $PAGE | pandoc -f mediawiki -t rst)
	// TODO: instead of chan, copy-on-write?
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
		logger.Errorf("\n------\nSTDIN\n------\n%s\n------\nSTDERR\n------\n%s", page.LatestRevision().Text, errAll)
		return "", err
	}

	// replace '<br>' with something rst understands
	render := strings.ReplaceAll(string(outAll), "<br>", ":raw-html:`<br/>`")

	return fmt.Sprintf(directives, title, render), nil
}
