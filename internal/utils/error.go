package utils

import (
	"os"

	"willpittman.net/x/logger"
	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

var osRemove = os.Remove

// file expects a os.File
func RmFileOn(file interfaces.OsFile, err error) {
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
