package database

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/utils"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed init.sql
var initSQL string

// driverName for sql
const driverName = "sqlite3"

// globalDB for internal/database useage
var globalDB *sql.DB

// InitDB : init the database file, Use log.Fatal to handle the err it returns
// It will ensure the file and dir exist then create tables
// There is no need to ensure path again
func initDB(databasePath string) error {
	funcErrDescription := "error when init database"

	var err error
	err = utils.EnsureDir(path.Dir(databasePath))
	if err != nil {
		return fmt.Errorf("%s, path: %s, err: %w", funcErrDescription, databasePath, err)
	}

	err = utils.EnsureFile(databasePath)
	if err != nil {
		return fmt.Errorf("%s, path: %s, err: %w", funcErrDescription, databasePath, err)
	}

	var db *sql.DB
	db, err = sql.Open(driverName, databasePath)
	if err != nil {
		return fmt.Errorf("%s, path: %s, driverName: %s, err: %w", funcErrDescription, databasePath, "sqlite3", err)
	}

	defer db.Close()

	_, err = db.Exec(initSQL)
	if err != nil {
		return fmt.Errorf("%s, exec initSQL, sql: %s\nerr: %w", funcErrDescription, initSQL, err)
	}

	return nil
}

// OpenLocalDB : Open database (sqlite) by the path from configuration
// if database file not exist, it will use the func initDB to create a new one
// The db is assert to the variant globalDB
// Remember to call CloseLocalDB using defer
func OpenLocalDB(databasePath string) {
	funcErrDescription := "error when open database"

	if globalDB != nil {
		slog.Warn("deplicate open database, ignore it")
		return
	}

	var err error

	if _, err = os.Stat(databasePath); os.IsNotExist(err) {
		utils.TextLogger.Info("database file not found, create a new one", "path", databasePath)
		err = initDB(databasePath)
		if err != nil {
			// Raise Fatal if fail to init
			log.Fatal(err)
		}
	}

	globalDB, err = sql.Open(driverName, databasePath)
	if err != nil {
		log.Fatal(funcErrDescription, err)
	}
}

// CloseLocalDB :  the global database
// return the err of Close() function
func CloseLocalDB() error {
	return globalDB.Close()
}
