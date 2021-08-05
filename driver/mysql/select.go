package mysql

import (
	"fmt"
	"strings"

	"github.com/eqto/dbm/stmt"
)

func selectStatement(s *stmt.Select) string {
	// sql := strings.Builder{}
	// fields := []string{}
	// tables := []string{}

	fieldStrs := []string{}
	tableStrs := []string{}

	fields := stmt.FieldsOf(s)

	for i, name := range fields.Names() {
		if alias := fields.AliasByIndex(i); alias != `` {
			fieldStrs = append(fieldStrs, fmt.Sprintf(`%s AS %s`, name, alias))
		} else {
			fieldStrs = append(fieldStrs, name)
		}
	}

	tables := stmt.TablesOf(s)
	for i, name := range tables.Names() {
		table := tables.TableByIndex(i)
		tableStr := name
		if table.Alias != `` {
			tableStr += ` ` + table.Alias
		}
		if i > 0 {
			join := `INNER`
			switch table.Join {
			case stmt.LeftJoin:
				join = `LEFT`
			case stmt.RightJoin:
				join = `RIGHT`
			}
			tableStr = fmt.Sprintf(`%s JOIN %s ON %s`, join, tableStr, table.JoinOn)
		}
		tableStrs = append(tableStrs, tableStr)
	}

	sql := strings.Builder{}
	sql.WriteString(fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(fieldStrs, `, `), strings.Join(tableStrs, ` `)))

	wheres := stmt.WheresOf(s)
	if len(wheres) > 0 {
		sql.WriteString(` WHERE `)
		for i, where := range wheres {
			if i > 0 {
				if where.Or {
					sql.WriteString(` OR `)
				} else {
					sql.WriteString(` AND `)
				}
			}
			sql.WriteString(where.Condition)
		}
	}

	if groupBy := stmt.GroupByOf(s); groupBy != `` {
		sql.WriteString(fmt.Sprintf(` GROUP BY %s`, groupBy))
	}
	if orderBy := stmt.OrderByOf(s); orderBy != `` {
		sql.WriteString(fmt.Sprintf(` ORDER BY %s`, orderBy))
	}
	if offset, count := stmt.LimitOf(s); count > 0 {
		sql.WriteString(fmt.Sprintf(` LIMIT %d, %d`, offset, count))
	}

	return sql.String()
}
