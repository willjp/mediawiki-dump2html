package utils

import (
	"encoding/xml"
	"io"
	"os/exec"

	"github.com/lithammer/dedent"
	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/html"
)

// Struct of pandoc CLI options.
type PandocOptions struct {
	From       string
	To         string
	Standalone bool
}

// Wraper around a pandoc conversion.
func PandocConvert(page *elements.Page, opts *PandocOptions) (string, error) {
	// raw=$(cat $PAGE | pandoc -f mediawiki -t rst)
	// TODO: instead of chan, mv-on-write?
	args := []string{"-f", opts.From, "-t", opts.To}
	if opts.Standalone == true {
		args = append(args, "--standalone")
	}
	cmd := exec.Command("pandoc", args...)
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

// Extracts pandoc's generated CSS from a page render.
//
// When pandoc is called using the `--standalone` param, it renders CSS into each page.
// This extracts that CSS, so that you could dump it to a file and reference it within each page.
func PandocExtractCss(page *elements.Page) (rendered string, err error) {
	var htmlNode html.Html
	opts := PandocOptions{From: "mediawiki", To: "html", Standalone: true}
	raw, err := PandocConvert(page, &opts)
	if err != nil {
		return "", err
	}
	xml.Unmarshal([]byte(raw), &htmlNode)
	css := dedent.Dedent(htmlNode.Head.Style)
	return css, nil
}
