package main

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/models"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/utils"
	"github.com/gin-gonic/gin"
)

var pathSet models.PathSet = *models.DefaultPathSet()

// 20 MiB for package.zip
const (
	maxPackageSize = 5 << 22
	packageName    = "package.zip"
)

func getZipCacheDir() string {
	cacheDir := fmt.Sprintf("package_%v", time.Now().Unix())

	return path.Join(pathSet.ZipCacheDir, cacheDir)
}

func postUploadHandler(c *gin.Context) {
	// Receive single zip package

	uploadPackage, err := c.FormFile("package")
	if err != nil {
		utils.TextLogger.Error("error when get package for endpoint", "endpoint", "/posts/upload", "err", err)
		c.JSON(http.StatusBadRequest, models.NewBadResponse("fail to get package from form", err))
		return
	}

	if uploadPackage.Size > maxPackageSize {
		utils.TextLogger.Error("upload package to large, reject it automatically")
		c.JSON(http.StatusBadRequest, models.NewBadResponse("max package size: 20 MB", fmt.Errorf("package too large, max size: %s", "20 MiB")))
		return
	}

	dstDir := getZipCacheDir()

	err = utils.EnsureDir(dstDir)
	if err != nil {
		utils.TextLogger.Error("error when ensure dir when create zip cache path", "err", err, "path", dstDir)
		c.JSON(http.StatusInternalServerError, models.NewBadResponse("internal message", fmt.Errorf("error when create tmp dir for package")))
		return
	}

	packagePath := path.Join(dstDir, packageName)
	err = c.SaveUploadedFile(uploadPackage, packagePath, 0o755)
	if err != nil {
		utils.TextLogger.Error("error when receive uploaded package", "err", err)
		c.JSON(http.StatusInternalServerError, models.NewBadResponse("error", fmt.Errorf("error when receive uploaded package: %w", err)))
		return
	}

	c.JSON(http.StatusOK, models.NewDataResponse(fmt.Sprintf("upload success, size: %d", uploadPackage.Size)))

	// TODO: push new task to chan
}
