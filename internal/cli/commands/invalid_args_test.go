package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/willjp/mediawiki-dump2html/internal/interfaces"
)

func TestInvalidArgsImplementsInterface(t *testing.T) {
	implementsInterface := func(interfaces.CliCommand) bool {
		return true
	}
	assert.True(t, implementsInterface(&InvalidArgs{}))
}

func TestInvalidArgsReturnsError(t *testing.T) {
	cmd := InvalidArgs{}
	err := cmd.Call()
	assert.Error(t, InvalidArgsError, err)
}
