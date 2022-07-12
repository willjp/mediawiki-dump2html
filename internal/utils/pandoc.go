package utils

import (
	"io"
	"os/exec"

	"willpittman.net/x/logger"
)

// Interface for a Cmd
type Cmd interface {
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
}

// Struct representing a pandoc command.
// (primarily intended for format conversions).
type PandocOpts struct {
	From       string
	To         string
	Standalone bool
}

// Commandline Arguments for Pandoc
func (this *PandocOpts) args() []string {
	args := []string{"-f", this.From, "-t", this.To}
	if this.Standalone == true {
		args = append(args, "--standalone")
	}
	return args
}

func (this *PandocOpts) Command() *PandocCmd {
	args := this.args()
	return &PandocCmd{Cmd: exec.Command("pandoc", args...)}
}

// Wraps exec.Cmd and adds methods related to executing a pandoc command.
type PandocCmd struct {
	Cmd
}

// Low-Level method to invoke pandoc on CLI (testable seam).
func (this *PandocCmd) Execute(stdin io.Reader) (string, error) {
	stdout, err := this.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer stdout.Close()
	stderr, err := this.StderrPipe()
	if err != nil {
		return "", err
	}
	defer stderr.Close()
	stdinW, err := this.StdinPipe()
	if err != nil {
		return "", err
	}
	ch := make(chan error, 1)
	go func(ch chan<- error) {
		defer stdinW.Close()
		data, err := io.ReadAll(stdin)
		if err != nil {
			ch <- err
			return
		}

		_, err = stdinW.Write(data)
		ch <- err
	}(ch)

	err = this.Start()
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
	if err = this.Wait(); err != nil {
		logger.Debugf("STDERR:\n%s", errAll)
		return "", err
	}
	return string(outAll), nil
}
