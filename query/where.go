package query

type Where struct {
	Condition
}

func (w *Where) And(query string) *Where {
	w.Condition.And(query)
	return w
}

func (w *Where) Or(query string) *Where {
	w.Condition.Or(query)
	return w
}

func (w *Where) GroupBy(groupBy string) *GroupBy {
	return w.stmt.(*TableStmt).GroupBy(groupBy)
}

func (w *Where) OrderBy(orderBy string) *OrderBy {
	return w.Condition.stmt.(*TableStmt).OrderBy(orderBy)
}
