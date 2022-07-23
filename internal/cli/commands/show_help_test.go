package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

func TestShowHelpImplementsInterface(t *testing.T) {
	implementsInterface := func(interfaces.CliCommand) bool {
		return true
	}
	assert.True(t, implementsInterface(&ShowHelp{}))
}

func TestShowHelpCall(t *testing.T) {
	cmd := ShowHelp{}
	err := cmd.Call()
	assert.Nil(t, err)
}
