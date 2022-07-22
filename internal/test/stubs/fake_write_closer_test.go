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

	t.Run("Returns bytes written to Writer on success", func(t *testing.T) {
		writer := strings.Builder{}
		wCloser := FakeWriteCloser{Writer: &writer}
		bytes, err := wCloser.Write([]byte("abc"))
		assert.Nil(t, err)
		assert.Equal(t, 3, bytes)
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

func TestFakeWriterString(t *testing.T) {
	t.Run("Returns written string", func(t *testing.T) {
		writer := strings.Builder{}
		wCloser := FakeWriteCloser{Writer: &writer}
		_, err := wCloser.Write([]byte("abc"))
		assert.Nil(t, err)
		res := wCloser.String()
		assert.Equal(t, "abc", res)
	})
}
