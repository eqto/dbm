package query

type WhereStmt struct {
	ConditionStmt
}

func (w *WhereStmt) And(query string) *WhereStmt {
	w.ConditionStmt.And(query)
	return w
}

func (w *WhereStmt) Or(query string) *WhereStmt {
	w.ConditionStmt.Or(query)
	return w
}

func (w *WhereStmt) GroupBy(groupBy string) *GroupByStmt {
	return w.stmt.(*TableStmt).GroupBy(groupBy)
}

func (w *WhereStmt) OrderBy(orderBy string) *OrderByStmt {
	return w.ConditionStmt.stmt.(*TableStmt).OrderBy(orderBy)
}
