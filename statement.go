package dbq

import (
	"github.com/eqto/dbq/stmt"
)

func Select(fields string) *stmt.SelectFields {
	return stmt.Build().Select(fields)
}

func InsertInto(table, fields string) *stmt.Insert {
	return stmt.Build().InsertInto(table, fields)
}

func Update(table string) *stmt.Update {
	return stmt.Build().Update(table)
}
