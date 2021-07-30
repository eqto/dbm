package query

type WhereStmt struct {
	table      *TableStmt
	conditions []string
}

func (w *WhereStmt) And(query string) *WhereStmt {
	w.conditions = append(w.conditions, `AND `+query)
	return w
}

func (w *WhereStmt) Or(query string) *WhereStmt {
	w.conditions = append(w.conditions, `OR `+query)
	return w
}

func (w *WhereStmt) OrderBy(order string) *OrderByStmt {
	return w.table.OrderBy(order)
}
