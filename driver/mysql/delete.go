package mysql

import (
	"fmt"
	"strings"

	"github.com/eqto/dbm/stmt"
)

func deleteStatement(s *stmt.Delete) string {
	sql := strings.Builder{}
	tableName := stmt.TableOf(s)

	sql.WriteString(fmt.Sprintf(`DELETE FROM %s`, tableName))

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

	return sql.String()
}
