package utils

import (
	"os"

	"willpittman.net/x/logger"
)

func RmFileOn(file *os.File, err error) {
	if err != nil {
		logger.Errorf("Error encountered, removing: %s", file.Name())
		os.Remove(file.Name())
	}
}
