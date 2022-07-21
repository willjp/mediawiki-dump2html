package pandoc

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

// implements interfaces/CmdExecutor
type FakeCmd struct {
	args []string
}

func (this *FakeCmd) Execute(stdin io.Reader) (render string, errs []error) {
	return "", nil
}

func (this *FakeCmd) Args() []string {
	return this.args
}

// implements interfaces/PandocExecutor
type FakeCmdExecutor struct {
	cmd interfaces.CmdExecutor
}

func (this *FakeCmdExecutor) Execute(stdin io.Reader) (render string, errs []error) {
	result, _ := io.ReadAll(stdin)
	return string(result), nil
}

func (this *FakeCmdExecutor) Args() []string {
	return this.cmd.Args()
}

// tests
func TestExecute(t *testing.T) {
	t.Run("Forwards STDIN, and performs cmd.Execute()", func(t *testing.T) {
		val := "abc"
		stdin := strings.NewReader(val)
		cmd := FakeCmdExecutor{}
		executor := Executor{}
		result, errs := executor.Execute(&cmd, stdin)

		assert.Nil(t, errs)
		assert.Equal(t, val, result)
	})
}

func TestArgs(t *testing.T) {
	t.Run("Returns contained cmd.Args()", func(t *testing.T) {
		args := []string{"foo", "-a", "--verbose"}
		cmd := FakeCmd{args: args}
		pandocCmd := FakeCmdExecutor{cmd: &cmd}
		executor := Executor{cmd: &pandocCmd}

		assert.Equal(t, args, executor.Args())
	})
}
