package query

type GroupByStmt struct {
	table  *TableStmt
	groups []string
}

func (g *GroupByStmt) OrderBy(orders string) *OrderByStmt {
	return g.table.OrderBy(orders)
}

func (g *GroupByStmt) Limit(num ...int) *TableStmt {
	return g.table.Limit(num...)
}
