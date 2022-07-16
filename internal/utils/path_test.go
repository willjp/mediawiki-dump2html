package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSantizePath(t *testing.T) {
	cases := []struct {
		test    string
		path    []byte
		expects []byte
	}{
		{
			test:    "Valid Paths are not changed",
			path:    []byte("MyFile-0_1_1.txt"),
			expects: []byte("MyFile-0_1_1.txt"),
		},
		{
			test:    "Groups of InValid Characters are replaced with underscores",
			path:    []byte("Subject: Name > Foo!.txt"),
			expects: []byte("Subject__Name___Foo_.txt"),
		},
	}
	for _, tcase := range cases {
		t.Run(tcase.test, func(t *testing.T) {
			res := SanitizeFilename(tcase.path)
			assert.Equal(t, string(res), string(tcase.expects))
		})
	}
}
