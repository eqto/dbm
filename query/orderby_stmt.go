package query

type OrderByStmt struct {
	table  *TableStmt
	orders []string
}

func (o *OrderByStmt) Limit(num ...int) *TableStmt {
	return o.table.Limit(num...)
}
