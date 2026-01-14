package logs

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func generateTestDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil
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

	return db
}

func TestFilterJSONLogs_Success(t *testing.T) {
	rules := []filterRule{
		{
			field:    "level",
			value:    "error",
			operator: "=",
		},
	}

	db := generateTestDB()

	log := JSONLog{Level: "error", Timestamp: 1704562111, Message: "test error, varying values"}
	line := `{"level":"error","ts":1704562270,"msg":"test error, varying values"}`
	db.Exec(`INSERT INTO jsonlogs (level, timestamp, message, raw) VALUES ($1, $2, $3, $4)`,
		log.Level,
		log.Timestamp,
		log.Message,
		line,
	)

	filter := NewFilter(db)
	logs, err := filter.HandleJSON(rules)
	if err != nil {
		t.Fatalf("failed to filter: %s", err)
	}

	if len(logs) != 1 {
		t.Fatalf("wanted %d logs, got %d", 1, len(logs))
	}
}

func TestFilterJSONLogs_BadOperator_ShouldError(t *testing.T) {
	rules := []filterRule{
		{
			field:    "level",
			value:    "error",
			operator: "equals",
		},
	}

	db := generateTestDB()

	filter := NewFilter(db)
	_, err := filter.HandleJSON(rules)

	if err == nil {
		t.Fatalf("filter should fail with bad operator: %s", rules[0].operator)
	}
}

func TestFilterJSONLogs_BadField_ShouldError(t *testing.T) {
	rules := []filterRule{
		{
			field:    "blah",
			value:    "error",
			operator: "=",
		},
	}

	db := generateTestDB()

	filter := NewFilter(db)
	_, err := filter.HandleJSON(rules)
	if err == nil {
		t.Fatalf("filter should fail with bad field: %s", rules[0].field)
	}
}
