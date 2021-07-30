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

func (w *WhereStmt) OrderBy(order string) *OrderByStmt {
	return w.ConditionStmt.stmt.(*TableStmt).OrderBy(order)
}
