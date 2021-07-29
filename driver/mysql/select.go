package mysql

import (
	"fmt"
	"strings"

	"github.com/eqto/go-db/query"
)

func querySelect(stmt *query.SelectStmt) string {
	sql := strings.Builder{}

	fields := []string{}
	tables := []string{}

	for i, table := range query.TableOf(stmt) {
		for _, field := range table.Fields {
			if field.Alias != `` {
				fields = append(fields, fmt.Sprintf(`%s AS %s`, field.Name, field.Alias))
			} else {
				fields = append(fields, field.Name)
			}
		}
		tableStr := table.Name
		if table.Alias != `` {
			tableStr += ` ` + table.Alias
		}
		if i == 0 {
			tables = append(tables, tableStr)
		} else {
			tables = append(tables, fmt.Sprintf(`INNER JOIN %s ON %s`, tableStr, table.Condition))
		}
	}
	sql.WriteString(fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(fields, `, `), strings.Join(tables, ` `)))

	strs := query.WhereOf(stmt)
	if len(strs) > 0 {
		sql.WriteString(fmt.Sprintf(` WHERE %s`, strings.Join(strs, ` `)))
	}
	strs = query.OrderByOf(stmt)
	if len(strs) > 0 {
		sql.WriteString(fmt.Sprintf(` ORDER BY %s`, strings.Join(strs, `, `)))
	}
	offset, count := query.LimitOf(stmt)
	if count > 0 {
		sql.WriteString(fmt.Sprintf(` LIMIT %d, %d`, offset, count))
	}
	return sql.String()
}
