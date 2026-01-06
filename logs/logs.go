package logs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type L struct {
	id      int
	Message string
}

type JSONLog struct {
	Level     string                 `json:"level"`
	Ts        int                    `json:"ts"`
	Message   string                 `json:"msg"`
	RequestID string                 `json:"request_id"`
	Extra     map[string]interface{} `json:"-"`
}

func Read(input io.Reader) error {
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())

		// JSON Log parsing
		if bytes.HasPrefix(line, []byte("{")) {
			var jsonLog JSONLog

			if err := json.Unmarshal(line, &jsonLog); err != nil {
				fmt.Printf("failed to unmarshal JSON: %v", line)
			} else {
				if err := json.Unmarshal(line, &jsonLog.Extra); err != nil {
					fmt.Printf("failed to unmarshal JSON: %v", line)
				}
				delete(jsonLog.Extra, "level")
				delete(jsonLog.Extra, "ts")
				delete(jsonLog.Extra, "msg")
				delete(jsonLog.Extra, "request_id")
				fmt.Printf("JSON Data: %+v\n", jsonLog)
			}
		} else {
			fmt.Printf("%s\n", line)
		}

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}
