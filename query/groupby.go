package query

type GroupBy struct {
	table  *TableStmt
	groups []string
}

func (g *GroupBy) OrderBy(orders string) *OrderBy {
	return g.table.OrderBy(orders)
}

func (g *GroupBy) Limit(num ...int) *TableStmt {
	return g.table.Limit(num...)
}
