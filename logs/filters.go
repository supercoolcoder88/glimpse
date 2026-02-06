package logs

import (
	"fmt"
	"slices"

	"github.com/jmoiron/sqlx"
)

type Rule struct {
	field    string
	value    any
	operator string
}

type filter struct {
	Db *sqlx.DB
}

var (
	CommonFields      = []string{"level", "ts", "message"}
	AllowedOperations = []string{"=", "<", "<=", ">", ">=", "LIKE"}
)

func NewRule(f string, v any, o string) (*Rule, error) {
	if !slices.Contains(AllowedOperations, o) {
		return nil, fmt.Errorf("invalid operator: %s", o)
	}

	if !slices.Contains(CommonFields, f) {
		return nil, fmt.Errorf("invalid field to filter: %s", f)
	}

	return &Rule{
		field:    f,
		value:    v,
		operator: o,
	}, nil
}

func NewFilter(db *sqlx.DB) *filter {
	return &filter{
		Db: db,
	}
}

func (f *filter) Apply(rules []Rule) ([]Entry, error) {
	query := `SELECT * FROM logs WHERE 1=1 `

	values := make(map[string]interface{}) // this allows sqlx to handle the type for the query automatically
	for _, rule := range rules {

		values[rule.field] = rule.value
		if rule.operator == "LIKE" {
			query += fmt.Sprintf("AND %s LIKE :%s ", rule.field, rule.field)
		} else {
			query += fmt.Sprintf("AND %s%s:%s ", rule.field, rule.operator, rule.field)
		}
	}

	// Make query
	rows, err := f.Db.NamedQuery(query, values)

	if err != nil {
		fmt.Printf("error querying logs: %v", err)
	}

	var logs []Entry
	log := Entry{}
	for rows.Next() {
		err := rows.StructScan(&log)
		if err != nil {
			fmt.Printf("error scanning struct: %s", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
