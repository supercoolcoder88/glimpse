package logs

import (
	"glimpse/db"
	"testing"
)

func TestFilterJSONLogs_Success(t *testing.T) {
	rules := []FilterRule{
		{
			Field:    "level",
			Value:    "error",
			Operator: "=",
		},
	}

	db, _ := db.Initialise("tests")

	log := JSONLog{Level: "error", Timestamp: 1704562111, Message: "test error, varying Values"}
	line := `{"level":"error","ts":1704562270,"msg":"test error, varying Values"}`
	db.Exec(`INSERT INTO json_logs (level, timestamp, message, raw) VALUES ($1, $2, $3, $4)`,
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
	rules := []FilterRule{
		{
			Field:    "level",
			Value:    "error",
			Operator: "equals",
		},
	}

	db, _ := db.Initialise("tests")

	filter := NewFilter(db)
	_, err := filter.HandleJSON(rules)

	if err == nil {
		t.Fatalf("filter should fail with bad Operator: %s", rules[0].Operator)
	}
}

func TestFilterJSONLogs_BadField_ShouldError(t *testing.T) {
	rules := []FilterRule{
		{
			Field:    "blah",
			Value:    "error",
			Operator: "=",
		},
	}

	db, _ := db.Initialise("tests")

	filter := NewFilter(db)
	_, err := filter.HandleJSON(rules)
	if err == nil {
		t.Fatalf("filter should fail with bad Field: %s", rules[0].Field)
	}
}
