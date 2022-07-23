package cli

import (
	"errors"
)

var InvalidArgsError = errors.New("[ERROR] invalid args")

// Indicates that there was an error parsing commandline arguments.
// Concretion of interfaces.CliCommand
type InvalidArgs struct{}

func (this *InvalidArgs) Call() error {
	return InvalidArgsError
}
