package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func Initialise() (*sqlx.DB, error) {
	// Open up a temporary file to store log data
	db, err := sqlx.Open("sqlite3", "glimpse_temp.db")
	if err != nil {
		return nil, fmt.Errorf("error opening db: %v", err)
	}

	// Create tables
	var schema = `
	CREATE TABLE jsonlogs (
		level TEXT,
		timestamp NUMBER,
		message TEXT,
		raw TEXT
	);
	`

	db.MustExec(schema)

	return db, nil
}
