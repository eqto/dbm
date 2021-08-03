package query

type OrderBy struct {
	table  *TableStmt
	orders []string
}

func (o *OrderBy) Limit(num ...int) *TableStmt {
	return o.table.Limit(num...)
}
