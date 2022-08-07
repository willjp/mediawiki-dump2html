// Contains the Application's logger instance.
// Exported so that it can be overridden with a fake to make log assertions in tests.
package log

import (
	"os"

	"github.com/willjp/go-logger"
)

var Log logger.Interface

func init() {
	logRaw := logger.New(os.Stderr)
	Log = &logRaw
}
