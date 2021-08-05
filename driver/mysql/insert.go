package mysql

import (
	"fmt"
	"strings"

	"github.com/eqto/dbq/stmt"
)

func insertStatement(s *stmt.Insert) string {

	tableName := stmt.TableOf(s)
	fields := stmt.FieldsOf(s)

	fieldStrs := []string{}
	fieldValues := []string{}
	for i, name := range fields.Names() {
		fieldStrs = append(fieldStrs, name)
		if value := fields.ValueByIndex(i); value != `` {
			fieldValues = append(fieldValues, value)
		} else {
			fieldValues = append(fieldValues, `?`)
		}
	}

	return fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s)`, tableName, strings.Join(fieldStrs, `, `), strings.Join(fieldValues, `, `))
}
