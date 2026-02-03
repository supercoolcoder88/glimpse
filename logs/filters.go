package logs

import (
	"fmt"
	"slices"

	"github.com/jmoiron/sqlx"
)

type FilterRule struct {
	Field    string
	Value    any
	Operator string
}

type filter struct {
	Db *sqlx.DB
}

var (
	CommonFields      = []string{"level", "ts", "message"}
	AllowedOperations = []string{"=", "<", "<=", ">", ">=", "LIKE"}
)

func NewFilterRule(f string, v any, o string) (*FilterRule, error) {
	if !slices.Contains(AllowedOperations, o) {
		return nil, fmt.Errorf("invalid Operator: %s", o)
	}

	if !slices.Contains(CommonFields, f) {
		return nil, fmt.Errorf("invalid field to filter: %s", f)
	}

	return &FilterRule{
		Field:    f,
		Value:    v,
		Operator: o,
	}, nil
}

func NewFilter(db *sqlx.DB) *filter {
	return &filter{
		Db: db,
	}
}

// Filters JSON logs using Operators provided by sqlite
func (f *filter) HandleJSON(rules []FilterRule) ([]JSON, error) {
	query := `SELECT * FROM logs WHERE 1=1 `

	Values := make(map[string]interface{}) // this allows sqlx to handle the type for the query automatically
	for _, rule := range rules {

		Values[rule.Field] = rule.Value
		if rule.Operator == "LIKE" {
			query += fmt.Sprintf("AND %s LIKE :%s ", rule.Field, rule.Field)
		} else {
			query += fmt.Sprintf("AND %s%s:%s ", rule.Field, rule.Operator, rule.Field)
		}
	}

	// Make query
	rows, err := f.Db.NamedQuery(query, Values)

	if err != nil {
		fmt.Printf("error querying logs: %v", err)
	}

	var logs []JSON
	log := JSON{}
	for rows.Next() {
		err := rows.StructScan(&log)
		if err != nil {
			fmt.Printf("error scanning struct: %s", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
