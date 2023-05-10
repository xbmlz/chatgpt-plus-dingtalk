package utils

import (
	"os"
	"path/filepath"
)

func Mkdir(path string) (err error) {
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0755)
	return
}
