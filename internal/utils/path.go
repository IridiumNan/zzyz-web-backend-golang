package utils

import (
	"os"
	"path"
)

// EnsureDir : Ensure the dir exist, if not exist it will create
// return error when create this dir
func EnsureDir(dirPath string) (err error) {
	err = os.MkdirAll(dirPath, 0o755)

	return
}

// EnsureFile : Ensure the file exist, if not exist it will create
// return error what ever it encounterd
func EnsureFile(filePath string) (err error) {
	if _, err = os.Stat(filePath); err == nil {
		return
	}

	err = os.MkdirAll(path.Dir(filePath), 0o755)
	if err != nil {
		return
	}

	_, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY, 0o644)

	return
}
