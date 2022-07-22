package utils

import (
	"os"

	"willpittman.net/x/mediawiki-to-sphinxdoc/internal/interfaces"
)

var osCreate = func(path string) (interfaces.OsFile, error) {
	return os.Create(path)
}

func FileReplace(contents string, filepath string) (errs []error) {
	file, err := osCreate(filepath)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}()

	_, err = file.WriteString(contents)
	if err != nil {
		RmFileOn(file, err)
		errs = append(errs, err)
		return errs
	}
	return nil
}
