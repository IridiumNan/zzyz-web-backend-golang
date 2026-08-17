package main

import (
	"log"
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

	go func() {
		// NOTE: internalRouter -> Use goroutine to run it
		internalRouter := getInternalRouter()
		err := internalRouter.Run(config.GlobalWebConfig.InternalAddress)
		if err != nil {
			utils.TextLogger.Error("error to start internal router", "err", err)
		}
	}()

	// mainRouter
	mainRouter := getMainRouter()

	err = mainRouter.Run(config.GlobalWebConfig.MainAddress)
	if err != nil {
		log.Fatal(err)
	}
}
