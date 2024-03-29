package pandoc

import (
	"io"
	"os/exec"
	"sync"

	"github.com/willjp/mediawiki-dump2html/internal/interfaces"
	"github.com/willjp/mediawiki-dump2html/internal/log"
)

// Wraps exec.Cmd and adds methods related to executing a pandoc command.
type Cmd struct {
	interfaces.ExecCmd
	args                []string
	skipCheckExecutable bool // (test seam) skip checking executable.
}

// Array of command/params passed on CLI
func (this *Cmd) Args() []string {
	return this.args
}

// Executes pandoc on CLI
func (this *Cmd) Execute(stdin io.Reader) (render string, errs []error) {
	if !this.skipCheckExecutable {
		this.checkExecutable()
	}

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

	stdinW, stdout, stderr, errs := this.buildPipes()
	if len(errs) > 0 {
		return "", errs
	}

	go this.writeStdin(stdinW, stdin, ch, &wg)

	render, errs = this.start(stdout, stderr)
	if len(errs) > 0 {
		return "", errs
	}
	return render, errs
}

// Builds stdin(Writer), stdout(Reader), stderr(Reader) connected to this pandoc process
func (this *Cmd) buildPipes() (stdin io.WriteCloser, stdout io.Reader, stderr io.Reader, errs []error) {
	var err error
	stdin, err = this.StdinPipe()
	if err != nil {
		errs = append(errs, err)
	}
	stdout, err = this.StdoutPipe()
	if err != nil {
		errs = append(errs, err)
	}
	stderr, err = this.StderrPipe()
	if err != nil {
		errs = append(errs, err)
	}
	return stdin, stdout, stderr, errs
}

// Pipes `reader` (mediawiki-src) into `writer` (pandoc-stdin)
func (this *Cmd) writeStdin(writer io.WriteCloser, reader io.Reader, ch chan<- error, wg *sync.WaitGroup) {
	defer func() {
		err := writer.Close()
		if err != nil {
			ch <- err
		}
		close(ch)
		wg.Done()
	}()
	data, err := io.ReadAll(reader)
	if err != nil {
		ch <- err
		return
	}
	if _, err = writer.Write(data); err != nil {
		ch <- err
	}
}

// Runs the pandoc process, capturing stdout/stderr
func (this *Cmd) start(stdout io.Reader, stderr io.Reader) (render string, errs []error) {
	var err error

	// run the command
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
		log.Log.Debugf("PANDOC ARGS: %s", this.Args())
		log.Log.Debugf("STDERR:\n%s", errAll)
		errs = append(errs, err)
	}
	return string(outAll), errs
}

// Panic if executable does not exist
func (this *Cmd) checkExecutable() {
	// panic if pandoc not available
	executablePath := this.Args()[0]
	_, err := exec.LookPath(executablePath)
	if err != nil {
		panic(err)
	}
}
