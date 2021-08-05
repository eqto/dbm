package sqlserver

import (
	"fmt"
	"strings"

	"github.com/eqto/dbm/stmt"
)

func insertStatement(s *stmt.Insert) string {
	tableName := stmt.TableOf(s)
	fields := stmt.FieldsOf(s)

	fieldStrs := []string{}
	fieldValues := []string{}
	counter := 0
	for i, name := range fields.Names() {
		fieldStrs = append(fieldStrs, name)
		if value := fields.ValueByIndex(i); value != `` {
			fieldValues = append(fieldValues, value)
		} else {
			counter++
			fieldValues = append(fieldValues, fmt.Sprintf(`@p%d`, counter))
		}
	}

	return fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s)`, tableName, strings.Join(fieldStrs, `, `), strings.Join(fieldValues, `, `))
}
