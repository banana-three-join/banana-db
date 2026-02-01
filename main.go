package main

import (
	"fmt"
	"math/rand/v2"
	"os"
)

/*tmp files assure that the only files that will be created, are the files that complete their execution*/
func SaveData(path string, data []byte) error {
	tmp := fmt.Sprintf("%s.tmp.%d", path, rand.Uint())
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		return err
	}
	defer func() {
		fp.Close()
		if err != nil {
			os.Remove(tmp)
		}
	}()

	if _, err := fp.Write(data); err != nil {
		return err
	}

	if err := fp.Sync(); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}
