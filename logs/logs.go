package logs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Logs will be converted to this struct for processing and display
type Entry struct {
	Level          string                 `json:"level" db:"level"`
	Timestamp      int                    `json:"ts" db:"timestamp"`
	Message        string                 `json:"msg" db:"message"`
	Raw            string                 `json:"-" db:"raw"`
	AdditionalFields map[string]interface{} // TODO: Implement this in future
	isJSON         bool
}

// Read will process the logs. Logs are categorised as either JSON format or unformatted (will look to add more processing)
// it will return the log as an Entry class via the "output" channel to allow for processing of the log for UI display purposes. Logs
// are passed into the "dbCh" channel to allow for filtering of the data when the user requests via a sqlite db.
func Read(input io.Reader, output chan Entry, db *sqlx.DB) error {
	scanner := bufio.NewScanner(input)

	// goroutine to handle database operations
	dbCh := make(chan Entry, 100)
	go func() {
		for entry := range dbCh {
			if entry.isJSON {
				db.Exec(`INSERT INTO logs (level, timestamp, message, raw) VALUES ($1, $2, $3, $4)`,
					entry.Level,
					entry.Timestamp,
					entry.Message,
					entry.Raw,
				)
			} else {
				db.Exec(`INSERT INTO logs (raw) VALUES ($1)`,
					entry.Raw,
				)
			}

		}
	}()

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		raw := string(line)

		// JSON Log parsing
		if bytes.HasPrefix(line, []byte("{")) {
			var jsonLog Entry

			if err := json.Unmarshal(line, &jsonLog); err != nil {
				return fmt.Errorf("failed to unmarshal JSON: %v", err)
			}

			dbCh <- jsonLog
			output <- jsonLog
		} else {
			e := Entry{Raw: raw}
			dbCh <- e
			output <- e
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}
