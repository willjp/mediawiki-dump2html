package cli

import (
	"errors"
)

var InvalidArgsError = errors.New("[ERROR] invalid args")

type InvalidArgs struct{}

func (this *InvalidArgs) Call() error {
	return InvalidArgsError
}
