package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	test "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
)

func TestRmFileOn(t *testing.T) {
	t.Run("Deletes file on error", func(t *testing.T) {
		var removed string
		osRemove = func(file string) error {
			removed = file
			return nil
		}
		file := test.FakeFile{Path: "/var/tmp/foo.txt"}
		RmFileOn(&file, fmt.Errorf("error"))
		assert.Equal(t, "/var/tmp/foo.txt", removed)
	})

	t.Run("Does not delete file when no error", func(t *testing.T) {
		var removed string
		osRemove = func(file string) error {
			removed = file
			return nil
		}
		file := test.FakeFile{Path: "/var/tmp/foo.txt"}
		RmFileOn(&file, nil)
		assert.Equal(t, "", removed)
	})
}
