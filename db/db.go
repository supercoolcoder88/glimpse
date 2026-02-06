package db

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
)

func Initialise() (*sqlx.DB, error) {
	var db *sqlx.DB

	db, err := sqlx.Open("sqlite3", "glimpse_temp.db")

	if err != nil {
		return nil, fmt.Errorf("error opening db: %v", err)
	}

	// Create tables
	var schema = `
	CREATE TABLE logs (
		raw TEXT NOT NULL,
		level TEXT,
		timestamp NUMBER,
		message TEXT
	)
	`

	db.MustExec(schema)

	return db, nil
}

func InitialiseTest(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	var schema = `
	CREATE TABLE logs (
		raw TEXT NOT NULL,
		level TEXT,
		timestamp NUMBER,
		message TEXT
	)
	`

	db.MustExec(schema)

	t.Cleanup(func() {
		db.Close()
	})

	return db
}
