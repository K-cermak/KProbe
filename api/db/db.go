package db

import (
	"database/sql"
	"errors"
	"os"

	"KProbeAPI/helpers"
	_ "modernc.org/sqlite"
)

const (
	DB_CONNECTION_FAILED  string = "Failed to connect to database"
	DB_QUERY_FAILED       string = "Failed to get value from database"
	DB_SCAN_NEWEST_FAILED string = "Failed to get newest scan from database"
	RES_OK                string = "OK"
)

// Default for testing. Override in build process using: go build -ldflags "-X KProbeAPI/db.dbPath=/opt/kprobe/db.sqlite"
var dbPath = "../cli/db.sqlite"

func getDatabaseConnection() (*sql.DB) {
	if !databaseExist() {
		helpers.PrintError("Database does not exist, run CLI app with <kprobe db init> first")
		return nil
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		helpers.PrintError("Failed to connect to database")
		return nil
	}

	return db
}

func closeDatabase(db *sql.DB) {
	if db == nil {
		return
	}

	if err := db.Close(); err != nil {
		helpers.PrintError("Failed to close database connection (" + err.Error() + ")")
	}
}

func databaseExist() bool {
	if _, err := os.Stat(dbPath); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		helpers.PrintError("Failed to check database existence (" + err.Error() + ")")
	}

	return false
}

func GetValue(key string) (string, string) {
	db := getDatabaseConnection()
	if db == nil {
		return "", DB_CONNECTION_FAILED
	}
	defer closeDatabase(db)

	var value string
	query := `
	SELECT value
	FROM keys
	WHERE name = ?;`

	err := db.QueryRow(query, key).Scan(&value)
	if err != nil {
		helpers.PrintError("Failed to get value from database (" + err.Error() + ")")
		return "", DB_QUERY_FAILED
	}

	return value, RES_OK
}

func GetScanNewest(scanName string) (helpers.ScanRes, string) {
	db := getDatabaseConnection()
	if db == nil {
		return helpers.ScanRes{}, DB_CONNECTION_FAILED
	}
	defer closeDatabase(db)

	query := `
	SELECT generated, passed
	FROM history
	WHERE scan_name = ?
	ORDER BY generated DESC
	LIMIT 1;
	`

	var res helpers.ScanRes

	err := db.QueryRow(query, scanName).Scan(&res.Generated, &res.Passed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.ScanRes{}, DB_SCAN_NEWEST_FAILED
		}
		helpers.PrintError("Failed to get data from database (" + err.Error() + ")")
		return helpers.ScanRes{}, DB_CONNECTION_FAILED
	}

	return res, RES_OK
}
