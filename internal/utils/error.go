package utils

import (
	"os"

	"willpittman.net/x/logger"
)

var osRemove = os.Remove

type Namer interface {
	Name() string
}

// file expects a os.File
func RmFileOn(file Namer, err error) {
	if err != nil {
		logger.Errorf("Error encountered, removing: %s", file.Name())
		osRemove(file.Name())
	}
}

func PanicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func LogWarnOn(err error) {
	if err != nil {
		logger.Warnf("Ignored Error: %s", err)
	}
}
