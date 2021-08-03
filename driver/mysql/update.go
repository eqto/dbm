package mysql

import (
	"fmt"
	"strings"

	"github.com/eqto/dbq/query"
)

func queryUpdate(stmt *query.UpdateStmt) string {
	table := query.TableOf(stmt)
	fieldStrs := []string{}
	for _, field := range table.Fields {
		fieldStrs = append(fieldStrs, field.Name+` = `+field.Placeholder)
	}
	sql := strings.Builder{}
	sql.WriteString(fmt.Sprintf(`UPDATE %s SET %s`, table.Name, strings.Join(fieldStrs, `, `)))
	if conditions := query.WhereOf(stmt); len(conditions) > 0 {
		sql.WriteString(fmt.Sprintf(` WHERE %s`, strings.Join(conditions, ` `)))
	}

	return sql.String()
}
