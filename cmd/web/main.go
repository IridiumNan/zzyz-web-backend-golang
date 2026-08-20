package main

import (
	"log"
	"log/slog"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/config"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/database"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/utils"
)

func main() {
	// WARN: This func panic if fail
	config.InitGlobalConfig()

	// fmt.Println(webConfig)

	// Init the logger and use defer to close log file
	logFilePath := config.GetLogFilePath()
	logFile, err := utils.InitLogger(logFilePath)
	if err != nil {
		slog.Error("error when init logger, use Stdout as default logger", "error", err)
	}

	defer func() {
		// Close logFile before the process end
		if logFile != nil {
			err := logFile.Close()
			if err != nil {
				slog.Error("error when close log file", "err", err)
			}
		}
	}()

	// NOTE: automatical create if not exist
	// Open database
	database.OpenLocalDB(config.GlobalWebConfig.DatabasePath)

	// starting the taker for handling db for member
	go database.RunDBWriteTasker()

	defer func() {
		err := database.CloseLocalDB()
		if err != nil {
			utils.TextLogger.Error("error when close database", "err", err)
		}
	}()

	// Run the internal router for member modify, check the [getInternalRouter] for endpoints
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
