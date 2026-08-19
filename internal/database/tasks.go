package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

type sqlWriteTaskFunc func(context.Context, *sql.Tx) error

type sqlWriteTask struct {
	execFunc sqlWriteTaskFunc

	// NOTE: sqlStr is for log recored
	// WARN: Don't use it as the transcation query
	sqlStr string
}

var (
	// TODO: write a func to listen the stop signal
	sqlWriteStopChan chan int = make(chan int)

	sqlWriteTaskChan chan sqlWriteTask = make(chan sqlWriteTask, 500)
)

// RunDBWriteTasker : This Tasker just handle the request for post update and
// NOTE: Run it before you start router
func RunDBWriteTasker() {
	utils.TextLogger.Info("RunDBWriteTasker: starting task")
	for {
		select {
		case <-sqlWriteStopChan:
			utils.TextLogger.Info("RunDBWriteTasker: received stop signal, exiting")

			// TODO: store remaining write tasks to disk and load it next time
			return

		case task, ok := <-sqlWriteTaskChan:
			if !ok {
				utils.TextLogger.Warn("RunDBWriteTasker: channel closed, exiting")
				return
			}

			func() {
				utils.TextLogger.Info("RunDBWriteTasker: executing sql", "sqlStr", task.sqlStr)

				tx, err := globalDB.Begin()
				if err != nil {
					utils.TextLogger.Error("RunDBWriteTasker: error when begin transcation", "err", err)
					return
				}

				defer tx.Rollback()

				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
				defer cancel()

				err = task.execFunc(ctx, tx)
				if err != nil {
					utils.TextLogger.Error("RunDBWriteTasker: error when exec sql", "err", err, "sqlStr", task.sqlStr)
				}

				if err := tx.Commit(); err != nil {
					utils.TextLogger.Error("RunDBWriteTasker: error when commit transcation", "err", err)
				}

				utils.TextLogger.Info("RunDBWriteTasker: exec success", "sqlStr", task.sqlStr)
			}()
		}
	}
}
