package utils

import (
	"io"

	"willpittman.net/x/logger"
)

// Wraps exec.Cmd and adds methods related to executing a pandoc command.
type PandocCmd struct {
	Cmd
	Args []string
}

// Low-Level method to invoke pandoc on CLI (testable seam).
func (this *PandocCmd) Execute(stdin io.Reader) (render string, errs []error) {
	// build pipes
	stdout, err := this.StdoutPipe()
	if err != nil {
		return "", append(errs, err)
	}
	defer func() {
		err := stdout.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}()

	stderr, err := this.StderrPipe()
	if err != nil {
		return "", append(errs, err)
	}
	defer func() {
		err = stderr.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}()

	stdinW, err := this.StdinPipe()
	if err != nil {
		errs = append(errs, err)
		return "", errs
	}
	ch := make(chan error, 2)
	go func(ch chan<- error) {
		defer func() {
			ch <- stdinW.Close()
			close(ch)
		}()
		data, err := io.ReadAll(stdin)
		if err != nil {
			ch <- err
			return
		}

		_, err = stdinW.Write(data)
		ch <- err
	}(ch)

	// run command
	err = this.Start()
	if err != nil {
		errs = append(errs, err)
		return "", errs
	}
	for {
		err, ok := <-ch
		if !ok {
			break
		} else if err != nil {
			errs = append(errs, err)
			return "", errs
		}
	}
	outAll, err := io.ReadAll(stdout)
	if err != nil {
		errs = append(errs, err)
		return "", errs
	}
	errAll, err := io.ReadAll(stderr)
	if err != nil {
		errs = append(errs, err)
		return "", errs
	}
	if err = this.Wait(); err != nil {
		logger.Debugf("STDERR:\n%s", errAll)
		errs = append(errs, err)
		return "", errs
	}
	return string(outAll), nil
}
