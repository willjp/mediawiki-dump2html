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
	args := this.args()
	cmd := cmd{Cmd: exec.Command("pandoc", args...)}
	return this.execute(&cmd)
}

// Commandline Arguments for Pandoc
func (this *Pandoc) args() []string {
	args := []string{"-f", this.From, "-t", this.To}
	if this.Standalone == true {
		args = append(args, "--standalone")
	}
	return args
}

// Invokes pandoc on CLI (testable seam).
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
