package sqlserver

import (
	"fmt"
	"strings"

	"github.com/eqto/dbm/stmt"
)

func insertStatement(s *stmt.Insert) string {
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
	sql := strings.Builder{}
	sql.WriteString(fmt.Sprintf(`INSERT INTO %s(%s)`, stmt.TableOf(s), strings.Join(fieldStrs, `, `)))
	if output := stmt.OutputOf(s); output != `` {
		sql.WriteString(fmt.Sprintf(` OUTPUT %s`, output))
	}

	sql.WriteString(fmt.Sprintf(` VALUES(%s)`, strings.Join(fieldValues, `, `)))
	return sql.String()
}
