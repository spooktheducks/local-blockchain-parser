package blockdb

import (
	"fmt"
)

type DataNotIndexedError struct {
	Index string
}

func (e DataNotIndexedError) Error() string {
	return fmt.Sprintf(`It seems that you need to build the %v index before running this command.  Try running the "builddb %v" command on your .dat files.`, e.Index, e.Index)
}
