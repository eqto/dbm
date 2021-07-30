package mysql

import (
	"fmt"
	"strings"

	"github.com/eqto/go-db/query"
)

func queryInsert(stmt *query.InsertStmt) string {
	table := query.TableOf(stmt)

	fieldStrs := []string{}
	valueStrs := []string{}
	values := query.ValueOf(stmt)

	for i, field := range table.Fields {
		fieldStrs = append(fieldStrs, field.Name)
		if len(values) > i {
			valueStrs = append(valueStrs, values[i])
		} else {
			valueStrs = append(valueStrs, `?`)
		}
	}
	return fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s)`, table.Name, strings.Join(fieldStrs, `, `), strings.Join(valueStrs, `, `))
}
