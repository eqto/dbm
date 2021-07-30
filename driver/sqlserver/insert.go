package sqlserver

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

	paramCounter := 0
	for i, field := range table.Fields {
		fieldStrs = append(fieldStrs, field.Name)
		if len(values) > i && values[i] != `?` {
			valueStrs = append(valueStrs, values[i])
		} else {
			paramCounter++
			valueStrs = append(valueStrs, fmt.Sprintf(`@p%d`, paramCounter))
		}
	}
	return fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s)`, table.Name, strings.Join(fieldStrs, `, `), strings.Join(valueStrs, `, `))
}
