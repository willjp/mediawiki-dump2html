package renderers

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/lithammer/dedent"
	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

// Has methods for conversion, and keeps state used during conversion
type RST struct{}

func (rst *RST) Filename(page *elements.Page) string {
	fileName := fmt.Sprint(page.Title, ".rst")
	return string(utils.SanitizePath([]byte(fileName)))
}

// Converts mediawiki text to rst, with tweaks so it behaves well with sphinx-docs.
func (rst *RST) Render(page *elements.Page) (rendered string, err error) {
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
	pandocRender, err := rst.pandocWikiToRst(page)
	if err != nil {
		return "", err
	}

	// replace '<br>' with something rst understands
	render := strings.ReplaceAll(pandocRender, "<br>", ":raw-html:`<br/>`")

	return fmt.Sprint(directives, string(title), render), nil
}

// Uses pandoc to convert mediawiki to rst (without additional modifications)
func (rst *RST) pandocWikiToRst(page *elements.Page) (rendered string, err error) {
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
