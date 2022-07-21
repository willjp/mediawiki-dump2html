package interfaces

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecCmd(t *testing.T) {
	var implementsIface = func(i ExecCmd) bool {
		return true
	}
	assert.True(t, implementsIface(&exec.Cmd{}))
}
