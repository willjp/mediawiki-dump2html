package utils

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
	stubs "willpittman.net/x/mediawiki-to-sphinxdoc/internal/test/stubs"
)

func TestFileReplace(t *testing.T) {
	t.Run("Writes File", func(t *testing.T) {
		file := stubs.NewFakeFile()
		osCreate = func(path string) (interfaces.OsFile, error) {
			return &file, nil
		}
		errs := FileReplace("abc", "/var/tmp/foo.txt")
		assert.Nil(t, errs)

		res, err := io.ReadAll(file.Buffer)
		assert.Nil(t, err)
		assert.Equal(t, []byte("abc"), res)
	})
}
