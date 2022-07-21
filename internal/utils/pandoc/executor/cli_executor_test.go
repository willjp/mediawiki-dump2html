package utils

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeCmdExecutor struct{}

func (this *FakeCmdExecutor) Execute(stdin io.Reader) (render string, errs []error) {
	result, _ := io.ReadAll(stdin)
	return string(result), nil
}

func TestCliExecutor(t *testing.T) {
	val := "abc"
	stdin := strings.NewReader(val)
	cmd := FakeCmdExecutor{}
	result, errs := cmd.Execute(stdin)

	assert.Nil(t, errs)
	assert.Equal(t, val, result)
}
