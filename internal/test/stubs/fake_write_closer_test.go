package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeWriterWrite(t *testing.T) {
	t.Run("Returns WriteError when provided", func(t *testing.T) {
		var ExpectedError = errors.New("Expected")
		writer := strings.Builder{}
		wCloser := FakeWriteCloser{Writer: &writer, WriteError: ExpectedError}
		_, err := wCloser.Write([]byte("abc"))
		assert.Error(t, err, ExpectedError)
	})
}

func TestFakeWriterClose(t *testing.T) {
	t.Run("Returns CloseError when provided", func(t *testing.T) {
		var ExpectedError = errors.New("Expected")
		wCloser := FakeWriteCloser{CloseError: ExpectedError}
		err := wCloser.Close()
		assert.Error(t, err, ExpectedError)
	})
}
