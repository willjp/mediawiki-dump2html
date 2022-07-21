package pandoc

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
	stubs "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
)

// implements interfaces/PandocExecutor
type FakePandocExecutor struct {
	cmd interfaces.CmdExecutor
}

func (this *FakePandocExecutor) Execute(stdin io.Reader) (render string, errs []error) {
	result, _ := io.ReadAll(stdin)
	return string(result), nil
}

func (this *FakePandocExecutor) Args() []string {
	return this.cmd.Args()
}

// tests
func TestExecute(t *testing.T) {
	t.Run("Forwards STDIN, and performs cmd.Execute()", func(t *testing.T) {
		val := "abc"
		stdin := strings.NewReader(val)
		cmd := FakePandocExecutor{}
		executor := Executor{}
		result, errs := executor.Execute(&cmd, stdin)

		assert.Nil(t, errs)
		assert.Equal(t, val, result)
	})
}

func TestArgs(t *testing.T) {
	t.Run("Returns contained cmd.Args()", func(t *testing.T) {
		args := []string{"foo", "-a", "--verbose"}
		cmd := stubs.FakeCmd{CliArgs: args}
		pandocCmd := FakePandocExecutor{cmd: &cmd}
		executor := Executor{cmd: &pandocCmd}

		assert.Equal(t, args, executor.Args())
	})
}
