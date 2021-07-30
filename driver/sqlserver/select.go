package sqlserver

import (
	"fmt"
	"strings"

	"github.com/eqto/go-db/query"
)

func querySelect(stmt *query.SelectStmt) string {
	sql := strings.Builder{}

	fields := []string{}
	tables := []string{}

	for _, field := range query.FieldsOf(stmt) {
		if field.Alias != `` {
			fields = append(fields, fmt.Sprintf(`%s AS %s`, field.Name, field.Alias))
		} else {
			fields = append(fields, field.Name)
		}
	}

	tableStmt := query.TableStmtOf(stmt)
	for {
		table := query.TableOf(tableStmt)
		tableStr := table.Name
		if table.Alias != `` {
			tableStr += ` ` + table.Alias
		}
		if len(tables) == 0 {
			tables = append(tables, tableStr)
		} else {
			tables = append(tables, fmt.Sprintf(`INNER JOIN %s ON %s`, tableStr, table.Condition))
		}

		tableStmt = query.JoinOf(tableStmt)
		if tableStmt == nil {
			break
		}
	}
	strFields := strings.Join(fields, `, `)
	if strFields == `` {
		strFields = `*`
	}
	sql.WriteString(fmt.Sprintf(`SELECT %s FROM %s`, strFields, strings.Join(tables, ` `)))

	if wheres := query.WhereOf(stmt); len(wheres) > 0 {
		sql.WriteString(fmt.Sprintf(` WHERE %s`, strings.Join(wheres, ` `)))
	}
	if groupBy := query.GroupByOf(stmt); len(groupBy) > 0 {
		sql.WriteString(fmt.Sprintf(` GROUP BY %s`, strings.Join(groupBy, ` `)))
	}
	if orderBys := query.OrderByOf(stmt); len(orderBys) > 0 {
		sql.WriteString(fmt.Sprintf(` ORDER BY %s`, strings.Join(orderBys, `, `)))
	}

	if offset, count := query.LimitOf(stmt); count > 0 {
		sql.WriteString(fmt.Sprintf(` OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, offset, count))
	}
	return sql.String()
}
