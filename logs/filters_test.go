package logs

import (
	"glimpse/db"
	"testing"
)

func TestFilterJSONLogs_Success(t *testing.T) {
	rules := []Rule{}

	rule, _ := NewRule("level", "error", "=")
	rules = append(rules, *rule)

	db := db.InitialiseTest(t)

	log := Entry{Level: "error", Timestamp: 1704562111, Message: "test error, varying Values"}
	line := `{"level":"error","ts":1704562270,"msg":"test error, varying Values"}`
	db.Exec(`INSERT INTO logs (level, timestamp, message, raw) VALUES ($1, $2, $3, $4)`,
		log.Level,
		log.Timestamp,
		log.Message,
		line,
	)

	filter := NewFilter(db)
	logs, err := filter.Apply(rules)
	if err != nil {
		t.Fatalf("failed to filter: %s", err)
	}

	if len(logs) != 1 {
		t.Fatalf("wanted %d logs, got %d", 1, len(logs))
	}
}

func TestFilterJSONLogs_BadOperator_ShouldError(t *testing.T) {
	_, err := NewRule("level", "ERROR", "equals")

	if err == nil {
		t.Fatalf("rule should fail with bad operator: equals")
	}
}

func TestFilterJSONLogs_BadField_ShouldError(t *testing.T) {
	_, err := NewRule("error", "ERROR", "=")

	if err == nil {
		t.Fatalf("rule should fail with bad field: error")
	}
}
