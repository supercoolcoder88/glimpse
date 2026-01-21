package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func Initialise(env string) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	if env == "test" {
		// Hold file in memory for tests
		db, err = sqlx.Open("sqlite3", ":memory:")
	} else {
		// Open up a temporary file to store log data
		db, err = sqlx.Open("sqlite3", "glimpse_temp.db")
	}

	if err != nil {
		return nil, fmt.Errorf("error opening db: %v", err)
	}

	// Create tables
	var schema = `
	CREATE TABLE json_logs (
		level TEXT,
		timestamp NUMBER,
		message TEXT,
		raw TEXT
	);
	CREATE TABLE unformatted_logs (
		raw TEXT
	);
	`

	db.MustExec(schema)

	return db, nil
}
