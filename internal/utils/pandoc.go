package utils

import (
	"io"
	"os/exec"

	"willpittman.net/x/logger"
)

// Struct representing a pandoc command.
// (primarily intended for format conversions).
type Pandoc struct {
	Stdin      io.Reader
	From       string
	To         string
	Standalone bool
}

// Executes a pandoc command.
func (this *Pandoc) Execute() (string, error) {
	args := this.Args()
	cmd := Cmd{Cmd: exec.Command("pandoc", args...)}
	return this.ExecuteCmd(&cmd)
}

// Commandline Arguments for Pandoc
func (this *Pandoc) Args() []string {
	args := []string{"-f", this.From, "-t", this.To}
	if this.Standalone == true {
		args = append(args, "--standalone")
	}
	return args
}

// Low-Level method to invoke pandoc on CLI (testable seam).
func (this *Pandoc) ExecuteCmd(c *Cmd) (string, error) {
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

// Low-Level test Seam that wraps exec.Cmd
type Cmd struct {
	*exec.Cmd
}
