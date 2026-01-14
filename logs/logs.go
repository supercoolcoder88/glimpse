package logs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"glimpse/db"
	"io"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type L struct {
	Message string
}

type JSONLog struct {
	Level     string `json:"level" db:"level"`
	Timestamp int    `json:"ts" db:"timestamp"`
	Message   string `json:"msg" db:"message"`
	Raw       string `json:"-" db:"raw"`
}

func Read(input io.Reader) error {
	scanner := bufio.NewScanner(input)

	conn, err := db.Initialise()

	if err != nil {
		return fmt.Errorf("db error: %v", err)
	}

	defer conn.Close()

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())

		// JSON Log parsing
		if bytes.HasPrefix(line, []byte("{")) {
			var jsonLog JSONLog

			if err := json.Unmarshal(line, &jsonLog); err != nil {
				return fmt.Errorf("failed to unmarshal JSON: %v", err)
			}

			conn.Exec(`INSERT INTO jsonlogs (level, timestamp, message, raw) VALUES ($1, $2, $3, $4)`,
				jsonLog.Level,
				jsonLog.Timestamp,
				jsonLog.Message,
				line,
			)
			rules := []filterRule{
				{
					field:    "level",
					value:    "error",
					operator: "LIKE",
				},
			}
			// PrintFilteredLog(conn, rules)
		} else {
			//fmt.Printf("%s\n", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}

func PrintJSONLogs(db *sqlx.DB) {
	var logs []JSONLog
	if err := db.Select(&logs, "SELECT level, timestamp, message, raw FROM jsonlogs"); err != nil {
		fmt.Printf("error printing logs: %v", err)
	}
	fmt.Println(logs)
}

// func PrintFilteredLog(db *sqlx.DB, rules []filterRule) {
// 	logs, _ := FilterJSONLog(db, rules)
// 	fmt.Println(logs)
// }
