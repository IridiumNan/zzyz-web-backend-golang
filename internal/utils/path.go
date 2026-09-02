package utils

import (
	"fmt"
	"os"
	"path"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/models"
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

func EnsureDataDirs() (err error) {
	reportFunc := func(targetPath string, err error) error {
		return fmt.Errorf("error when ensure dir, path: %s, err: %w", targetPath, err)
	}

	pathSet := models.DefaultPathSet()

	err = EnsureDir(pathSet.HTMLPendingDir)
	if err != nil {
		return reportFunc(pathSet.HTMLPendingDir, err)
	}

	err = EnsureDir(pathSet.HTMLReleaseDir)
	if err != nil {
		return reportFunc(pathSet.HTMLReleaseDir, err)
	}

	err = EnsureDir(pathSet.IDDir)
	if err != nil {
		return reportFunc(pathSet.IDDir, err)
	}

	err = EnsureDir(pathSet.TitleDir)
	if err != nil {
		return reportFunc(pathSet.TitleDir, err)
	}

	err = EnsureDir(pathSet.RawDir)
	if err != nil {
		return reportFunc(pathSet.RawDir, err)
	}

	err = EnsureDir(pathSet.ZipCacheDir)
	if err != nil {
		return reportFunc(pathSet.ZipCacheDir, err)
	}

	return
}
