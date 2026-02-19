package logs

import (
	"glimpse/db"
	"strings"
	"testing"
)

func TestLogsRead_Success(t *testing.T) {
	input := `2026-01-06T17:30:01Z INFO [auth_service] User login successful: user_id=8823
{"level":"info","ts":1704562270,"msg":"structured log record","request_id":"req-99"}
ts=2026-01-06T17:31:40Z level=INFO component=worker_node_2 msg="Health check passed"`

	reader := strings.NewReader(input)

	conn := db.InitialiseTest(t)
	ch := make(chan Entry)

	go func() {
		if err := Read(reader, ch, conn); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}()

	var received []string
	for entry := range ch {
		received = append(received, entry.Raw)
	}

	expected := 3
	if len(received) != expected {
		t.Fatalf("expected %d entries, got %d", expected, len(received))
	}

	if received[0] != "2026-01-06T17:30:01Z INFO [auth_service] User login successful: user_id=8823" {
		t.Errorf("incorrect value, got '%s'", received[0])
	}
}
