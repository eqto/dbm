package stmt

type OrderBy struct {
	stmt *Select
}

func (o *OrderBy) Limit(num ...int) *Select {
	return o.stmt.Limit(num...)
}
