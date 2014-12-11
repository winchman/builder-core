package dockerfile

import (
	"github.com/rafecolton/go-fileutils"
)

/*
Func ideas:

- env vars
- inbox
- from (or regex)
- docker (ps / images)
*/

func Plus(src, dest string) error {
	//copy
	if err := fileutils.CpWithArgs(src, dest, fileutils.CpArgs{PreserveModTime: true}); err != nil {
		return err
	}

	// prerocess

	// return preprocessingFunc()
	return nil
}
