package utils

import (
	"encoding/xml"
	"io"
	"os/exec"

	"github.com/lithammer/dedent"
	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/html"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/elements/mwdump"
)

// Struct of pandoc CLI options.
type Pandoc struct {
	Stdin      io.Reader
	From       string
	To         string
	Standalone bool
}

func (this *Pandoc) Execute() (string, error) {
	args := this.args()
	cmd := cmd{Cmd: exec.Command("pandoc", args...)}
	return this.execute(&cmd)
}

func (this *Pandoc) args() []string {
	args := []string{"-f", this.From, "-t", this.To}
	if this.Standalone == true {
		args = append(args, "--standalone")
	}
	return args
}

func (this *Pandoc) execute(c *cmd) (string, error) {
	stdout, err := c.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer stdout.Close()
	stderr, err := c.StderrPipe()
	if err != nil {
		return "", err
	}
	defer stderr.Close()
	stdin, err := c.StdinPipe()
	if err != nil {
		return "", err
	}
	ch := make(chan error, 1)
	go func(ch chan<- error) {
		defer stdin.Close()
		data, err := io.ReadAll(this.Stdin)
		if err != nil {
			ch <- err
			return
		}

		_, err = stdin.Write(data)
		ch <- err
	}(ch)

	err = c.Start()
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
	if err = c.Wait(); err != nil {
		logger.Debugf("STDERR:\n%s", errAll)
		return "", err
	}
	return string(outAll), nil
}

// Test Seam that wraps exec.Cmd
type cmd struct {
	*exec.Cmd
}

// Wraper around a pandoc conversion.
func PandocConvert(page *mwdump.Page, opts *Pandoc) (string, error) {
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

	err = cmd.Start()
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
func PandocExtractCss(page *mwdump.Page) (rendered string, err error) {
	var htmlNode html.Html
	opts := Pandoc{From: "mediawiki", To: "html", Standalone: true}
	raw, err := PandocConvert(page, &opts)
	if err != nil {
		return "", err
	}
	xml.Unmarshal([]byte(raw), &htmlNode)
	css := dedent.Dedent(htmlNode.Head.Style)
	return css, nil
}
