package log

import (
	"os"

	"willpittman.net/x/logger"
)

var Log logger.Interface

func init() {
	logRaw := logger.New(os.Stderr)
	Log = &logRaw
}