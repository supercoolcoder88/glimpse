package logs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"glimpse/db"
	"io"
	"reflect"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Display is what will be shown to users in the TUI
type Entry struct {
	Raw            string
	ExpectedFields map[string]string
	SpecialFields  map[string]string
}

type JSON struct {
	Level     string `json:"level" db:"level"`
	Timestamp int    `json:"ts" db:"timestamp"`
	Message   string `json:"msg" db:"message"`
	Raw       string `json:"-" db:"raw"`
}

type Unformatted struct {
	Raw string `json:"-" db:"raw"`
}

func NewEntry(r string, e map[string]string, s map[string]string) *Entry {
	return &Entry{
		Raw:            r,
		ExpectedFields: e,
		SpecialFields:  s,
	}
}

func Read(input io.Reader, ch chan Entry) error {
	scanner := bufio.NewScanner(input)

	conn, err := db.Initialise("")

	if err != nil {
		return fmt.Errorf("db error: %v", err)
	}

	defer conn.Close()

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())

		// JSON Log parsing
		if bytes.HasPrefix(line, []byte("{")) {
			var jsonLog JSON

			if err := json.Unmarshal(line, &jsonLog); err != nil {
				return fmt.Errorf("failed to unmarshal JSON: %v", err)
			}

			// Send jsonLog to db
			conn.Exec(`INSERT INTO logs (level, timestamp, message, raw) VALUES ($1, $2, $3, $4)`,
				jsonLog.Level,
				jsonLog.Timestamp,
				jsonLog.Message,
				line,
			)

			expected := make(map[string]string)
			special := make(map[string]string) // TODO: Add special fields

			v := reflect.ValueOf(jsonLog)
			t := v.Type()

			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				tag := field.Tag.Get("json")

				if tag == "-" || tag == "" {
					continue
				}

				value := fmt.Sprintf("%v", v.Field(i).Interface())

				expected[tag] = value
			}

			e := NewEntry(
				jsonLog.Raw,
				expected,
				special,
			)

			ch <- *e
		} else {
			conn.Exec(`INSERT INTO logs (raw) VALUES ($1)`,
				line,
			)

			e := &Entry{Raw: string(line)}

			ch <- *e
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}

func PrintJSONLogs(db *sqlx.DB) {
	var logs []JSON
	if err := db.Select(&logs, "SELECT level, timestamp, message, raw FROM logs"); err != nil {
		fmt.Printf("error printing logs: %v", err)
	}
	fmt.Println(logs)
}

// func PrintFilteredLog(db *sqlx.DB, rules []FilterRule) {
// 	logs, _ := FilterJSONLog(db, rules)
// 	fmt.Println(logs)
// }
