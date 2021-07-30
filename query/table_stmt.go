/*
 * @Author: tuxer
 * @Date: 2021-07-29 12:42:08
 * @Last Modified by: tuxer
 * @Last Modified time: 2021-07-30 10:09:58
 */

package query

import (
	"strings"
)

type Table struct {
	Name      string
	Alias     string
	Condition string
	fields    []Field
}

type TableStmt struct {
	stmt  interface{}
	table Table
	join  *TableStmt
}

func (t *TableStmt) InnerJoin(table, condition string) *TableStmt {
	fields := FieldsOf(t.stmt)
	t.join = parseTable(t.stmt, table, fields)
	t.join.table.Condition = condition
	return t.join
}

func (t *TableStmt) Where(condition string) *WhereStmt {
	where := &WhereStmt{table: t, conditions: []string{condition}}
	assignWhere(t.stmt, where)
	return where
}

//OrderBy
//query: "title" => Select books.* From books ORDER BY title
//query: "title DESC" => Select books.* From books ORDER BY title DESC
func (t *TableStmt) OrderBy(order string) *OrderBy {
	o := &OrderBy{table: t}
	split := strings.Split(order, `,`)
	for _, order := range split {
		o.orders = append(o.orders, strings.TrimSpace(order))
	}
	assignOrderBy(t.stmt, o)
	return o
}

//Limit used by MySQL. Parameters 'num' can be single int for "LIMIT n" or double for "LIMIT n1, n2"
//Ex:
//SELECT * FROM books LIMIT 1. offset = 0, length = 1
//SELECT * FROM books LIMIT 0, 10. offset = 0, length = 10
func (t *TableStmt) Limit(num ...int) *TableStmt {
	assignLimit(t.stmt, num...)
	return t
}

func parseTable(stmt interface{}, query string, fields []Field) *TableStmt {
	tableStmt := &TableStmt{stmt: stmt}
	split := strings.SplitN(query, ` `, 2)
	prefix := ``
	table := Table{}
	if len(split) == 2 {
		table.Name, table.Alias = strings.TrimSpace(split[0]), strings.TrimSpace(split[1])
		prefix = table.Alias + `.`
	} else {
		table.Name = strings.TrimSpace(query)
		prefix = table.Name + `.`
	}
	for _, field := range fields {
		if strings.HasPrefix(field.Name, prefix) {
			table.fields = append(table.fields, field)
		}
	}
	tableStmt.table = table
	return tableStmt
}
