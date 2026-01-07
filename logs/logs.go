package logs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"glimpse/db"
	"io"
	"os"
	"os/signal"
	"syscall"

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

	// Cleanup function
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nShutting down... cleaning up temporary database.")
		os.Remove("glimpse_temp.db")
		os.Exit(0)
	}()

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

func PrintFilteredLog(db *sqlx.DB) {
	var logs []JSONLog
	if err := db.Select(&logs, "SELECT level, timestamp, message, raw FROM jsonlogs WHERE level=$1", "error"); err != nil {
		fmt.Printf("error printing logs: %v", err)
	}
	fmt.Println(logs)
}
