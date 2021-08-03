package query

import (
	"strings"
)

type Table struct {
	Name      string
	Alias     string
	Condition string
	Fields    []Field
}

type TableStmt struct {
	stmt     interface{}
	table    Table
	joinTo   *TableStmt
	joinKind string
}

func (t *TableStmt) Join(joinKind, table, condition string) *TableStmt {
	t.joinTo = parseTable(t.stmt, table)
	t.joinTo.table.Condition = condition
	t.joinKind = joinKind
	return t.joinTo
}

func (t *TableStmt) InnerJoin(table, condition string) *TableStmt {
	return t.Join(`INNER`, table, condition)
}

func (t *TableStmt) LeftJoin(table, condition string) *TableStmt {
	return t.Join(`LEFT`, table, condition)
}

func (t *TableStmt) RightJoin(table, condition string) *TableStmt {
	return t.Join(`RIGHT`, table, condition)
}

func (t *TableStmt) Where(condition string) *Where {
	where := &Where{Condition: Condition{stmt: t, conditions: []string{condition}}}
	assignWhere(t.stmt, where)
	return where
}

func (t *TableStmt) GroupBy(groupBy string) *GroupBy {
	split := strings.Split(groupBy, `,`)
	groups := []string{}
	for _, s := range split {
		groups = append(groups, strings.TrimSpace(s))
	}
	g := &GroupBy{groups: groups}
	assignGroupBy(t.stmt, g)
	return g
}

//OrderBy
//query: "title" => Select books.* From books ORDER BY title
//query: "title DESC" => Select books.* From books ORDER BY title DESC
func (t *TableStmt) OrderBy(orders string) *OrderBy {
	o := &OrderBy{table: t}
	split := strings.Split(orders, `,`)
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

func parseTable(stmt interface{}, query string) *TableStmt {
	tableStmt := &TableStmt{stmt: stmt}
	split := strings.SplitN(query, ` `, 2)
	table := Table{}
	if len(split) == 2 {
		table.Name, table.Alias = strings.TrimSpace(split[0]), strings.TrimSpace(split[1])
	} else {
		table.Name = strings.TrimSpace(query)
	}
	tableStmt.table = table
	return tableStmt
}
