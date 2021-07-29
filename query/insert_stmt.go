package query

type InsertStmt struct {
	table *Table
}

func InsertInto(table string, fields string) *InsertStmt {
	t := &Table{Name: table, Fields: parseFields(fields)}
	return &InsertStmt{table: t}
}
