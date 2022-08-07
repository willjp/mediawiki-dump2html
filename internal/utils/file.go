package utils

import "github.com/willjp/mediawiki-dump2html/internal/appfs"

func FileReplace(contents string, filepath string) (errs []error) {
	file, err := appfs.AppFs.Create(filepath)
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
