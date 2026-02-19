package logs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Logs will be converted to this struct for processing and display
type Entry struct {
	Level            string                 `json:"level" db:"level"`
	Timestamp        int                    `json:"ts" db:"timestamp"`
	Message          string                 `json:"msg" db:"message"`
	Raw              string                 `json:"-" db:"raw"`
	AdditionalFields map[string]interface{} // TODO: Implement this in future
	isUnknownFormat  bool
}

// Read will process the logs. Logs are categorised as either JSON format or unformatted (will look to add more processing)
// it will return the log as an Entry class via the "output" channel to allow for processing of the log for UI display purposes. Logs
// are passed into the "dbCh" channel to allow for filtering of the data when the user requests via a sqlite db.
func Read(input io.Reader, output chan Entry, db *sqlx.DB) error {
	defer close(output)
	scanner := bufio.NewScanner(input)

	// goroutine to handle database operations
	dbCh := make(chan Entry)
	go func() {
		for entry := range dbCh {
			if entry.isUnknownFormat {
				db.Exec(`INSERT INTO logs (raw) VALUES ($1)`,
					entry.Raw,
				)
			} else {
				db.Exec(`INSERT INTO logs (level, timestamp, message, raw) VALUES ($1, $2, $3, $4)`,
					entry.Level,
					entry.Timestamp,
					entry.Message,
					entry.Raw,
				)
			}

		}
	}()

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		raw := string(line)

		// Check for logfmt
		re := regexp.MustCompile(`(\w+)=("[^"]*"|[^"\s]+)`)
		matches := re.FindAllStringSubmatch(raw, -1)

		isLogFmt := len(matches) > 0
		// JSON Log parsing
		if bytes.HasPrefix(line, []byte("{")) {
			var jsonLog Entry

			if err := json.Unmarshal(line, &jsonLog); err != nil {
				return fmt.Errorf("failed to unmarshal JSON: %v", err)
			}
			jsonLog.Raw = raw

			dbCh <- jsonLog
			output <- jsonLog

		} else if isLogFmt {
			entry := Entry{
				Raw:              raw,
				AdditionalFields: make(map[string]interface{}),
			}

			for _, match := range matches {
				key := match[1]
				value := match[2]

				if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
					value = value[1 : len(value)-1]
				}

				// TODO: Make this more dynamic later
				switch key {
				case "level":
					entry.Level = value
				case "ts":
					if ts, err := strconv.Atoi(value); err == nil {
						entry.Timestamp = ts
					}
				case "msg":
					entry.Message = value
				default:
					entry.AdditionalFields[key] = value
				}
			}

			dbCh <- entry
			output <- entry
		} else {
			e := Entry{Raw: raw, isUnknownFormat: true}
			dbCh <- e
			output <- e
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}
