package main

import (
	"log/slog"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/config"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/utils"
)

func main() {
	// panic if fail
	config.InitGlobalConfig()

	// fmt.Println(webConfig)

	// Init the logger and use defer to close log file
	logFilePath := config.GetLogFilePath()
	logFile, err := utils.InitLogger(logFilePath)
	if err != nil {
		slog.Error("error when init logger, use Stdout as default logger", "error", err)
	}

	defer func() {
		if logFile != nil {
			err := logFile.Close()
			if err != nil {
				slog.Error("error when close log file", "err", err)
			}
		}
	}()

	// Router
	router := getGinRouter()

	router.Run()
}
