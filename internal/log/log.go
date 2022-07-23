package log

import (
	"os"

	"willpittman.net/x/logger"
)

// Seam so we can stub logger in tests
var Log = logger.New(os.Stderr)
