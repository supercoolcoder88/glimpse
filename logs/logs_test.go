package logs

import (
	"strings"
	"testing"
)

func TestLogIngest_Success(t *testing.T) {
	input := "this is a testing line"
	reader := strings.NewReader(input)

	if err := Read(reader); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFilter_Level_Success(t *testing.T) {

}
