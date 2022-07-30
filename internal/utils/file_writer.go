package utils

import "github.com/spf13/afero"

type FileWriter struct{}

func (this *FileWriter) Write(file afero.File, b []byte) (n int, err error) {
	return file.Write(b)
}
