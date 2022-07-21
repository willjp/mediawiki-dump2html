package interfaces

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func implementsIface(i ExecCmd) bool {
	return true
}

func TestInterface(t *testing.T) {
	// interface should cover a real exec.Cmd{}
	assert.True(t, implementsIface(&exec.Cmd{}))
}
