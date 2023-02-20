package mysql

import (
	"fmt"
	"strings"

	"github.com/eqto/dbm/stmt"
)

func updateStatement(s *stmt.Update) string {
	sql := strings.Builder{}
	tableName := stmt.TableOf(s)

	sql.WriteString(fmt.Sprintf(`UPDATE %s SET %s`, tableName, strings.Join(stmt.NameValueOf(s), `, `)))

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
	if orderBy := stmt.OrderByOf(s); orderBy != `` {
		sql.WriteString(fmt.Sprintf(` ORDER BY %s`, orderBy))
	}
	if _, count := stmt.LimitOf(s); count > 0 {
		sql.WriteString(fmt.Sprintf(` LIMIT  %d`, count))
	}

	return sql.String()
}
