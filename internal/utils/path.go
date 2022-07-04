package utils

import (
	"fmt"
	"os"
	"regexp"
)

var nonPosixChars *regexp.Regexp

func init() {
	var err error
	nonPosixChars, err = regexp.Compile(fmt.Sprint("[^", `A-Za-z0-9\._\-`, os.PathSeparator, "]"))
	PanicOn(err)
}

func SanitizePath(path []byte) []byte {
	return nonPosixChars.ReplaceAll(path, []byte("-"))
}
