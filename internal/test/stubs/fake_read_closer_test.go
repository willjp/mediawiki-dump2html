package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeReadCloserRead(t *testing.T) {
	t.Run("Returns ReadError when provided", func(t *testing.T) {
		var ExpectedError = errors.New("Expected")
		buf := make([]byte, 0)
		rCloser := FakeReadCloser{Reader: strings.NewReader("abc"), ReadError: ExpectedError}
		_, err := rCloser.Read(buf)
		assert.Error(t, err, ExpectedError)
	})

	t.Run("Returns Written when provided", func(t *testing.T) {
		buf := make([]byte, 3)
		rCloser := FakeReadCloser{Reader: strings.NewReader("abc")}
		_, err := rCloser.Read(buf)
		assert.Nil(t, err)
		assert.Equal(t, []byte("abc"), buf)
	})
}

func TestFakeReadCloserClose(t *testing.T) {
	t.Run("Returns CloseError when provided", func(t *testing.T) {
		var ExpectedError = errors.New("Expected")
		rCloser := FakeReadCloser{CloseError: ExpectedError}
		err := rCloser.Close()
		assert.Error(t, err, ExpectedError)
	})
}
