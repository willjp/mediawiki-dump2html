package utils

import (
	"fmt"
	"os"
	"regexp"
)

var nonPosixChars *regexp.Regexp

func init() {
	nonPosixChars = regexp.MustCompile(fmt.Sprint("[^", `A-Za-z0-9\._\-`, os.PathSeparator, "]"))
}

func SanitizePath(path []byte) []byte {
	return nonPosixChars.ReplaceAll(path, []byte("_"))
}
