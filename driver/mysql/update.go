package mysql

import (
	"fmt"
	"strings"

	"github.com/eqto/go-db/query"
)

func queryUpdate(stmt *query.UpdateStmt) string {
	tables := query.TableOf(stmt)
	if len(tables) == 0 {
		return ``
	}
	table := tables[0]
	fieldStrs := []string{}
	for _, field := range table.Fields {
		fieldStrs = append(fieldStrs, field.Name+` = `+field.Placeholder)
	}

	return fmt.Sprintf(`UPDATE %s SET %s`, table.Name, strings.Join(fieldStrs, `, `))
	// valueStrs := []string{}
	// values := query.ValueOf(stmt)

	// for i, field := range table.Fields {
	// 	fieldStrs = append(fieldStrs, field.Name)
	// 	if len(values) > i {
	// 		valueStrs = append(valueStrs, values[i])
	// 	} else {
	// 		valueStrs = append(valueStrs, `?`)
	// 	}
	// }
	// return fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s)`, tables[0].Name, strings.Join(fieldStrs, `, `), strings.Join(valueStrs, `, `))
}
