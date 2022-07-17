package utils

import (
	"io"
	"sync"

	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/utils"
)

// Wraps exec.Cmd and adds methods related to executing a pandoc command.
type Cmd struct {
	utils.Cmd
	Args []string
}

// Low-Level method to invoke pandoc on CLI (testable seam).
func (this *Cmd) Execute(stdin io.Reader) (render string, errs []error) {
	// record goroutine errors
	wg := sync.WaitGroup{}
	wg.Add(1)
	ch := make(chan error, 2)
	defer func(ch <-chan error) {
		wg.Wait()
		err := <-ch
		if err != nil {
			errs = append(errs, err)
		}
	}(ch)

	// build pipes
	stdout, err := this.StdoutPipe()
	if err != nil {
		errs = append(errs, err)
	}
	stderr, err := this.StderrPipe()
	if err != nil {
		errs = append(errs, err)
	}
	stdinW, err := this.StdinPipe()
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return "", errs
	}

	// write stdin
	go func(ch chan<- error) {
		defer func() {
			err := stdinW.Close()
			if err != nil {
				ch <- err
			}
			close(ch)
			wg.Done()
		}()
		data, err := io.ReadAll(stdin)
		if err != nil {
			ch <- err
			return
		}
		if _, err = stdinW.Write(data); err != nil {
			ch <- err
		}
	}(ch)

	// run command
	err = this.Start()
	if err != nil {
		errs = append(errs, err)
		return "", errs
	}

	outAll, err := io.ReadAll(stdout)
	if err != nil {
		errs = append(errs, err)
	}
	errAll, err := io.ReadAll(stderr)
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return "", errs
	}

	err = this.Wait()
	if err != nil {
		logger.Debugf("STDERR:\n%s", errAll)
		errs = append(errs, err)
	}
	return string(outAll), errs
}
