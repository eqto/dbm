package stmt

type GroupBy struct {
	stmt *Select
}

func (g *GroupBy) OrderBy(orderBy string) *OrderBy {
	return g.stmt.OrderBy(orderBy)
}

func (g *GroupBy) Limit(num ...int) *Select {
	return g.stmt.Limit(num...)
}
