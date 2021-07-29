package mysql

import (
	"fmt"
	"strings"

	"github.com/eqto/go-db/query"
)

func queryInsert(stmt *query.InsertStmt) string {
	tables := query.TableOf(stmt)
	if len(tables) == 0 {
		return ``
	}
	table := tables[0]
	fields := []string{}
	values := []string{}
	for _, field := range table.Fields {
		fields = append(fields, field.Name)
		values = append(values, `?`)
	}
	return fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s)`, tables[0].Name, strings.Join(fields, `, `), strings.Join(values, `, `))
}
