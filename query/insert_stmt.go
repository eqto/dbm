package query

import "strings"

type InsertStmt struct {
	table  *Table
	values []string
}

func (i *InsertStmt) Values(query string) *InsertStmt {
	values := strings.Split(query, `,`)
	i.values = []string{}
	for _, val := range values {
		val = strings.TrimSpace(val)
		i.values = append(i.values, val)
	}
	return i
}

func InsertInto(table string, fields string) *InsertStmt {
	t := &Table{Name: table, Fields: parseFields(fields)}
	return &InsertStmt{table: t}
}
