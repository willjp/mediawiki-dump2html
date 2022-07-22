package appfs

import "github.com/spf13/afero"

// Interface to use in place of os.* (has in-memory stub interface)
var AppFs = afero.NewOsFs()
