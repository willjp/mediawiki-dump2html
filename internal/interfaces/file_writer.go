package interfaces

import "github.com/spf13/afero"

type FileWriter interface {
	Write(file afero.File, b []byte) (n int, err error)
}
