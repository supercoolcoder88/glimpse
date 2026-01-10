package logs

import (
	"fmt"
	"slices"

	"github.com/jmoiron/sqlx"
)

type filterRule struct {
	field    string
	value    any
	operator string
}

var operators = []string{"=", "<", "<=", ">", ">=", "LIKE"}
var allowedFields = []string{"level", "ts", "message"}

// Filters JSON logs using operators provided by sqlite
func FilterJSONLog(db *sqlx.DB, rules []filterRule) ([]JSONLog, error) {
	query := `SELECT * FROM jsonlogs WHERE 1=1 `

	values := make(map[string]interface{}) // this allows sqlx to handle the type for the query automatically
	for _, rule := range rules {
		if !slices.Contains(operators, rule.operator) {
			return nil, fmt.Errorf("invalid operator: %s", rule.operator)
		}

		// Check if filters are correct
		if !slices.Contains(allowedFields, rule.field) {
			return nil, fmt.Errorf("invalid field to filter: %s", rule.field)
		}

		values[rule.field] = rule.value
		if rule.operator == "LIKE" {
			query += fmt.Sprintf("AND %s LIKE :%s ", rule.field, rule.field)
		} else {
			query += fmt.Sprintf("AND %s%s:%s ", rule.field, rule.operator, rule.field)
		}
	}

	// Make query
	rows, err := db.NamedQuery(query, values)

	if err != nil {
		fmt.Printf("error querying logs: %v", err)
	}

	var logs []JSONLog
	log := JSONLog{}
	for rows.Next() {
		err := rows.StructScan(&log)
		if err != nil {
			fmt.Printf("error scanning struct: %s", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
