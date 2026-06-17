package db

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"KProbeCLI/helpers"

	_ "modernc.org/sqlite"
)

// const dbPath = "db.sqlite" //FOR TESTING, CHANGE TO BELOW
const dbPath = "/opt/kprobe/db.sqlite"
const dbVersion = "v1.1"

var DB *sql.DB

func connectDatabase() {
	if !DatabaseExist() {
		helpers.PrintError(true, "Database does not exist, run <kprobe db init> first")
	}

	var err error

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		helpers.PrintError(true, "Failed to connect to database ("+err.Error()+")")
	}

	if !DatabaseVersionCheck() {
		helpers.PrintError(true, "Database version mismatch, you probably upgraded the application. Run <kprobe db reset> to reinitialize the database. Migration of the database is not supported.")
	}
}

func InitDatabase() {
	var err error

	if DatabaseExist() {
		err = os.Remove(dbPath)
		if err != nil {
			helpers.PrintError(true, "Failed to delete database ("+err.Error()+")")
		}
	}

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		helpers.PrintError(true, "Failed to connect to database ("+err.Error()+")")
	}

	createTableQuery := `
	CREATE TABLE keys (
		name VARCHAR(32) PRIMARY KEY,
		value VARCHAR(4096) NOT NULL
	);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		helpers.PrintError(true, "Failed to create table ("+err.Error()+")")
	}

	createTableQuery = `
	CREATE TABLE scans (
		name VARCHAR(32) PRIMARY KEY,
		type VARCHAR(4) CHECK(type IN ('ping', 'http')),
		address VARCHAR(256) NOT NULL,
		timeout INTEGER,
		status_code VARCHAR(256),
		keyword TEXT
	);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		helpers.PrintError(true, "Failed to create table ("+err.Error()+")")
	}

	createTableQuery = `
	CREATE TABLE history (
		scan_name VARCHAR(32) NOT NULL,
		generated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		passed BOOLEAN NOT NULL,
		delete_after DATETIME NOT NULL,
		PRIMARY KEY (scan_name, generated)
	);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		helpers.PrintError(true, "Failed to create table ("+err.Error()+")")
	}

	InsertValue("probe_name", "New Probe")
	InsertValue("db_version", dbVersion)
	InsertValue("db_init_time", time.Now().String())

	InsertValue("config_set", "false")
	InsertValue("delete_after", "7")

	InsertValue("api_port", "80")
	InsertValue("editor_endpoint", "true")

	InsertValue("ping_retries", "5")
	InsertValue("ignore_ssl_errors", "false")
}

func DatabaseVersionCheck() bool {
	if DB == nil {
		connectDatabase()
	}

	var dbRealVersion string
	
	query := `
	SELECT value
	FROM keys
	WHERE name = 'db_version';
	`

	err := DB.QueryRow(query).Scan(&dbRealVersion)
	if err != nil {
		helpers.PrintError(true, "Failed to get database version ("+err.Error()+")")
	}

	return dbRealVersion == dbVersion
}



func DatabaseExist() bool {
	if _, err := os.Stat(dbPath); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		helpers.PrintError(true, "Failed to check database existence ("+err.Error()+")")
	}

	return false
}

func GetValue(key string) string {
	if DB == nil {
		connectDatabase()
	}

	var value string

	query := `
	SELECT value
	FROM keys
	WHERE name = ?;
	`

	err := DB.QueryRow(query, key).Scan(&value)
	if err != nil {
		helpers.PrintError(true, "Failed to get data from database ("+err.Error()+")")
	}

	return value
}

func InsertValue(key string, value string) {
	if DB == nil {
		connectDatabase()
	}

	insertQuery := `
	INSERT INTO keys (name, value) 
	VALUES (?, ?) 
	ON CONFLICT(name) 
	DO UPDATE SET value = excluded.value;
	`

	_, err := DB.Exec(insertQuery, key, value)
	if err != nil {
		helpers.PrintError(true, "Failed to insert data into database ("+err.Error()+")")
	}
}

func GetScans() []helpers.Scan {
	if DB == nil {
		connectDatabase()
	}

	var scans []helpers.Scan

	query := `
	SELECT name, type, address, timeout, status_code, keyword
	FROM scans;
	`

	rows, err := DB.Query(query)
	if err != nil {
		helpers.PrintError(true, "Failed to get data from database ("+err.Error()+")")
	}

	for rows.Next() {
		var scan helpers.Scan

		err = rows.Scan(&scan.Name, &scan.Type, &scan.Address, &scan.Timeout, &scan.StatusCode, &scan.Keyword)
		if err != nil {
			helpers.PrintError(true, "Failed to scan data from database ("+err.Error()+")")
		}

		scans = append(scans, scan)
	}

	return scans
}

func AddScan(scan helpers.Scan) {
	if DB == nil {
		connectDatabase()
	}

	insertQuery := `
	INSERT INTO scans (name, type, address, timeout, status_code, keyword) 
	VALUES (?, ?, ?, ?, ?, ?);
	`

	_, err := DB.Exec(insertQuery, scan.Name, scan.Type, scan.Address, scan.Timeout, scan.StatusCode, scan.Keyword)
	if err != nil {
		helpers.PrintError(true, "Failed to insert data into database ("+err.Error()+")")
	}
}

func DeleteScans() {
	if DB == nil {
		connectDatabase()
	}

	deleteQuery := `
	DELETE FROM scans;
	`

	_, err := DB.Exec(deleteQuery)
	if err != nil {
		helpers.PrintError(true, "Failed to delete scans from database ("+err.Error()+")")
	}
}

func GetScanRes(scanName string, start string, end string) []helpers.ScanRes {
	if DB == nil {
		connectDatabase()
	}

	var scanRes []helpers.ScanRes

	query := `
	SELECT generated, passed
	FROM history
	WHERE scan_name = ? AND generated BETWEEN ? AND ?
	ORDER BY generated ASC;
	`

	rows, err := DB.Query(query, scanName, start, end)
	if err != nil {
		helpers.PrintError(true, "Failed to get data from database ("+err.Error()+")")
	}

	for rows.Next() {
		var res helpers.ScanRes

		err = rows.Scan(&res.Generated, &res.Passed)
		if err != nil {
			helpers.PrintError(true, "Failed to scan data from database ("+err.Error()+")")
		}

		scanRes = append(scanRes, res)
	}

	return scanRes
}

func AddScanRes(scanName string, passed bool) {
	if DB == nil {
		connectDatabase()
	}

	daysToDelete := GetValue("delete_after")

	insertQuery := `
	INSERT INTO history (scan_name, passed, delete_after) 
	VALUES (?, ?, datetime('now', 'localtime', '+' || ? || ' days'));
	`

	_, err := DB.Exec(insertQuery, scanName, passed, daysToDelete)
	if err != nil {
		helpers.PrintError(true, "Failed to insert data into database ("+err.Error()+")")
	}
}

func DeleteOldScanRes() {
	if DB == nil {
		connectDatabase()
	}

	deleteQuery := `
	DELETE FROM history
	WHERE delete_after < datetime('now', 'localtime');
	`

	_, err := DB.Exec(deleteQuery)
	if err != nil {
		helpers.PrintError(true, "Failed to delete old results from database ("+err.Error()+")")
	}
}

func GetScanNewest(scanName string) (helpers.ScanRes, bool) {
	if DB == nil {
		connectDatabase()
	}

	query := `
	SELECT generated, passed
	FROM history
	WHERE scan_name = ?
	ORDER BY generated DESC
	LIMIT 1;
	`

	var res helpers.ScanRes

	err := DB.QueryRow(query, scanName).Scan(&res.Generated, &res.Passed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helpers.ScanRes{}, false
		}
		helpers.PrintError(true, "Failed to get data from database ("+err.Error()+")")
	}

	return res, true
}
