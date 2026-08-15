package utils

import (
	"io"
	"log/slog"
	"os"
	"path"
)

var TextLogger *slog.Logger

func InitLogger(logFilePath string) (logFile *os.File, err error) {
	err = EnsureDir(path.Dir(logFilePath))
	if err != nil {
		return
	}

	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		slog.Error("fail to open log file", "error", err)

		return
	}

	mw := io.MultiWriter(os.Stdout, logFile)

	handler := slog.NewTextHandler(mw, nil)
	TextLogger = slog.New(handler)

	slog.Info("system logger initialized", "log_file_path", logFilePath)

	return
}

func CloseLogger() {
}
