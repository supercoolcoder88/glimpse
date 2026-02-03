package logs

import (
	"strings"
	"testing"
)

func TestLogsRead_Success(t *testing.T) {
	input := "this is a testing line"
	reader := strings.NewReader(input)

	ch := make(chan Entry)
	if err := Read(reader, ch); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
