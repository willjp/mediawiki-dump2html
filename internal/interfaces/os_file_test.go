package interfaces

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsFile(t *testing.T) {
	var implementsIface = func(i OsFile) bool {
		return true
	}
	assert.True(t, implementsIface(&os.File{}))
}
